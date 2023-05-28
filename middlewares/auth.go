package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func (m Middlewares) ForAuth(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		if err == http.ErrNoCookie {
			m.error(w, "unauthorized", http.StatusForbidden)
		} else {
			m.error(w, fmt.Sprintf("bad request: %s", err.Error()), http.StatusForbidden)
		}
		return false
	}
	sessionToken := cookie.Value

	m.logger.Info(fmt.Sprintf("sessionToken: %s", sessionToken))

	if sessionToken == "" {
		m.error(w, "unauthorized", http.StatusForbidden)
		return false
	}

	result, err := m.redis.Get(context.Background(), sessionToken).Result()
	if err != nil && err.Error() != "redis: nil" {
		m.error(w, fmt.Sprintf("bad request: %s", err.Error()), http.StatusForbidden)
		return false
	}

	if result == "" {
		http.SetCookie(w, &http.Cookie{
			Name:    "sessionID",
			Value:   "",
			Expires: time.Now().Add(365 * 24 * time.Hour),
		})

		m.logger.Info(fmt.Sprintf("user: %s", result))

		m.error(w, "unauthorized", http.StatusForbidden)
		return false
	}

	return true
}

func (m Middlewares) ForUnauth(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("sessionID")
	if err != nil && err != http.ErrNoCookie {
		m.error(w, fmt.Sprintf("bad request: %s", err.Error()), http.StatusForbidden)
		return false
	} else if err != nil {
		return true
	}

	sessionToken := cookie.Value
	if sessionToken == "" {
		return true
	}

	m.logger.Info(fmt.Sprintf("sessionToken: %s", sessionToken))

	result, err := m.redis.Get(context.Background(), sessionToken).Result()
	if result == "" || (err != nil && err.Error() == "redis: nil") {
		return true
	}

	m.logger.Info(fmt.Sprintf("user: %s", result))

	m.error(w, "already auth", http.StatusForbidden)
	return false
}

func (m Middlewares) error(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	w.Write([]byte(fmt.Sprintf(`{"success": false, "error": %s}`, message)))
}
