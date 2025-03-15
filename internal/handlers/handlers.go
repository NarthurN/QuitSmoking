package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"text/template"
	"time"

	"github.com/NarthurN/QuitSmoking/internal/helpers"
	"github.com/NarthurN/QuitSmoking/internal/middleware"
	"github.com/NarthurN/QuitSmoking/internal/mocks"
	"github.com/NarthurN/QuitSmoking/internal/models"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	db     *sql.DB
	Logger *slog.Logger
	Mw     *middleware.Middleware
}

func New(db *sql.DB, logger *slog.Logger) *Handlers {
	return &Handlers{
		db:     db,
		Logger: logger,
		Mw:     middleware.New(logger, helpers.NewTokener()),
	}
}

// Home отображает стартовую страницу
func (h *Handlers) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl, err := template.ParseFiles("static/templates/index.html")
		if err != nil {
			h.Logger.Error("handlers.Home.ParseFIles", helpers.SlogErr(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	}
}

// Signin записывает JWT-токен в заголовок Authorization и проверяет корректность username и password
func (h *Handlers) Signin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var creds models.Credentials
		creds.Username = username
		creds.Password = password

		smoker, ok := mocks.Smokers[creds.Username]
		if !ok {
			h.Logger.Debug("handlers.Signin.CheckSmokerInBase", helpers.SlogDebug("smoker in base not found"))
			http.Error(w, "Пользователя с таким username не существует", http.StatusBadRequest)
			return
		}

		expectedPassword := creds.Password
		if expectedPassword != smoker.Password {
			h.Logger.Debug("handlers.Signin.CheckPasswordOfSmokerInBase", helpers.SlogDebug("password of smoker isn't correct"))
			http.Error(w, "Пароль неверный", http.StatusUnauthorized)
			return
		}

		tokenString, err := h.Mw.Tokener.GetJwtToken(creds.Username)
		if err != nil {
			h.Logger.Error("handlers.Signin.GetJwtToken", helpers.SlogErr(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   "Bearer " + tokenString,
			Expires: time.Now().UTC().Add(5 * time.Minute),
			Path:    "/",
		})

		w.WriteHeader(http.StatusOK)
		tmpl, err := template.ParseFiles("static/templates/signin.html")
		if err != nil {
			h.Logger.Error("handlers.Signin.ParseFIles", helpers.SlogErr(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, smoker)
	}
}

func (h *Handlers) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   "",
			Path:     "/",
			Expires: time.Now(),
			MaxAge:   -1,
		})
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// GetSmokers отображает всех Smokers в формате JSON
func (h *Handlers) GetSmokers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		smokers, err := json.Marshal(&mocks.Smokers)
		if err != nil {
			h.Logger.Error("handlers.GetSmokers.Marshal", helpers.SlogErr(err))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(smokers)
	}
}

func (h *Handlers) GetForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("token")
		if err != nil {
			tmpl, err := template.ParseFiles("static/templates/form.html")
			if err != nil {
				h.Logger.Error("handlers.GetForm.ParseFIles", helpers.SlogErr(err))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, nil)
			return
		}
		http.Redirect(w, r, `/profile`, http.StatusFound)
	}
}

// GetSmokerProfile отображает данные одного Smoker по его id
func (h *Handlers) GetSmokerProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		username, ok := r.Context().Value(models.ContextString("smoker.name")).(string)
		if !ok {
			h.Logger.Error("handlers.GetSmokerProfile.ctxNameToString")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		smoker := mocks.Smokers[username]

		w.WriteHeader(http.StatusOK)
		tmpl, err := template.ParseFiles("static/templates/profile.html")
		if err != nil {
			h.Logger.Error("handlers.GetSmokerProfile.ParseFIles", helpers.SlogErr(err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, smoker)
	}
}

// PostSmoker создаёт нового Smoker
func PostSmoker(w http.ResponseWriter, r *http.Request) {
	var smoker models.Smoker

	if err := json.NewDecoder(r.Body).Decode(&smoker); err != nil {
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
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

// DeleteSmoker удаляет Smoker по id
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
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

// PutSmoker обновляет данные курильщика по id
func PutSmoker(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, ok := mocks.Smokers[id]; !ok {
		http.Error(w, "Такого курильщика не существует", http.StatusBadRequest)
		return
	}

	var smoker models.Smoker

	if err := json.NewDecoder(r.Body).Decode(&smoker); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	mocks.Smokers[id].ID = smoker.ID
	mocks.Smokers[id].Name = smoker.Name
	mocks.Smokers[id].StoppedSmoking = smoker.StoppedSmoking

	message := map[string]string{"message": "Данные пользователя изменены", "id": id}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(message); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

// GetSmokersDiffTime возвращает промежуток времени между StoppedSmoking и времени "сейчас"
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

	message := map[string]string{"message": "Вы не курили", "id": id, "stoppedSmoking": timePassed}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
