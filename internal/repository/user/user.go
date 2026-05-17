package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	userModel "userwalletservice/internal/model/user"

	"github.com/jmoiron/sqlx"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type UserRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Register(ctx context.Context, user *userModel.User) error {
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`

	err := r.db.QueryRowxContext(ctx, query, user.Email, user.PasswordHash).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("register user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*userModel.User, error) {
	var user userModel.User
	query := `SELECT id, email, password_hash FROM users WHERE id=$1`

	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*userModel.User, error) {
	var user userModel.User

	query := `SELECT id, email, password_hash FROM users WHERE email=$1`

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return &user, nil
}
