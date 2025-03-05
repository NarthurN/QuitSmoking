package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Smoker struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Username       string    `json:"username"`
	Password       string    `json:"password"`
	StoppedSmoking time.Time `json:"stoppedSmoking"`
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Чтобы в контекст передавать не тип string
type ContextString string