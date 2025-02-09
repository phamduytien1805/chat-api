package main

import (
	"context"

	"github.com/phamduytien1805/package/server"
)

type Router struct {
	httpServer server.HttpServer
}

func NewRouter(httpServer server.HttpServer) *Router {
	return &Router{httpServer}
}

func (r *Router) Run() {
	r.httpServer.RegisterRoutes()
	r.httpServer.Run()
}
func (r *Router) GracefulStop(ctx context.Context) error {
	return r.httpServer.GracefulStop(ctx)
}
