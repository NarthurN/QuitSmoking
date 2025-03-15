package helpers

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/NarthurN/QuitSmoking/internal/configs"
	"github.com/NarthurN/QuitSmoking/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type Tokener struct{}

func NewTokener() *Tokener {
	return &Tokener{}
}

func (t *Tokener) GetJwtToken(username string) (string, error) {
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

func (t *Tokener) VerifyUser(token string) (*models.Claims, error) {
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

func (t *Tokener) AllowedPath(path string, m map[string]struct{}) bool {
	if _, ok := m[path]; ok || strings.HasPrefix(path, "/static/") {
		return true
	}
	return false
}

func (t *Tokener) CheckPermision(username, path string) bool {
	requiredRoles, ok := configs.PathsRoles[path]
	if ok {
		rolesOfUser, ok := configs.UserRoles[username]
		if !ok {
			return false
		}
		for _, requiredRole := range requiredRoles {
			for _, roleOfUser := range rolesOfUser {
				if requiredRole == roleOfUser {
					return true
				}
			}
		}
	}
	return true
}

func SlogErr(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func SlogDebug(str string) slog.Attr {
	return slog.Attr{
		Key:   "debug",
		Value: slog.StringValue(str),
	}
}
