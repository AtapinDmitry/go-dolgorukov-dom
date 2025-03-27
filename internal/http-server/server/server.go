package server

import (
	"log/slog"
	"net/http"

	"github.com/AtapinDmitry/go-dolgorukov-dom/internal/config"
	"github.com/AtapinDmitry/go-dolgorukov-dom/internal/storage/postgres"
)

func New(log *slog.Logger, cfg *config.Config, storage *postgres.Storage) *http.Server {
	router := NewRouter(log, storage)

	server := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	return server
}
