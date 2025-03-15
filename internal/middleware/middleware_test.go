package middleware

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NarthurN/QuitSmoking/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockVerifier struct {
	mock.Mock
}

func (m *MockVerifier) VerifyUser(token string) (*models.Claims, error) {
	args := m.Called(token)
	return args.Get(0).(*models.Claims), args.Error(1)
}

func (m *MockVerifier) GetJwtToken(username string) (string, error) {
	args := m.Called(username)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockVerifier) AllowedPath(path string, allowedPaths map[string]struct{}) bool {
	args := m.Called(path, allowedPaths)
	return args.Get(0).(bool)
}

func (m *MockVerifier) CheckPermision(username, path string) bool {
	args := m.Called(username, path)
	return args.Get(0).(bool)
}

func TestJwtAuthWhenPathIsAllowedWithRefresh(t *testing.T) {
	expectedToken := "1111"
	allowedPath := "/smoker"
	user := "arthur"
	// Создаем мок Verifier
	mockVerifier := new(MockVerifier)
	// Настраиваем мок, чтобы он возвращал успешный результат
	mockVerifier.On("AllowedPath", allowedPath, mock.Anything).Return(true)
	// mockVerifier.On("VerifyUser", expectedToken).Return(&models.Claims{
	// 	Username: user,
	// 	RegisteredClaims: jwt.RegisteredClaims{
	// 		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(20 * time.Second)),
	// 	},
	// }, nil)
	mockVerifier.On("GetJwtToken", user).Return(expectedToken, nil)
	//mockVerifier.On("CheckPermision", user, allowedPath).Return(true)

	// Создаем middleware с моком Verifier
	middleware := New(slog.Default(), mockVerifier)

	// Тестовый обработчик
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Тестовый запрос с валидным токеном
	req := httptest.NewRequest("GET", allowedPath, nil)
	rr := httptest.NewRecorder()
	http.SetCookie(rr, &http.Cookie{
		Name:    "token",
		Value:   "Bearer " + expectedToken,
		Expires: time.Now().UTC().Add(20 * time.Second),
		Path:    "/",
	})

	// Вызываем middleware
	middleware.JwtAuth(handler).ServeHTTP(rr, req)

	// Проверяем, что статус 200
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем, что VerifyUser был вызван с правильным аргументом
	//mockVerifier.AssertCalled(t, "CheckPermision", user, allowedPath)
	//mockVerifier.AssertCalled(t, "VerifyUser", expectedToken)
	//mockVerifier.AssertCalled(t, "GetJwtToken", user)
}

func TestJwtAuthWhenPathIsAllowed(t *testing.T) {
	expectedToken := "1111"
	allowedPath := "/smoker"
	user := "arthur"
	// Создаем мок Verifier
	mockVerifier := new(MockVerifier)
	// Настраиваем мок, чтобы он возвращал успешный результат
	mockVerifier.On("VerifyUser", expectedToken).Return(&models.Claims{
		Username: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(5 * time.Minute)),
		},
	}, nil)
	mockVerifier.On("AllowedPath", allowedPath, mock.Anything).Return(false)
	// Создаем middleware с моком Verifier
	middleware := New(slog.Default(), mockVerifier)

	// Тестовый обработчик
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Тестовый запрос с валидным токеном
	req := httptest.NewRequest("GET", allowedPath, nil)
	req.Header.Set("Authorization", "Bearer "+expectedToken)
	rr := httptest.NewRecorder()

	// Вызываем middleware
	middleware.JwtAuth(handler).ServeHTTP(rr, req)

	// Проверяем, что статус 200
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем, что VerifyUser был вызван с правильным аргументом
	mockVerifier.AssertCalled(t, "VerifyUser", expectedToken)
}

func TestJwtAuthWhenPathIsNotAllowed(t *testing.T) {
	expectedToken := "1111"
	notAllowedPath := "/static/abc"
	user := "arthur"

	// Создаем мок Verifier
	mockVerifier := new(MockVerifier)
	// Настраиваем мок, чтобы он возвращал успешный результат
	mockVerifier.On("AllowedPath", notAllowedPath, mock.Anything).Return(true)
	// Настраиваем мок, чтобы он возвращал успешный результат
	mockVerifier.On("VerifyUser", expectedToken).Return(&models.Claims{
		Username: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(5 * time.Minute)),
		},
	}, nil)

	// Создаем middleware с моком Verifier
	middleware := New(slog.Default(), mockVerifier)

	// Тестовый обработчик
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Тестовый запрос с валидным токеном
	req := httptest.NewRequest("GET", notAllowedPath, nil)
	req.Header.Set("Authorization", "Bearer "+expectedToken)
	rr := httptest.NewRecorder()

	// Вызываем middleware
	middleware.JwtAuth(handler).ServeHTTP(rr, req)

	// Проверяем, что статус 200
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем, что VerifyUser был вызван с правильным аргументом
	mockVerifier.AssertCalled(t, "AllowedPath", notAllowedPath, mock.Anything)
}
