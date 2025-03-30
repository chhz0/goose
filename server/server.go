package server

import (
	"context"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"
)

type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

type ServerPlur struct {
	servers []Server
	eg      *errgroup.Group
}

func NewServerPlur() *ServerPlur {
	group, _ := errgroup.WithContext(context.Background())
	return &ServerPlur{
		eg: group,
	}
}

func (s *ServerPlur) AddServer(server Server) {
	s.servers = append(s.servers, server)
}

func (s *ServerPlur) StartAll() error {
	for _, server := range s.servers {
		s.eg.Go(server.ListenAndServe)
	}

	return s.eg.Wait()
}

func (s *ServerPlur) ShutdownAll(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for _, server := range s.servers {
		s.eg.Go(func() error {
			return server.Shutdown(ctx)
		})
	}
	return s.eg.Wait()
}

func (s *ServerPlur) RunOrDie(sig ...os.Signal) error {
	if err := s.StartAll(); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, sig...)
	<-quit

	return s.ShutdownAll(5 * time.Second)
}
