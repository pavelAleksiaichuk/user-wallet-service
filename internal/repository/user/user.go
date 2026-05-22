package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	userModel "userwalletservice/internal/model/user"
	"userwalletservice/internal/repository"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// 🔥 Возвращает int (id нового пользователя) и error
func (r *Repository) Register(ctx context.Context, user userModel.User, walletRepository repository.Wallet) (int, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	var userID int
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`

	err = tx.QueryRowxContext(ctx, query, user.Email, user.PasswordHash).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("register user insert: %w", err)
	}

	if err := walletRepository.CreateTx(ctx, tx, userID); err != nil {
		return 0, fmt.Errorf("register user wallet create: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit transaction: %w", err)
	}

	return userID, nil // 🔥 Возвращаем сгенерированный базой ID наружу
}

func (r *Repository) GetByID(ctx context.Context, id int) (*userModel.User, error) {
	var user userModel.User
	query := `SELECT id, email, password_hash FROM users WHERE id=$1`

	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, userModel.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*userModel.User, error) {
	var user userModel.User
	query := `SELECT id, email, password_hash FROM users WHERE email=$1`

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, userModel.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return &user, nil
}
