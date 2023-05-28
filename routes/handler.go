package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"main/middlewares"
	"main/models"
	"main/services"
	"main/utils/logging"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/redis/go-redis/v9"
)

type Router struct {
	Router      *httprouter.Router
	Services    *services.Services
	middlewares *middlewares.Middlewares
	logger      *logging.Logger
	redis       *redis.Client
}

func NewRouter(services *services.Services, rdb *redis.Client, middlewares *middlewares.Middlewares, logger *logging.Logger) *Router {
	router := httprouter.New()

	return &Router{
		Router:      router,
		Services:    services,
		middlewares: middlewares,
		logger:      logger,
		redis:       rdb,
	}
}

func (r *Router) Register() {
	userHandler := NewUsersHandler(r)
	taskHandler := NewTasksHandler(r)

	userHandler.RegisterUsersRoutes()
	taskHandler.RegisterTasksRoutes()
}

func (router *Router) getUser(r *http.Request) (u *models.User, err error) {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		return u, err
	}
	sessionToken := cookie.Value

	if sessionToken == "" {
		return u, err
	}

	router.logger.Info(fmt.Sprintf("sessionToken: %s", sessionToken))

	result, err := router.redis.Get(context.Background(), sessionToken).Result()
	if err != nil && err.Error() != "redis: nil" {
		return u, err
	}

	router.logger.Info(fmt.Sprintf("user: %s", result))

	if result == "" {
		return u, err
	}

	err = json.Unmarshal([]byte(result), &u)

	return u, err
}

func (router *Router) send(w http.ResponseWriter, result string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	w.Write([]byte(fmt.Sprintf(`{"success": true, "result": %s}`, result)))
}

func (router *Router) error(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	w.Write([]byte(fmt.Sprintf(`{"success": false, "error": %s}`, message)))
}
