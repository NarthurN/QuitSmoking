package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/NarthurN/QuitSmoking/internal/configs"
	"github.com/NarthurN/QuitSmoking/internal/helpers"
	"github.com/NarthurN/QuitSmoking/internal/mocks"
	"github.com/NarthurN/QuitSmoking/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

type Handlers struct {
	db     *sql.DB
	Logger *slog.Logger
}

func New(db *sql.DB, logger *slog.Logger) *Handlers {
	return &Handlers{
		db:     db,
		Logger: logger,
	}
}

// Home отображает стартовую страницу
func (h *Handlers) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value(models.ContextString("username")).(string)
		w.WriteHeader(http.StatusOK)
		if !ok {
			w.Write([]byte("Привет, Гость! Это приложение для тех, кто бросает курить!"))
			return
		} 
		w.Write(fmt.Appendf(nil, "Привет, %s! Это приложение для тех, кто бросает курить!", username))
	}
}

// Signin выдаёт JWT-токен по user и password
func (h *Handlers) Signin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds models.Credentials

		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		tokenString, err := helpers.GetJwtToken(&creds)
		if err != nil {
			h.Logger.Debug(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		smoker, ok := mocks.Smokers[creds.Username]
		if !ok {
			http.Error(w, "Пользователя с таким username не существует", http.StatusBadRequest)
			return
		}
	
		expectedPassword := creds.Password
		if expectedPassword != smoker.Password {
			http.Error(w, "Пароль неверный", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Authorization", "Bearer "+tokenString)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Вы авторизированы!"))
	}
}

func (h *Handlers) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		tknStr := c.Value

		claims := &models.Claims{}

		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(toke *jwt.Token) (any, error) {
			return []byte(configs.JwtKey), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if !tkn.Valid {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		expirationTime := time.Now().UTC().Add(5 * time.Minute)
		claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(configs.JwtKey))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})
	}
}

func (h *Handlers) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Expires: time.Now().UTC(),
		})
	}
}

// GetSmokers отображает всех Smokers в формате JSON
func (h *Handlers) GetSmokers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		smokers, err := json.Marshal(&mocks.Smokers)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(smokers)
	}
}

// GetSmoker отображает данные одного Smoker по его id
func GetSmoker(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	smoker, ok := mocks.Smokers[id]
	if !ok {
		http.Error(w, "Такого курильщика у нас нет", http.StatusBadRequest)
		return
	}

	smokerBytes, err := json.Marshal(&smoker)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(smokerBytes)
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
