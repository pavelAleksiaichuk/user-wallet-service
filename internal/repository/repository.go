package repository

import (
	"context"
	"userwalletservice/internal/model/wallet"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type Wallet interface {
	Create(ctx context.Context, userID int) error
	CreateTx(ctx context.Context, tx *sqlx.Tx, userID int) error // 🔥 ДОБАВИЛИ В ИНТЕРФЕЙС
	GetByUserID(ctx context.Context, userID int) (*wallet.Wallet, error)
	Deposit(ctx context.Context, userID int, amount decimal.Decimal) (*wallet.Wallet, error)
	Withdraw(ctx context.Context, userID int, amount decimal.Decimal) (*wallet.Wallet, error)
	Transfer(ctx context.Context, fromUserID int, toUserID int, amount decimal.Decimal) error
}
