package server

import (
	"net/http"
	"time"

	"github.com/NarthurN/QuitSmoking/internal/handlers"
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

func SetupRoutes(h *handlers.Handlers) *http.ServeMux{
	mux := http.NewServeMux()
	mux.Handle(`GET /`, h.Home())
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
