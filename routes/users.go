package routes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"main/middlewares"
	"main/models"
	"main/services"
	"main/utils/logging"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/redis/go-redis/v9"
	"gopkg.in/validator.v2"
)

type UsersHandler struct {
	Parent      *Router
	Router      *httprouter.Router
	Services    *services.Services
	middlewares *middlewares.Middlewares
	logger      *logging.Logger
	redis       *redis.Client
}

func NewUsersHandler(router *Router) *UsersHandler {
	return &UsersHandler{
		Parent:      router,
		Router:      router.Router,
		Services:    router.Services,
		middlewares: router.middlewares,
		logger:      router.logger,
		redis:       router.redis,
	}
}

func (h UsersHandler) RegisterUsersRoutes() {
	h.Router.HandlerFunc(http.MethodGet, "/users/current", h.middlewares.ApplyMiddlewares(
		h.GetCurrent,
		h.middlewares.ForAuth,
	))
	h.Router.HandlerFunc(http.MethodPost, "/users/register", h.middlewares.ApplyMiddlewares(
		h.RegisterUser,
		h.middlewares.ForUnauth,
	))
	h.Router.HandlerFunc(http.MethodPost, "/users/login", h.middlewares.ApplyMiddlewares(
		h.LoginUser,
		h.middlewares.ForUnauth,
	))
}

func (h UsersHandler) GetCurrent(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		if err == http.ErrNoCookie {
			h.Parent.error(w, "unauthorized", http.StatusUnauthorized)
		} else {
			h.Parent.error(w, fmt.Sprintf("bad request: %s", err.Error()), http.StatusBadRequest)
		}
		return
	}
	sessionToken := cookie.Value

	if sessionToken == "" {
		h.Parent.error(w, "unauthorized", http.StatusInternalServerError)
		return
	}

	h.logger.Info(fmt.Sprintf("sessionToken: %s", sessionToken))

	result, err := h.redis.Get(context.Background(), sessionToken).Result()
	if err != nil && err.Error() != "redis: nil" {
		h.Parent.error(w, fmt.Sprintf("bad request: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if result == "" {
		http.SetCookie(w, &http.Cookie{
			Name:    "sessionID",
			Value:   "",
			Expires: time.Now().Add(365 * 24 * time.Hour),
			Path:    "/",
		})

		h.logger.Info(fmt.Sprintf("user: %s", result))

		h.Parent.error(w, "unauthorized", http.StatusInternalServerError)
		return
	}

	var user models.User

	err = json.Unmarshal([]byte(result), &user)
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("can not unmarshal user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	jsonResp, _ := json.Marshal(user)
	h.Parent.send(w, string(jsonResp), http.StatusOK)
}

func (h UsersHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var createUserDTO models.CreateUserDTO
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&createUserDTO)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			h.Parent.error(w, fmt.Sprintf("bad request: wrong type provided for field - %s", unmarshalErr.Field), http.StatusBadRequest)
		} else {
			h.Parent.error(w, fmt.Sprintf("bad request: %s", err.Error()), http.StatusBadRequest)
		}
		return
	}

	if err := validator.Validate(createUserDTO); err != nil {
		h.Parent.error(w, fmt.Sprintf("validataion error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	isUserExist, err := h.Services.Users.FindUserByEmail(context.Background(), createUserDTO.Email)
	if err != nil && isUserExist.ID != "" {
		h.Parent.error(w, fmt.Sprintf("can not check is user exist: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if isUserExist.ID != "" {
		h.Parent.error(w, "user with this email already exist", http.StatusBadRequest)
		return
	}

	buildedUser, err := createUserDTO.BuildUser()
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("can not build user: %s", err.Error()), http.StatusBadRequest)
		return
	}

	oid, err := h.Services.Users.CreateUser(context.Background(), buildedUser)
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("can not create user: %s", err.Error()), http.StatusBadRequest)
		return
	}

	h.Parent.send(w, fmt.Sprintf("\"%s\"", oid), http.StatusOK)
}

func (h UsersHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var LoginUserDTO models.LoginUserDTO
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&LoginUserDTO)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			h.Parent.error(w, fmt.Sprintf("bad request: wrong type provided for field - %s", unmarshalErr.Field), http.StatusBadRequest)
		} else {
			h.Parent.error(w, fmt.Sprintf("bad request: %s", err.Error()), http.StatusBadRequest)
		}
		return
	}

	if err := validator.Validate(LoginUserDTO); err != nil {
		h.Parent.error(w, fmt.Sprintf("validataion error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	user, err := h.Services.Users.FindUserByEmail(context.Background(), LoginUserDTO.Email)
	if err != nil && user.ID != "" {
		h.Parent.error(w, "wrong email or password", http.StatusBadRequest)
		return
	}

	match := user.CompareHashAndPassword(LoginUserDTO.Password)
	if !match {
		h.Parent.error(w, "wrong email or password", http.StatusBadRequest)
		return
	}

	sessionToken := uuid.NewString()

	userBytes, err := json.Marshal(user)
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("error then marshal user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	err = h.redis.Set(context.Background(), sessionToken, userBytes, 0).Err()
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("error then create session: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "sessionID",
		Value:   sessionToken,
		Expires: time.Now().Add(365 * 24 * time.Hour),
		Path:    "/",
	})

	h.Parent.send(w, string(userBytes), http.StatusOK)
}
