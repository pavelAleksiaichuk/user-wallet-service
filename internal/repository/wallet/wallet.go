package wallet

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"userwalletservice/internal/model/wallet"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// CreateTx — реализация для работы внутри транзакции
func (r *Repository) CreateTx(ctx context.Context, tx *sqlx.Tx, userID int) error {
	query := `
		INSERT INTO wallets (user_id, balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`
	now := time.Now()
	_, err := tx.ExecContext(ctx, query, userID, 0.0, now, now)
	return err
}

// Твой метод Create оставим, он полезен для тестов или отдельных вызовов
func (r *Repository) Create(ctx context.Context, userID int) error {
	query := `
		INSERT INTO wallets (user_id, balance, created_at, updated_at)
		VALUES (:user_id, :balance, :created_at, :updated_at)
	`
	now := time.Now()
	data := map[string]interface{}{
		"user_id":    userID,
		"balance":    0.00,
		"created_at": now,
		"updated_at": now,
	}
	_, err := r.db.NamedExecContext(ctx, query, data)
	return err
}

// GetByUserID находит кошелек по ID пользователя
func (r *Repository) GetByUserID(ctx context.Context, userID int) (*wallet.Wallet, error) {
	query := `
		SELECT id, user_id, balance, created_at, updated_at 
		FROM wallets 
		WHERE user_id = $1
	`

	var w wallet.Wallet

	// 🔥 Магия sqlx: метод GetContext сам выполняет SQL
	// и автоматически раскладывает колонки по полям структуры w (благодаря тегам db:"...")
	err := r.db.GetContext(ctx, &w, query, userID)
	if err != nil {
		// sqlx внутри использует стандартную ошибку sql.ErrNoRows, если ничего не найдено
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("wallet not found")
		}
		return nil, err
	}

	return &w, nil
}

func (r *Repository) Deposit(ctx context.Context, userID int, amount decimal.Decimal) (*wallet.Wallet, error) {
	var w wallet.Wallet

	// Обновляем баланс и сразу запрашиваем обновленные данные через RETURNING
	query := `
		UPDATE wallets 
		SET balance = balance + $1, updated_at = NOW() 
		WHERE user_id = $2 
		RETURNING id, user_id, balance, created_at, updated_at`

	err := r.db.GetContext(ctx, &w, query, amount, userID)
	if err != nil {
		return nil, err
	}

	return &w, nil
}

func (r *Repository) Withdraw(ctx context.Context, userID int, amount decimal.Decimal) (*wallet.Wallet, error) {
	var w wallet.Wallet

	// Списываем, только если текущий баланс больше или равен сумме списания
	query := `
		UPDATE wallets 
		SET balance = balance - $1, updated_at = NOW() 
		WHERE user_id = $2 AND balance >= $1
		RETURNING id, user_id, balance, created_at, updated_at`

	err := r.db.GetContext(ctx, &w, query, amount, userID)
	if err != nil {
		// Если условие balance >= $1 не выполнилось, sqlx вернет sql.ErrNoRows,
		// потому что UPDATE не найдет и не обновит ни одной строки.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("insufficient funds or wallet not found")
		}
		return nil, err
	}

	return &w, nil
}

// Меняем возвращаемые значения на (*wallet.Wallet, error)
func (r *Repository) Transfer(ctx context.Context, fromUserID int, toUserID int, amount decimal.Decimal) (*wallet.Wallet, error) {
	// Начало транзакции через sqlx
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Авто-откат при панике или ошибке
	defer tx.Rollback()

	// Переменная, куда зальётся обновленный кошелек отправителя
	var senderWallet wallet.Wallet

	// 1. Списываем деньги у отправителя и СРАЗУ возвращаем всю структуру кошелька
	withdrawQuery := `
		UPDATE wallets 
		SET balance = balance - $1, updated_at = NOW() 
		WHERE user_id = $2 AND balance >= $1
		RETURNING id, user_id, balance, created_at, updated_at`

	// Магия sqlx: tx.GetContext выполнит запрос и сам заполнит senderWallet
	err = tx.GetContext(ctx, &senderWallet, withdrawQuery, amount, fromUserID)
	if err != nil {
		// Если денег не хватило или юзер не найден, sqlx вернет sql.ErrNoRows
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("insufficient funds or sender wallet not found")
		}
		return nil, err
	}

	// 2. Начисляем деньги получателю
	depositQuery := `
		UPDATE wallets 
		SET balance = balance + $1, updated_at = NOW() 
		WHERE user_id = $2`

	res, err := tx.ExecContext(ctx, depositQuery, amount, toUserID)
	if err != nil {
		return nil, err
	}

	// Проверяем, существует ли вообще кошелек получателя
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, errors.New("recipient wallet not found")
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Возвращаем указатель на обновленный кошелек отправителя
	return &senderWallet, nil
}
