package users

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/AtapinDmitry/go-dolgorukov-dom/internal/http-server/handlers"
	"github.com/AtapinDmitry/go-dolgorukov-dom/internal/lib/api/response"
	"github.com/AtapinDmitry/go-dolgorukov-dom/internal/storage/dto"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Users interface {
	AddUser(name, email string) (uint, error)
	GetUser(id uint) (*dto.User, error)
	GetUsersList(filter *dto.UsersListFilter) ([]*dto.User, error)
	UpdateUser(id uint, name, email string) error
	DeleteUser(id uint) error
}

type Handler struct {
	Log *slog.Logger
	Users
}

type UserRequest struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserResponse struct {
	response.Response
	ID uint64 `json:"id"`
}

func UserRoutes(uh *Handler) chi.Router {
	r := chi.NewRouter()
	r.Get("/{page}/{pageSize}", uh.GetUsersListHandler())
	r.Post("/", uh.AddUserHandler())
	r.Get("/{id}", uh.GetUserHandler())
	r.Put("/{id}", uh.UpdateUserHandler())
	r.Delete("/{id}", uh.DeleteUserHandler())
	return r
}

func (uh *Handler) GetUsersListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Users.GetUsersList"

		log := uh.Log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		page, err := strconv.Atoi(chi.URLParam(r, "page"))
		if err != nil {
			log.Info("error getting page number")
			render.JSON(w, r, err)
		}

		pageSize, err := strconv.Atoi(chi.URLParam(r, "pageSize"))
		if err != nil {
			log.Info("error getting page size")
			render.JSON(w, r, err)
		}

		filter := &dto.UsersListFilter{
			Page:     page,
			PageSize: pageSize,
		}

		resUsersList, err := uh.GetUsersList(filter)
		if err != nil {
			log.Info("error getting users list")

			render.JSON(w, r, err)
		}

		render.JSON(w, r, resUsersList)
	}
}

func (uh *Handler) GetUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Users.GetUser"

		log := uh.Log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Info("error getting user id")
			render.JSON(w, r, err)
		}

		user, err := uh.GetUser(uint(id))
		if err != nil {
			log.Info("error getting user")
			render.JSON(w, r, err)
		}

		render.JSON(w, r, user)
	}
}

func (uh *Handler) AddUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Users.AddUser"

		log := uh.Log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req UserRequest

		if err := handlers.DecodeJSONBody(r, &req); err != nil {
			log.Info("error parsing body")
			render.JSON(w, r, err)

			return
		}

		id, err := uh.AddUser(req.Name, req.Email)
		if err != nil {
			log.Info("error adding user")
			render.JSON(w, r, err)

			return
		}

		render.JSON(w, r, id)
	}
}

func (uh *Handler) UpdateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Users.UpdateUser"

		log := uh.Log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Info("error getting user id")
			render.JSON(w, r, err)

			return
		}

		var req UserRequest

		if err := handlers.DecodeJSONBody(r, &req); err != nil {
			log.Info("error parsing body")
			render.JSON(w, r, err)

			return
		}

		err = uh.UpdateUser(uint(id), req.Name, req.Email)
		if err != nil {
			log.Info("error updating user")
			render.JSON(w, r, err)

			return
		}

		render.JSON(w, r, http.StatusOK)
	}
}

func (uh *Handler) DeleteUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Users.DeleteUser"

		log := uh.Log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Info("error getting user id")
		}

		err = uh.DeleteUser(uint(id))
		if err != nil {
			log.Info("error deleting user")
		}

		render.JSON(w, r, http.StatusOK)
	}
}
