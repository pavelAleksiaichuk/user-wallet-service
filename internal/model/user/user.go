package user

import "time"

type User struct {
	ID           int       `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"` // Никогда не возвращай это в JSON!
	CreatedAt    time.Time `db:"created_at"`
}
