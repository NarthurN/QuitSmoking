package server

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/NarthurN/QuitSmoking/internal/handlers"
	"github.com/NarthurN/QuitSmoking/internal/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
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
	mux := http.NewServeMux()
	mux.Handle(`GET /`, h.Home())
	mux.Handle(`GET /form`, h.GetForm())
	mux.Handle(`POST /signin`, h.Signin())
	mux.Handle("GET /logout", h.Logout())
	mux.Handle("POST /register", h.Register())
	mux.Handle(`GET /smokers`, h.GetSmokers())
	mux.Handle(`GET /profile`, h.GetSmokerProfile())

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// mux.Handle("GET /static/", http.FileServer(http.Dir("static")))
	return h.Mw.Log(h.Mw.JwtAuth(mux))
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

func SetupDB() *sql.DB {
	connStr := "user=arthur password=arthur dbname=Smokers host=localhost port=5432 sslmode=disable"

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to db")

	return db
}
