package main

import (
	"dolgorukov-dom/internal/config"
	"dolgorukov-dom/internal/lib/logger/sl"
	"dolgorukov-dom/internal/storage/postgres"
	"fmt"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	var configPath string

	rootCmd := cobra.Command{
		Use:     "rest-service-example",
		Version: "v1.0",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("configPath == %s", configPath)
		},
	}

	rootCmd.Flags().StringVarP(&configPath, "config", "c", "",
		"Config file path")

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}

	// init config
	cfg := config.MustLoad(configPath)

	// init logger
	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env)) // к каждому сообщению будет добавляться поле с информацией о текущем окружении

	// init storage
	storage, err := postgres.New(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
	}

	defer func() {
		err = storage.Close()
		if err != nil {
			log.Error("failed to close storage", sl.Err(err))
		}
	}()

	//
	log.Info("initializing server", slog.String("address", cfg.Address)) // Помимо сообщения выведем параметр с адресом
	log.Debug("logger debug mode enabled")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
