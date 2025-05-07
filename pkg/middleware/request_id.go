package middleware

import (
	"context"
	"finance-backend/pkg/utils"
	"net/http"

	"github.com/google/uuid"
)

func RequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-Id")

		if requestID == "" {
			requestID = uuid.New().String()
			r.Header.Set("X-Request-Id", requestID)
		}
		w.Header().Set("X-Request-Id", requestID)
		ctx := context.WithValue(r.Context(), utils.ContextKeyRequestId, requestID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
