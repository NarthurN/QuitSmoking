package main

import (
	"log"
	"net/http"

	"github.com/NarthurN/QuitSmoking/internal/handlers"
)

func main() {
	addr := ":8080"
	log.Printf("Server is listening on %s ...", addr)

	http.HandleFunc("/", handlers.Home)

	http.ListenAndServe(addr, nil)
}