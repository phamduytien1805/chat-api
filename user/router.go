package main

import (
	"context"

	"github.com/phamduytien1805/package/server"
)

type Router struct {
	grpcServer server.GrpcServer
	httpServer server.HttpServer
}

func NewRouter(grpcServer server.GrpcServer, httpServer server.HttpServer) *Router {
	return &Router{grpcServer, httpServer}
}

func (r *Router) Run() {
	r.grpcServer.Register()
	r.grpcServer.Run()

	r.httpServer.RegisterRoutes()
	r.httpServer.Run()
}
func (r *Router) GracefulStop(ctx context.Context) error {
	r.httpServer.GracefulStop(ctx)
	return r.grpcServer.GracefulStop()
}
