package user

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrEmailRequired      = errors.New("email is required")
	ErrPasswordRequired   = errors.New("password is required")
	ErrPasswordTooShort   = errors.New("password is too short (min 3 chars)")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

// User описывает структуру пользователя в базе данных
type User struct {
	ID           int       `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"` // Никогда не возвращается в JSON
	CreatedAt    time.Time `db:"created_at"`
}

// Validate проверяет бизнес-сущность перед сохранением в БД
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

// --- DTO: Структуры запросов и ответов (Транспортный слой) ---

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate проверяет корректность полей входного запроса
func (c CreateUserRequest) Validate() error {
	if c.Email == "" {
		return ErrEmailRequired
	}
	if !strings.Contains(c.Email, "@") {
		return ErrInvalidEmail
	}
	if c.Password == "" {
		return ErrPasswordRequired
	}
	if len(c.Password) < 3 {
		return ErrPasswordTooShort
	}
	return nil
}

type UserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}
