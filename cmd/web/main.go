package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NarthurN/QuitSmoking/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", handlers.Home)
	r.Get("/smokers", handlers.GetSmokers)
	r.Get("/smokers/{id}", handlers.GetSmoker)
	r.Post("/smokers", handlers.PostSmoker)
	r.Delete("/smokers/{id}", handlers.DeleteSmoker)
	r.Put("/smokers/{id}", handlers.PutSmoker)
	r.Get("/smokers/{id}", handlers.GetSmokersDiffTime)

	addr := ":8080"
	log.Printf("Server is listening on %s ...", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		fmt.Printf("Ошибка при запуске сервера %s", err.Error())
	}
}
