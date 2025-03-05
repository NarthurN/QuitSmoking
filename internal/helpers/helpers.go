package helpers

import (
	"fmt"
	"time"

	"github.com/NarthurN/QuitSmoking/internal/configs"
	"github.com/NarthurN/QuitSmoking/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func GetJwtToken(username string) (string, error) {
	op := "helpers.GetJwtToken"
	expirationTime := time.Now().UTC().Add(5 * time.Minute)

	claims := &models.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(configs.JwtKey))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return tokenString, nil
}

func VerifyUser(token string) (*models.Claims, error) {
	op := "helpers.VerifyUser"
	claims := &models.Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(configs.JwtKey), nil
	})
	if err != nil {
		return claims, fmt.Errorf("%s: %w", op, err)
	}

	if !tkn.Valid {
		return claims, fmt.Errorf("%s: %w", op, err)
	}

	return claims, nil
}
