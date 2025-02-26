package main

import (
	"log"
	"net/http"
	"time"

	"github.com/NarthurN/QuitSmoking/internal/handlers"
)

func main() {

	mux := http.NewServeMux()
	srv := &http.Server{
		Addr: ":8080",
		Handler: mux,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	mux.Handle(`GET /`, handlers.NewHome(nil, nil))
	mux.Handle(`GET /smokers`, handlers.NewGetSmokers())

	// r.Get("/", handlers.Home)
	// r.Get("/smokers", handlers.GetSmokers)
	// r.Get("/smokers/{id}", handlers.GetSmoker)
	// r.Post("/smokers", handlers.PostSmoker)
	// r.Delete("/smokers/{id}", handlers.DeleteSmoker)
	// r.Put("/smokers/{id}", handlers.PutSmoker)
	// r.Get("/smokers/{id}", handlers.GetSmokersDiffTime)

	log.Printf("Server is listening on %s ...", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Ошибка при запуске сервера %s", err.Error())
	}
}
