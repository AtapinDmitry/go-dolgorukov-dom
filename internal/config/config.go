package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env        string `mapstructure:"env"`
	DB         `mapstructure:"db"`
	HTTPServer `mapstructure:"http_server"`
}

type DB struct {
	Host     string `mapstructure:"postgres_host"`
	Port     string `mapstructure:"postgres_port"`
	User     string `mapstructure:"postgres_user"`
	Password string `mapstructure:"postgres_password"`
	DBName   string `mapstructure:"db_name"`
}

type HTTPServer struct {
	Address     string        `mapstructure:"address"`
	Timeout     time.Duration `mapstructure:"timeout"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`
}

func MustLoad(configPath string) *Config {
	if configPath == "" {
		viper.AddConfigPath("./config")
	} else {
		viper.AddConfigPath(configPath)
	}

	viper.SetConfigName("local")
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			panic(fmt.Errorf("config file not found: %s", err))
		} else {
			panic(fmt.Errorf("config file error: %s", err))
		}
	}

	var cfg Config

	err := viper.Unmarshal(&cfg)
	if err != nil {
		panic(fmt.Errorf("error reading config file: %s", err))
	}

	return &cfg
}
