package main

import (
	"log"
	"net/http"

	"github.com/NarthurN/QuitSmoking/internal/handlers"
	"github.com/NarthurN/QuitSmoking/internal/server"
)

func main() {
	h := handlers.New(nil, nil)

	mux := http.NewServeMux()

	server.SetupRoutes(mux, h)

	srv := server.New(mux)


	log.Printf("Server is listening on %s ...", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Ошибка при запуске сервера %s", err.Error())
	}
}
