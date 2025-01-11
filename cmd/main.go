package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

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

	s, err := ServerBuilder()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init server: %v\n", err)
		os.Exit(1)
	}

	s.Serve()
}
