package main

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Welcome string
}

func mustReadViperEnvConfig() Config {
	config, err := readViperEnvConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return config
}

const (
	EnvKeyWelcome = "welcome"
)

func readViperEnvConfig() (Config, error) {
	config := viper.New()
	config.SetEnvPrefix("fcs")
	err := config.BindEnv(EnvKeyWelcome)
	if err != nil {
		return Config{}, err
	}
	return Config{
		Welcome: config.GetString(EnvKeyWelcome),
	}, nil
}
