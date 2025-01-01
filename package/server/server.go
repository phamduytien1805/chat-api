package server

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HttpServer interface {
	RegisterRoutes()
	Run()
	GracefulStop(ctx context.Context) error
}

type GrpcServer interface {
	Register()
	Run()
	GracefulStop() error
}

type Router interface {
	Run()
	GracefulStop(ctx context.Context) error
}

type InfraCloser interface {
	Close() error
}

type Server struct {
	router      Router
	infraCloser InfraCloser
}

func NewServer(router Router, infraCloser InfraCloser) *Server {
	return &Server{router, infraCloser}
}

func (s *Server) Serve() {
	s.router.Run()

	done := make(chan bool, 1)
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.GracefulStop(ctx, done)
	}()

	<-done
}

func (s *Server) GracefulStop(ctx context.Context, done chan bool) {
	err := s.router.GracefulStop(ctx)
	if err != nil {
		slog.Error(err.Error())
	}

	if err = s.infraCloser.Close(); err != nil {
		slog.Error(err.Error())
	}

	slog.Info("gracefully shutdowned")
	done <- true
}
