package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/NarthurN/QuitSmoking/internal/helpers"
	"github.com/NarthurN/QuitSmoking/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

// Для получения статуса ответа
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
		case duration > 500*time.Microsecond:
			m.logger.Warn("slow request", attrs...)
		default:
			m.logger.Info("request completed", attrs...)
		}

		fmt.Println()
	})
}

func (m *Middleware) JwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/signin" || r.URL.Path == "/" {
			next.ServeHTTP(w, r)
			return
		}
		var authHeaderValue string
        
        // Пытаемся получить токен из куки
        cookie, err := r.Cookie("token")
        if err == nil {
            authHeaderValue = cookie.Value
        } else {
            authHeaderValue = r.Header.Get("Authorization")
        }

		if authHeaderValue == "" {
			m.logger.Debug("middleware.jwtAuth.authHeaderValue", helpers.SlogDebug("authHeaderValue is empty"))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeaderValue, " ")
		if len(bearerToken) != 2 {
			m.logger.Debug("middleware.jwtAuth.bearerToken", helpers.SlogDebug("format bearerToken is not Bearer {jwt}"))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		claims, err := helpers.VerifyUser(bearerToken[1])
		if err != nil {
			m.logger.Error("middleware.jwtAuth.VerifyUser", helpers.SlogErr(err))
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if time.Until(claims.ExpiresAt.Time) < 30*time.Second {
			newToken, err := helpers.GetJwtToken(claims.Username)
			if err != nil {
				m.logger.Error("middleware.jwtAuth.GetJwtToken", helpers.SlogErr(err))
			}
			w.Header().Set("Authorization", "Bearer "+newToken)
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, models.ContextString("smoker.name"), claims.Username)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
