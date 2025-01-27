package main

import (
	"log/slog"
	"os"

	"github.com/gocql/gocql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/phamduytien1805/cmd/handlers"
	"github.com/phamduytien1805/internal/auth"
	"github.com/phamduytien1805/internal/platform/db"
	"github.com/phamduytien1805/internal/platform/mail"
	"github.com/phamduytien1805/internal/platform/redis_engine"
	"github.com/phamduytien1805/internal/platform/scylladb"
	"github.com/phamduytien1805/internal/taskq"
	"github.com/phamduytien1805/internal/user"
	"github.com/phamduytien1805/package/config"
	"github.com/phamduytien1805/package/hash_generator"
	"github.com/phamduytien1805/package/server"
	"github.com/phamduytien1805/package/token"
	"github.com/phamduytien1805/package/validator"
	"github.com/redis/go-redis/v9"
	"github.com/scylladb/gocqlx/v3"
)

type InfraStruct struct {
	pgConn      *pgxpool.Pool
	redisClient *redis.Client
	cqlSession  *gocql.Session
}

func (i *InfraStruct) Close() error {
	i.pgConn.Close()
	i.redisClient.Close()
	i.cqlSession.Close()
	return nil
}

func ServerBuilder() (*server.Server, error) {
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cluster := scylladb.NewClusterManager(configConfig.Scylla)
	cqlSession, err := gocqlx.WrapSession(gocql.NewSession(*cluster.Cluster))
	if err != nil {
		return nil, err
	}

	pgConn, err := db.NewPostgresql(configConfig)
	if err != nil {
		return nil, err
	}

	store := db.NewStore(pgConn)

	redisQuerier := redis_engine.NewRedis(configConfig)
	redisStore := redis_engine.NewRedisStore(redisQuerier)

	validator := validator.New()
	hashGen := hash_generator.NewArgon2idHash(configConfig)
	tokenMaker, err := token.NewJWTMaker(configConfig.Token.SecretKey)
	if err != nil {
		return nil, err
	}

	taskqProducer := taskq.NewTaskProducer(configConfig.Redis, logger)

	mailSvc := mail.NewMailService(configConfig.Mail, logger)
	authSvc := auth.NewAuthService(configConfig, logger, tokenMaker, taskqProducer, redisStore, store)
	userSvc := user.NewUserServiceImpl(store, configConfig, logger, hashGen)

	taskqServer := taskq.NewTaskConsumer(configConfig.Redis, logger, mailSvc)
	httpServer := handlers.NewHttpServer(configConfig, logger, validator, authSvc, userSvc)
	router := NewRouter(httpServer, taskqServer)

	infraCloser := &InfraStruct{
		pgConn:      pgConn,
		redisClient: redisQuerier,
		cqlSession:  cqlSession.Session,
	}

	return server.NewServer(router, infraCloser), nil

}
