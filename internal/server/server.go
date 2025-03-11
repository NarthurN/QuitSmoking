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
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func SetupRoutes(h *handlers.Handlers) http.Handler {
	mv := middleware.New(h.Logger)
	mux := http.NewServeMux()
	mux.Handle(`GET /`, h.Home())
	mux.Handle(`GET /form`, h.GetForm())
	mux.Handle(`POST /signin`, h.Signin())
	mux.Handle("POST /logout", h.Logout())
	mux.Handle(`GET /smokers`, h.GetSmokers())
	mux.Handle(`GET /profile`, h.GetSmokerProfile())
	
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	return mv.Log(mv.JwtAuth(mux))
}

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
