package infrastructure

import "userwalletservice/internal/model/auth"

// Объявляем контракт для работы с токенами
type TokenManager interface {
	GenerateToken(userID int) (string, error)
	// Добавляем второй метод, который обычно есть в JWT (например, Parse/Validate)
	ValidateToken(token string) (*auth.Claims, error)
}
