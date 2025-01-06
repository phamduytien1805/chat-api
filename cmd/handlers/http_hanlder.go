package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/phamduytien1805/internal/platform/redis_engine"
	"github.com/phamduytien1805/internal/user"
	"github.com/phamduytien1805/package/config"
	"github.com/phamduytien1805/package/http_utils"
	"github.com/phamduytien1805/package/token"
	"github.com/phamduytien1805/package/validator"
)

type HttpServer struct {
	httpServer *http.Server
	config     *config.Config
	logger     *slog.Logger
	validator  *validator.Validate
	httpPort   string
	router     *chi.Mux
	userSvc    user.UserSvc
	tokenMaker token.Maker
	redis      redis_engine.RedisQuerier
}

func NewHttpServer(config *config.Config, logger *slog.Logger, validator *validator.Validate, tokenMaker token.Maker, userSvc user.UserSvc, redis redis_engine.RedisQuerier) *HttpServer {
	return &HttpServer{
		config:     config,
		logger:     logger,
		validator:  validator,
		httpPort:   config.Web.Http.Server.Port,
		userSvc:    userSvc,
		tokenMaker: tokenMaker,
		redis:      redis,
	}
}

func (s *HttpServer) RegisterRoutes() {
	s.router = chi.NewRouter()
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"}, //NOTE: just for development purpose
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	s.router.Use(middleware.Heartbeat("/ping"))

	s.router.NotFound(http_utils.NotFoundResponse)
	s.router.MethodNotAllowed(http_utils.MethodNotAllowedResponse)

	s.router.Route(("/auth"), func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Post("/login", s.authenticateUserBasic)
			r.Post("/register", s.registerUser)
			r.Post("/token", s.refreshToken)
		})
	})
	s.router.Route("/user", func(r chi.Router) {
		r.Use(s.authenticator)
		r.Get("/", s.getUser)
	})

}
func (s *HttpServer) Run() {
	go func() {
		addr := ":" + s.httpPort
		s.httpServer = &http.Server{
			Addr:    addr,
			Handler: s.router,
		}
		s.logger.Info("http server listening", slog.String("addr", addr))
		err := s.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.logger.Error(err.Error())
			os.Exit(1)
		}
	}()
}

func (r *HttpServer) GracefulStop(ctx context.Context) error {

	err := r.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
