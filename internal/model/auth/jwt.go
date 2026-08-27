package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// Claims теперь лежит здесь. Слой моделей знает про внешнюю либу jwt/v5 — это нормально,
// так как это стандартные типы (RegisteredClaims), которые кочуют по всему приложению.
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}
