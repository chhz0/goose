package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultTimeout = 30 * time.Second
	maxRetries     = 3
	retryDelay     = 500 * time.Millisecond
)

const (
	ContentTypeJSON      = "application/json"
	ContentTypeForm      = "application/x-www-form-urlencoded"
	ContentTypeMultipart = "multipart/form-data"
	ContentTypeXML       = "application/xml"
	ContentTypeText      = "text/plain"
)

var defaultClient = NewClient()

type Client struct {
	baseURL    string
	httpClient *http.Client
	headers    map[string]string
}

type ClientOption func(*Client)

func NewClient(opts ...ClientOption) *Client {
	client := &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		headers: make(map[string]string),
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(baseURL, "/")
	}
}

func WithHeaders(headers map[string]string) ClientOption {
	return func(c *Client) {
		for k, v := range headers {
			c.headers[k] = v
		}
	}
}

// RequestBuilder is a builder for building HTTP requests.
type RequestBuilder struct {
	client      *Client
	method      string
	url         string
	headers     map[string]string
	queryParams url.Values
	pathParams  map[string]string
	body        interface{}
	bodyType    string
	formData    url.Values
	retries     int
	files       map[string]string
}

func (c *Client) newRequestBuilder(method, path string) *RequestBuilder {
	return &RequestBuilder{
		client:      c,
		method:      strings.ToUpper(method),
		url:         c.buildURL(path),
		headers:     make(map[string]string),
		queryParams: make(url.Values),
		pathParams:  make(map[string]string),
		formData:    make(url.Values),
		files:       make(map[string]string),
		retries:     maxRetries,
	}
}

func (c *Client) buildURL(path string) string {
	if c.baseURL == "" {
		return path
	}
	return fmt.Sprintf("%s/%s", c.baseURL, strings.TrimLeft(path, "/"))
}

func (c *Client) Get(url string) *RequestBuilder     { return c.newRequestBuilder("GET", url) }
func (c *Client) Post(url string) *RequestBuilder    { return c.newRequestBuilder("POST", url) }
func (c *Client) Put(url string) *RequestBuilder     { return c.newRequestBuilder("PUT", url) }
func (c *Client) Delete(url string) *RequestBuilder  { return c.newRequestBuilder("DELETE", url) }
func (c *Client) Patch(url string) *RequestBuilder   { return c.newRequestBuilder("PATCH", url) }
func (c *Client) Head(url string) *RequestBuilder    { return c.newRequestBuilder("HEAD", url) }
func (c *Client) Options(url string) *RequestBuilder { return c.newRequestBuilder("OPTIONS", url) }

func (rb *RequestBuilder) AddHeader(key, value string) *RequestBuilder {
	rb.headers[key] = value
	return rb
}

func (rb *RequestBuilder) AddQueryParam(key, value string) *RequestBuilder {
	rb.queryParams.Add(key, value)
	return rb
}

func (rb *RequestBuilder) AddPathParam(key, value string) *RequestBuilder {
	rb.pathParams[key] = value
	return rb
}

func (rb *RequestBuilder) SetJSONBody(body interface{}) *RequestBuilder {
	rb.body = body
	rb.bodyType = ContentTypeJSON
	return rb
}

func (rb *RequestBuilder) SetFormData(data map[string]string) *RequestBuilder {
	for k, v := range data {
		rb.formData.Add(k, v)
	}
	rb.bodyType = ContentTypeForm
	return rb
}

func (rb *RequestBuilder) AddFile(fileName, filePath string) *RequestBuilder {
	rb.files[fileName] = filePath
	rb.bodyType = ContentTypeMultipart
	return rb
}

func (rb *RequestBuilder) SetRetries(retries int) *RequestBuilder {
	rb.retries = retries
	return rb
}

func (rb *RequestBuilder) buildRequest() (*http.Request, error) {
	finalURL := rb.url

	// process path params
	for k, v := range rb.pathParams {
		param := ":" + k
		if !strings.Contains(finalURL, param) {
			return nil, fmt.Errorf("path parameter %s not found in url %s", k, finalURL)
		}
		finalURL = strings.ReplaceAll(finalURL, param, url.PathEscape(v))
	}

	// add query params
	if len(rb.queryParams) > 0 {
		finalURL += "?" + rb.queryParams.Encode()
	}

	// prepare request body
	var body io.Reader
	contentType := ""

	switch rb.bodyType {
	case ContentTypeJSON:
		if rb.body != nil {
			jsonData, err := json.Marshal(rb.body)
			if err != nil {
				return nil, err
			}
			body = bytes.NewBuffer(jsonData)
			contentType = ContentTypeJSON
		}
	case ContentTypeForm:
		if len(rb.formData) > 0 {
			body = strings.NewReader(rb.formData.Encode())
			contentType = ContentTypeForm
		}
	case ContentTypeMultipart:
		if len(rb.files) > 0 || len(rb.formData) > 0 {
			var buf bytes.Buffer
			writer := multipart.NewWriter(&buf)

			for k, values := range rb.formData {
				for _, v := range values {
					if err := writer.WriteField(k, v); err != nil {
						return nil, err
					}
				}
			}
			for field, filePath := range rb.files {
				file, err := os.Open(filePath)
				if err != nil {
					return nil, err
				}
				defer file.Close()

				part, err := writer.CreateFormFile(field, filepath.Base(filePath))
				if err != nil {
					return nil, err
				}

				if _, err := io.Copy(part, file); err != nil {
					return nil, err
				}
			}
			if err := writer.Close(); err != nil {
				return nil, err
			}

			body = &buf
			contentType = writer.FormDataContentType()
		}
	}

	req, err := http.NewRequest(rb.method, finalURL, body)
	if err != nil {
		return nil, err
	}

	mergeHeaders(req, rb.headers, rb.client.headers)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	return req, nil
}

func (rb *RequestBuilder) Do() (*Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= rb.retries; attempt++ {
		req, _ := rb.buildRequest()

		ctx, cancel := context.WithTimeout(context.Background(), rb.client.httpClient.Timeout)

		req = req.WithContext(ctx)

		resp, err = rb.client.httpClient.Do(req)
		cancel()

		if shouldRetry(err) && attempt < rb.retries {
			time.Sleep(retryDelay * time.Duration(1<<attempt))
			continue
		}
		break
	}

	if err != nil {
		return nil, fmt.Errorf("request failed after %d attempts: %w", rb.retries, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header.Clone(),
		body:       body,
	}, nil
}

func mergeHeaders(req *http.Request, headers ...map[string]string) {
	for _, header := range headers {
		for k, v := range header {
			if strings.EqualFold(k, "host") && req.Host == "" {
				req.Host = v
				continue
			}
			req.Header.Set(k, v)
		}
	}
}

func shouldRetry(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		if urlErr.Timeout() || urlErr.Temporary() {
			return true
		}
	}

	return false
}

type Response struct {
	StatusCode int
	Headers    http.Header
	body       []byte
}

func (r *Response) JSON(v interface{}) error {
	return json.Unmarshal(r.body, v)
}

func (r *Response) Text() string {
	return string(r.body)
}

func (r *Response) OK() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}
func (r *Response) Created() bool {
	return r.StatusCode == http.StatusCreated
}

func (r *Response) NoContent() bool {
	return r.StatusCode == http.StatusNoContent
}
func SetBaseURL(baseURL string) {
	defaultClient.baseURL = baseURL
}

type RequestOptions func(*RequestBuilder)

func WithPathParams(params map[string]string) RequestOptions {
	return func(rb *RequestBuilder) {
		for k, v := range params {
			rb.AddPathParam(k, v)
		}
	}
}

func WithQueryParams(params map[string]string) RequestOptions {
	return func(rb *RequestBuilder) {
		for k, v := range params {
			rb.AddQueryParam(k, v)
		}
	}
}

func WithRequestHeaders(headers map[string]string) RequestOptions {
	return func(rb *RequestBuilder) {
		for k, v := range headers {
			rb.AddHeader(k, v)
		}
	}
}
func WithJSONBody(body interface{}) RequestOptions {
	return func(rb *RequestBuilder) {
		rb.SetJSONBody(body)
	}
}

func WithFormData(data map[string]string) RequestOptions {
	return func(rb *RequestBuilder) {
		rb.SetFormData(data)
	}
}

func WithFile(fileName, filePath string) RequestOptions {
	return func(rb *RequestBuilder) {
		rb.AddFile(fileName, filePath)
	}
}

func Get(path string, opts ...RequestOptions) (*Response, error) {
	return doRequest(defaultClient.Get(path), opts...)
}

func Post(path string, opts ...RequestOptions) (*Response, error) {
	return doRequest(defaultClient.Post(path), opts...)
}
func Put(path string, opts ...RequestOptions) (*Response, error) {
	return doRequest(defaultClient.Put(path), opts...)
}

func Delete(path string, opts ...RequestOptions) (*Response, error) {
	return doRequest(defaultClient.Delete(path), opts...)
}
func Patch(path string, opts ...RequestOptions) (*Response, error) {
	return doRequest(defaultClient.Patch(path), opts...)
}
func doRequest(rb *RequestBuilder, opts ...RequestOptions) (*Response, error) {
	for _, opt := range opts {
		opt(rb)
	}

	return rb.Do()
}
