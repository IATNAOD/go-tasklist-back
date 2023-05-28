package middlewares

import (
	"main/utils/logging"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type Middlewares struct {
	redis  *redis.Client
	logger *logging.Logger
}

func NewMiddlewares(redis *redis.Client, logger *logging.Logger) *Middlewares {
	return &Middlewares{
		redis:  redis,
		logger: logger,
	}
}

func (m Middlewares) ApplyMiddlewares(handler http.HandlerFunc, middlewares ...func(w http.ResponseWriter, r *http.Request) bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := true

		for i := 0; i < len(middlewares) && pass; i++ {
			pass = middlewares[i](w, r)
		}

		if pass {
			handler(w, r)
		}
	}
}
