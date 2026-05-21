package user

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrEmailRequired    = errors.New("email is required")
	ErrPasswordRequired = errors.New("password is required")
	ErrPasswordTooShort = errors.New("password is too short (min 3 chars)")
	ErrInvalidEmail     = errors.New("invalid email")
)

type User struct {
	ID           int       `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"` // Никогда не возвращай это в JSON!
	CreatedAt    time.Time `db:"created_at"`
}

// 🔥 Метод валидации теперь живет прямо в структуре
func (u *User) Validate(password string) error {
	if u.Email == "" {
		return ErrEmailRequired
	}

	if !strings.Contains(u.Email, "@") {
		return ErrInvalidEmail
	}

	if password == "" {
		return ErrPasswordRequired
	}

	if len(password) < 3 {
		return ErrPasswordTooShort
	}

	return nil
}
