package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/phamduytien1805/cmd/handlers"
	"github.com/phamduytien1805/internal/platform/db"
	"github.com/phamduytien1805/internal/user"
	"github.com/phamduytien1805/package/config"
	"github.com/phamduytien1805/package/hash_generator"
	"github.com/phamduytien1805/package/server"
	"github.com/phamduytien1805/package/token"
	"github.com/phamduytien1805/package/validator"
	"github.com/spf13/viper"
)

var cfgFile string

type InfraStruct struct {
	pgConn *pgxpool.Pool
}

func (i *InfraStruct) Close() error {
	i.pgConn.Close()
	return nil
}

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

func initServer() (*server.Server, error) {
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

	validator := validator.New()
	hashGen := hash_generator.NewArgon2idHash(configConfig)
	tokenMaker, err := token.NewJWTMaker(configConfig.Token.SecretKey)
	if err != nil {
		return nil, err
	}
	userSvc := user.NewUserServiceImpl(store, configConfig, logger, hashGen)
	httpServer := handlers.NewHttpServer(configConfig, logger, validator, tokenMaker, userSvc)
	router := handlers.NewRouter(httpServer)

	infraCloser := &InfraStruct{
		pgConn: pgConn,
	}

	return server.NewServer(router, infraCloser), nil

}

func main() {
	initConfig()

	s, err := initServer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init server: %v\n", err)
		os.Exit(1)
	}

	s.Serve()
}
