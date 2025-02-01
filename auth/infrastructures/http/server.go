package http_adapter

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/phamduytien1805/auth/usecase"
	"github.com/phamduytien1805/package/common"
	"github.com/phamduytien1805/package/config"
	"github.com/phamduytien1805/package/http_utils"
	"github.com/phamduytien1805/package/server"
	"github.com/phamduytien1805/package/validator"
)

type Usecases struct {
	Login             *usecase.LoginUsecase
	Logout            *usecase.LogoutUsecase
	Register          *usecase.RegisterUsecase
	VerifyEmail       *usecase.VerifyEmailUsecase
	ResendEmail       *usecase.ResendEmailUsecase
	RefreshToken      *usecase.RefreshTokenUsecase
	VerifyAccessToken *usecase.VerifyAccessTokenUsecase
}

type HttpServer struct {
	httpServer *http.Server
	logger     common.HttpLog
	validator  *validator.Validate
	httpPort   string
	router     *chi.Mux

	// usecase
	uc *Usecases
}

func NewHttpServer(config *config.AuthConfig, validator *validator.Validate, uc *Usecases) server.HttpServer {
	return &HttpServer{
		logger:    common.NewHttpLog(),
		validator: validator,
		httpPort:  config.Http.Server.Port,
		uc:        uc,
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
			r.Post("/logout", s.logout)
			r.Post("/token", s.refreshToken)
			r.Post("/verify-email", s.verifyEmailUser)
			r.With(s.authenticator).Post("/resend-verification", s.resendEmailVerification)
		})
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
