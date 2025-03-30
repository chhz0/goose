package rest

import (
	"testing"
)

func Test_GetWithBaseURL(t *testing.T) {
	SetBaseURL("http://localhost:8080")
	resp, _ := Get("/ping/:id",
		WithPathParams(map[string]string{"id": "123"}),
		WithQueryParams(map[string]string{"details": "true"}),
		WithRequestHeaders(map[string]string{"Authorization": "Bearer token"}),
	)
	t.Log(resp.Text())
}

func Test_Get(t *testing.T) {
	resp, _ := Get("http://localhost:8080/ping/")
	t.Log(resp.Text())
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func Test_Post(t *testing.T) {
	SetBaseURL("http://localhost:8080")
	resp, _ := Post("/user",
		WithJSONBody(User{Name: "John", Email: "john@example.com"}),
		WithRequestHeaders(map[string]string{
			"X-Request-ID": "ncahdlai",
		}),
	)
	t.Log(resp.Text())
	var user User
	err := resp.JSON(&user)
	t.Log(err)
	t.Log(user.Name)
}
