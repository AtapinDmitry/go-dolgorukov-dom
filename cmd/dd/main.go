package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AtapinDmitry/go-dolgorukov-dom/internal/config"
	"github.com/AtapinDmitry/go-dolgorukov-dom/internal/http-server/server"
	"github.com/AtapinDmitry/go-dolgorukov-dom/internal/lib/logger/sl"
	"github.com/AtapinDmitry/go-dolgorukov-dom/internal/storage/postgres"
	"github.com/spf13/cobra"
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
		Run: func(_ *cobra.Command, _ []string) {
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

	// DB migration
	if err := storage.Migrate(); err != nil {
		log.Error("failed to migrate storage", sl.Err(err))
	}

	//
	log.Info("initializing server", slog.String("address", cfg.Address)) // Помимо сообщения выведем параметр с адресом
	log.Debug("logger debug mode enabled")

	// starting server
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := server.New(log, cfg, storage)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	// TODO: move timeout to config
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	log.Info("server stopped")
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
