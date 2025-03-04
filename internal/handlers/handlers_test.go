package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHomeWhenOk(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	h := New(nil, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Home())
	handler.ServeHTTP(responseRecorder, r)

	assert.Equal(t, responseRecorder.Code, http.StatusOK)

	expectedText := "Это приложение для тех, кто бросает курить!"
	body := responseRecorder.Body.String()
	assert.Equal(t, expectedText, body)
}

func TestGetSmokersWhenOk(t *testing.T) {
	r := httptest.NewRequest("GET", "/smokers", nil)
	h := New(nil, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetSmokers())
	handler.ServeHTTP(responseRecorder, r)

	assert.Equal(t, responseRecorder.Code, http.StatusOK)

}
