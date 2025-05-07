package middleware

import (
	"finance-backend/pkg/logger"
	"net/http"
)

func RecoverMiddleware(next http.Handler) http.Handler {
	log := logger.NewLogger()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				log.Info(r.Context(), "Recovered from panic", map[string]interface{}{"error": err})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
