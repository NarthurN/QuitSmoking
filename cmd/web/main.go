package main

import (
	"log"

	"github.com/NarthurN/QuitSmoking/internal/handlers"
	"github.com/NarthurN/QuitSmoking/internal/server"
)

func main() {
	logger := server.SetupLogger("debug")

	h := handlers.New(nil, logger)

	mux := server.SetupRoutes(h)

	srv := server.New(mux)


	log.Printf("Server is listening on %s ...", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Ошибка при запуске сервера %s", err.Error())
	}
}
