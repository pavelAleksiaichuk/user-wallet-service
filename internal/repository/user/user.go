package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	userModel "userwalletservice/internal/model/user"
	walletRepository "userwalletservice/internal/repository/wallet"

	"github.com/jmoiron/sqlx"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// 🔥 Теперь метод принимает интерфейс репозитория кошелька
func (r *Repository) Register(ctx context.Context, user *userModel.User, walletRepository walletRepository.Repository) error {
	// 1. Открываем транзакцию через sqlx
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// 2. Гарантируем откат в случае паники или ошибки
	defer tx.Rollback()

	// 3. Сохраняем пользователя, выполняя запрос СТРОГО через транзакцию tx, а не через r.db
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`
	err = tx.QueryRowxContext(ctx, query, user.Email, user.PasswordHash).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("register user insert: %w", err)
	}

	// 4. Создаем кошелек для пользователя в рамках ЭТОЙ ЖЕ транзакции
	// Передаем tx в новый метод CreateTx
	if err := walletRepository.CreateTx(ctx, tx, user.ID); err != nil {
		return fmt.Errorf("register user wallet create: %w", err)
	}

	// 5. Если всё прошло успешно — фиксируем изменения в БД
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int) (*userModel.User, error) {
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

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*userModel.User, error) {
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
