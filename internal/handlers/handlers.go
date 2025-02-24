package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/NarthurN/QuitSmoking/internal/mocks"
	"github.com/NarthurN/QuitSmoking/internal/models"
	"github.com/go-chi/chi/v5"
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

func GetSmoker(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	smoker, ok := mocks.Smokers[id]
	if !ok {
		http.Error(w, "Такого курильщика у нас нет", http.StatusBadRequest)
		return
	}

	smokerBytes, err := json.Marshal(&smoker)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(smokerBytes)
}

func PostSmoker(w http.ResponseWriter, r *http.Request) {
	var smoker models.Smoker
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &smoker); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := mocks.Smokers[smoker.ID]; ok {
		http.Error(w, "Такой курильщик уже существует", http.StatusBadRequest)
		return
	}

	mocks.Smokers[smoker.ID] = &smoker

	message := map[string]string{"message": "Пользователь записан", "id": smoker.ID}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func DeleteSmoker(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, ok := mocks.Smokers[id]; !ok {
		http.Error(w, "Такого курильщика не существует", http.StatusBadRequest)
		return
	}

	delete(mocks.Smokers, id)

	message := map[string]string{"message": "Пользователь удалён", "id": id}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func PutSmoker(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, ok := mocks.Smokers[id]; !ok {
		http.Error(w, "Такого курильщика не существует", http.StatusBadRequest)
		return
	}

	var smoker models.Smoker
	var buf bytes.Buffer

	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &smoker); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mocks.Smokers[id].ID = smoker.ID
	mocks.Smokers[id].Name = smoker.Name
	mocks.Smokers[id].Experience = smoker.Experience
	mocks.Smokers[id].StoppedSmoking = smoker.StoppedSmoking

	message := map[string]string{"message": "Данные пользователя изменены", "id": id}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func GetSmokersDiffTime(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, ok := mocks.Smokers[id]; !ok {
		http.Error(w, "Такого курильщика не существует", http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()

	diff := now.Sub(mocks.Smokers[id].StoppedSmoking)

	// Преобразуем разницу в годы, месяцы, дни и часы
	years := int(diff.Hours() / 24 / 365)
	remaining := diff - time.Duration(years)*365*24*time.Hour

	months := int(remaining.Hours() / 24 / 30)
	remaining -= time.Duration(months) * 30 * 24 * time.Hour

	days := int(remaining.Hours() / 24)
	remaining -= time.Duration(days) * 24 * time.Hour

	hours := int(remaining.Hours())

	timePassed := fmt.Sprintf("%d лет, %d месяцев, %d дней, %d часов", years, months, days, hours)

	message := map[string]string{"message": "Вы не курили", "id": id, "stoppedSmoking":timePassed}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}