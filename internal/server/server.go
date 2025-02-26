package server

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/NarthurN/QuitSmoking/internal/handlers"
	"github.com/NarthurN/QuitSmoking/internal/middleware"
)

func New(mux http.Handler) *http.Server {
	return &http.Server{
		Addr: ":8080",
		Handler: mux,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 60 * time.Second,
	}
}

func SetupRoutes(h *handlers.Handlers, mv *middleware.Middleware) *http.ServeMux{
	mux := http.NewServeMux()
	mux.Handle(`GET /`, mv.Log(
		h.Home()),
	)
	mux.Handle(`GET /smokers`, h.GetSmokers())
	return mux
}

	// r.Get("/", handlers.Home)
	// r.Get("/smokers", handlers.GetSmokers)
	// r.Get("/smokers/{id}", handlers.GetSmoker)
	// r.Post("/smokers", handlers.PostSmoker)
	// r.Delete("/smokers/{id}", handlers.DeleteSmoker)
	// r.Put("/smokers/{id}", handlers.PutSmoker)
	// r.Get("/smokers/{id}", handlers.GetSmokersDiffTime)

func SetupLogger(level string) *slog.Logger {
	var slogLevel slog.Level

	switch level {
	case "debug":
		slogLevel = slog.LevelDebug
	case "info":
		slogLevel = slog.LevelInfo
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: slogLevel,
    }))

	return logger
}