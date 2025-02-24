package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NarthurN/QuitSmoking/internal/mocks"
)

func Home(w http.ResponseWriter, r *http.Request) {
	msg := "Это приложение для тех, кто бросает курить!"
	w.Write([]byte(msg))
}

func GetSmokers(w http.ResponseWriter, r *http.Request) {
	smokers, err := json.Marshal(&mocks.Smokers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(smokers)
}