package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
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
}

func NewHttpServer(config *config.Config, logger *slog.Logger, validator *validator.Validate, tokenMaker token.Maker, userSvc user.UserSvc) *HttpServer {
	return &HttpServer{
		config:     config,
		logger:     logger,
		validator:  validator,
		httpPort:   config.Web.Http.Server.Port,
		userSvc:    userSvc,
		tokenMaker: tokenMaker,
	}
}

func (s *HttpServer) RegisterRoutes() {
	s.router = chi.NewRouter()
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	s.router.Use(middleware.Heartbeat("/ping"))

	s.router.NotFound(http_utils.NotFoundResponse)
	s.router.MethodNotAllowed(http_utils.MethodNotAllowedResponse)

	s.router.Route("/user", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(s.authenticator)
			r.Get("/", s.getUser)
		})
		r.Group(func(r chi.Router) {
			r.Post("/register", s.registerUser)
			r.Post("/auth", s.authenticateUserBasic)
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

type contextKey string

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = contextKey("authorization_payload")
)

func (s *HttpServer) authenticator(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authorizationHeader := r.Header.Get(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			s.logger.Error("missing authorization header")
			err := errors.New("authorization header is not provided")
			http_utils.InvalidAuthenticateResponse(w, r, err)
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			http_utils.InvalidAuthenticateResponse(w, r, err)
			return
		}
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			http_utils.InvalidAuthenticateResponse(w, r, err)
			return
		}
		accessToken := fields[1]
		payload, err := s.tokenMaker.VerifyToken(accessToken)
		if err != nil {
			err := errors.New("invalid authorization header token")
			http_utils.InvalidAuthenticateResponse(w, r, err)
			return
		}
		ctx = context.WithValue(ctx, authorizationPayloadKey, *payload)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
