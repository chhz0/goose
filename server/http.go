package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chhz0/goose/server/engines"
)

type HttpConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	TLS          *TLSConfig
}

func (hc *HttpConfig) check() {
	if hc.Addr == "" {
		hc.Addr = ":8080"
	}
	if hc.ReadTimeout == 0 {
		hc.ReadTimeout = 5 * time.Second
	}
	if hc.WriteTimeout == 0 {
		hc.WriteTimeout = 10 * time.Second
	}
}

type TLSConfig struct {
	Cert string
	Key  string
}

type httpServer struct {
	cfg    *HttpConfig
	server *http.Server
}

// Listen implements Server.
func (hs *httpServer) ListenAndServe() error {
	errChan := make(chan error, 1)
	defer close(errChan)
	go func(errChan chan error) {
		if hs.cfg.TLS != nil {
			if err := hs.server.ListenAndServeTLS(hs.cfg.TLS.Cert, hs.cfg.TLS.Key); err != nil &&
				err != http.ErrServerClosed {
				errChan <- err
			}
		} else {
			if err := hs.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errChan <- err
			}
		}
	}(errChan)

	if err := <-errChan; err != nil {
		return err
	}

	return hs.wait()
}

// Shutdown implements Server.
func (hs *httpServer) Shutdown(ctx context.Context) error {
	return hs.server.Shutdown(ctx)
}

func (hs *httpServer) wait() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(quit)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return hs.server.Shutdown(ctx)
}

func NewHttp(cfg *HttpConfig, engine engines.Handler) Server {
	cfg.check()
	return &httpServer{
		cfg: cfg,
		server: &http.Server{
			Addr:         cfg.Addr,
			Handler:      engine.Handler(),
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
	}
}
