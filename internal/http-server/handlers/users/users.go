package users

import (
	"dolgorukov-dom/internal/lib/api/response"
	"dolgorukov-dom/internal/storage/dto"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type Users interface {
	AddUser(name, email string) (int64, error)
	GetUser(id int64) (*dto.User, error)
	GetUsersList() ([]*dto.User, error)
	UpdateUser(id int64, name, email string) error
	DeleteUser(id int64) error
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
	r.Get("/", uh.GetUsersListHandler())
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

		resUsersList, err := uh.GetUsersList()
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

		user, err := uh.GetUser(int64(id))
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

		rawUser, err := r.GetBody()
		if err != nil {
			log.Info("error getting body")
			render.JSON(w, r, err)
		}

		var req UserRequest

		err = render.DecodeJSON(rawUser, req)
		if err != nil {
			log.Info("error decoding body")
			render.JSON(w, r, err)
		}

		id, err := uh.AddUser(req.Name, req.Email)
		if err != nil {
			log.Info("error adding user")
			render.JSON(w, r, err)
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
		}

		rawUser, err := r.GetBody()
		if err != nil {
			log.Info("error getting body")
			render.JSON(w, r, err)
		}

		var req UserRequest
		err = render.DecodeJSON(rawUser, &req)
		if err != nil {
			log.Info("error decoding body")
			render.JSON(w, r, err)
		}

		err = uh.UpdateUser(int64(id), req.Name, req.Email)
		if err != nil {
			log.Info("error updating user")
			render.JSON(w, r, err)
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

		err = uh.DeleteUser(int64(id))
		if err != nil {
			log.Info("error deleting user")
		}

		render.JSON(w, r, http.StatusOK)
	}
}
