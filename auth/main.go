package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	http_adapter "github.com/phamduytien1805/auth/infrastructures/http"
	"github.com/phamduytien1805/auth/infrastructures/mail"
	"github.com/phamduytien1805/auth/infrastructures/taskq"
	tokensvc "github.com/phamduytien1805/auth/infrastructures/token"
	"github.com/phamduytien1805/auth/infrastructures/userclient"
	"github.com/phamduytien1805/auth/usecase"
	"github.com/phamduytien1805/package/config"
	redis_engine "github.com/phamduytien1805/package/redis"
	"github.com/phamduytien1805/package/server"
	"github.com/phamduytien1805/package/validator"
	"github.com/spf13/viper"
)

var cfgFile string

func initConfig() {

	viper.SetConfigType("yaml")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func main() {
	flag.StringVar(&cfgFile, "config", "", "config file path")
	flag.Parse()

	initConfig()

	s, err := AppBuilder()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init server: %v\n", err)
		os.Exit(1)
	}

	s.Serve()
}

func AppBuilder() (*server.Server, error) {
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	redisQuerier := redis_engine.NewRedis(configConfig)
	redisStore := redis_engine.NewRedisStore(redisQuerier)

	validator := validator.New()
	tokenSvc, err := tokensvc.NewTokenService(configConfig.Token, redisStore)
	if err != nil {
		return nil, err
	}

	userSvc, err := userclient.NewUserClientService(configConfig.User)
	if err != nil {
		return nil, err
	}

	workerProducer := taskq.NewTaskProducer(configConfig.Redis, logger)

	mailSvc := mail.NewMailService(configConfig.Mail, logger, redisStore, workerProducer, userSvc)
	worker := taskq.NewTaskConsumer(configConfig.Redis, logger, mailSvc)

	infra := NewInfraCloser()
	httpServer := http_adapter.NewHttpServer(configConfig.Auth, validator, &http_adapter.Usecases{
		Login:             usecase.NewLoginUsecase(logger, userSvc, tokenSvc),
		Register:          usecase.NewRegisterUsecase(logger, userSvc, mailSvc, tokenSvc),
		VerifyEmail:       usecase.NewVerifyEmailUsecase(logger, mailSvc),
		ResendEmail:       usecase.NewResendEmailUsecase(logger, mailSvc, userSvc),
		RefreshToken:      usecase.NewRefreshTokenUsecase(logger, tokenSvc, userSvc),
		VerifyAccessToken: usecase.NewVerifyAccessTokenUsecase(logger, tokenSvc),
		Logout:            usecase.NewLogoutUsecase(logger, tokenSvc),
	})
	router := NewRouter(httpServer, worker)

	return server.NewServer(router, infra), nil

}
