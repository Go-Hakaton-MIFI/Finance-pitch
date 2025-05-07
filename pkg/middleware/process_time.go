package middleware

import (
	"finance-backend/pkg/logger"
	"net/http"
	"strconv"
	"time"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	start      time.Time
	log        *logger.Logger
	Duration   string
}

func (rw *responseWriterWrapper) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	processTime := time.Since(rw.start)
	rw.Duration = strconv.FormatFloat(processTime.Seconds(), 'f', 6, 64)
	rw.ResponseWriter.Header().Set("X-Process-Time", rw.Duration)
	rw.ResponseWriter.Header().Set("Content-Type", "application/json")
	rw.ResponseWriter.WriteHeader(statusCode)

}

func (rw *responseWriterWrapper) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func LoggingProcessTimeMiddleware(next http.Handler) http.Handler {
	log := logger.NewLogger()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Info(r.Context(), "request_started", map[string]interface{}{
			"path":   r.URL.Path,
			"method": r.Method,
		})

		wrappedWriter := &responseWriterWrapper{ResponseWriter: w, start: start, log: log}

		next.ServeHTTP(wrappedWriter, r)

		log.Info(r.Context(), "request_finished", map[string]interface{}{
			"path":     r.URL.Path,
			"method":   r.Method,
			"duration": wrappedWriter.Duration,
			"status":   wrappedWriter.statusCode,
		})
	})
}
