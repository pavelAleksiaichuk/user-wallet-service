package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"userwalletservice/internal/model/auth" // 🔥 Импортируем модели
)

type JWTManager struct {
	secretKey []byte
	issuer    string
}

func New(secret string) *JWTManager {
	return &JWTManager{
		secretKey: []byte(secret),
		issuer:    "userwalletservice",
	}
}

func (m *JWTManager) GenerateToken(userID int) (string, error) {
	// Инициализируем Claims из пакета моделей
	claims := &auth.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    m.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

// 🔥 Метод теперь возвращает (*auth.Claims, error) из слоя моделей!
func (m *JWTManager) ValidateToken(tokenString string) (*auth.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrInvalidToken
		}
		return m.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, auth.ErrExpiredToken
		}
		return nil, auth.ErrInvalidToken
	}

	claims, ok := token.Claims.(*auth.Claims)
	if !ok || !token.Valid {
		return nil, auth.ErrInvalidToken
	}

	return claims, nil
}
