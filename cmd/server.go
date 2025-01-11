package main

import (
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/phamduytien1805/cmd/handlers"
	"github.com/phamduytien1805/internal/auth"
	"github.com/phamduytien1805/internal/platform/db"
	"github.com/phamduytien1805/internal/platform/mail"
	"github.com/phamduytien1805/internal/platform/redis_engine"
	"github.com/phamduytien1805/internal/taskq"
	"github.com/phamduytien1805/internal/user"
	"github.com/phamduytien1805/package/config"
	"github.com/phamduytien1805/package/hash_generator"
	"github.com/phamduytien1805/package/server"
	"github.com/phamduytien1805/package/token"
	"github.com/phamduytien1805/package/validator"
)

type InfraStruct struct {
	pgConn *pgxpool.Pool
}

func (i *InfraStruct) Close() error {
	i.pgConn.Close()
	return nil
}

func ServerBuilder() (*server.Server, error) {
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	pgConn, err := db.NewPostgresql(configConfig)
	if err != nil {
		return nil, err
	}

	store := db.NewStore(pgConn)

	redisQuerier := redis_engine.NewRedis(configConfig)

	validator := validator.New()
	hashGen := hash_generator.NewArgon2idHash(configConfig)
	tokenMaker, err := token.NewJWTMaker(configConfig.Token.SecretKey)
	if err != nil {
		return nil, err
	}

	taskqProducer := taskq.NewTaskProducer(configConfig.Redis, logger)

	mailSvc := mail.NewMailService(configConfig.Mail, logger)
	authSvc := auth.NewAuthService(configConfig, logger, tokenMaker, taskqProducer, redisQuerier)
	userSvc := user.NewUserServiceImpl(store, configConfig, logger, hashGen)

	taskqServer := taskq.NewTaskConsumer(configConfig.Redis, logger, mailSvc)
	httpServer := handlers.NewHttpServer(configConfig, logger, validator, authSvc, userSvc)
	router := NewRouter(httpServer, taskqServer)

	infraCloser := &InfraStruct{
		pgConn: pgConn,
	}

	return server.NewServer(router, infraCloser), nil

}
