package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type customResponseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *customResponseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

type Middleware struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}

func (m *Middleware) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()

		m.logger.Debug(
			"Request started",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)

		rw := &customResponseWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		attrs := []any{
			"method", r.Method,
			"status", rw.status,
			"duration", duration.String(),
		}

		switch {
		case rw.status >= 500:
			m.logger.Error("request faild", attrs...)
		case rw.status >= 400:
			m.logger.Warn("request warning", attrs...)
		case duration > 500 * time.Microsecond:
			m.logger.Warn("slow request", attrs...)
		default:
			m.logger.Info("request completed", attrs...)
		}
	})
}