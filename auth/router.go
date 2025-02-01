package main

import (
	"context"

	"github.com/phamduytien1805/auth/domain"
	"github.com/phamduytien1805/package/server"
)

type Router struct {
	httpServer server.HttpServer
	worker     domain.TaskConsumer
}

func NewRouter(httpServer server.HttpServer, worker domain.TaskConsumer) *Router {
	return &Router{httpServer, worker}
}

func (r *Router) Run() {
	r.httpServer.RegisterRoutes()
	r.httpServer.Run()

	r.worker.RegisterHandlers()
	r.worker.Run()
}
func (r *Router) GracefulStop(ctx context.Context) error {
	err := r.worker.GracefulStop(ctx)
	if err != nil {
		return err
	}
	return r.httpServer.GracefulStop(ctx)
}
