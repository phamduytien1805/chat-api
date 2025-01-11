package main

import (
	"context"

	"github.com/phamduytien1805/package/server"
)

type Router struct {
	httpServer  server.HttpServer
	taskqServer server.TaskQServer
}

func NewRouter(httpServer server.HttpServer, taskq server.TaskQServer) *Router {
	return &Router{httpServer, taskq}
}

func (r *Router) Run() {
	r.httpServer.RegisterRoutes()
	r.taskqServer.RegisterHandlers()

	r.httpServer.Run()
	r.taskqServer.Run()
}
func (r *Router) GracefulStop(ctx context.Context) error {
	return r.httpServer.GracefulStop(ctx)
}
