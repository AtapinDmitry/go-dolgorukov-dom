package server

import (
	"log/slog"
	"net/http"

	"github.com/AtapinDmitry/go-dolgorukov-dom/internal/http-server/handlers/users"
	mwLogger "github.com/AtapinDmitry/go-dolgorukov-dom/internal/http-server/middleware/logger"
	"github.com/AtapinDmitry/go-dolgorukov-dom/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func NewRouter(log *slog.Logger, storage *postgres.Storage) chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, "Hello World!")
	})

	usersHandler := &users.Handler{
		Log:   log,
		Users: storage,
	}

	router.Mount("/users", users.UserRoutes(usersHandler))

	return router
}
