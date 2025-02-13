package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/phamduytien1805/hub/infras/db"
	"github.com/phamduytien1805/hub/infras/http"
	"github.com/phamduytien1805/hub/usecase"
	"github.com/phamduytien1805/package/config"
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

	pgConn, err := db.NewPostgresql(configConfig.DB)
	if err != nil {
		return nil, err
	}

	store := db.NewStore(pgConn)

	validator := validator.New()

	infra := NewInfraCloser()

	getDMUc := usecase.NewGetDMChannelUsecase(logger, store)
	createDMUc := usecase.NewCreateDMChannelUsecase(logger, store)

	httpServer := http.NewHttpServer(configConfig.Hub, validator, &http.Usecases{
		CreateDMChannel: createDMUc,
		GetDMChannel:    getDMUc,
	})
	router := NewRouter(httpServer)

	return server.NewServer(router, infra), nil

}
