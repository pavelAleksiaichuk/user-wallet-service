package wallet

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID        int             `json:"id" db:"id"`
	UserID    int             `json:"user_id" db:"user_id"`
	Balance   decimal.Decimal `json:"balance" db:"balance"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

// WalletResponse используется для отдачи красивого JSON фронтенду
type WalletResponse struct {
	ID        int             `json:"id"`
	UserID    int             `json:"user_id,omitempty"` // omitempty, если где-то скрываем
	Balance   decimal.Decimal `json:"balance"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}

// DepositRequest DTO для пополнения
type DepositRequest struct {
	Amount decimal.Decimal `json:"amount"`
}

type WithdrawRequest struct {
	Amount decimal.Decimal `json:"amount"`
}

type TransferRequest struct {
	ToUserID int             `json:"to_user_id"`
	Amount   decimal.Decimal `json:"amount"`
}

// 🔥 Выносим все бизнес-ошибки кошелька в модель
var (
	ErrAmountMustBePositive = errors.New("amount must be greater than zero")
	ErrInsufficientFunds    = errors.New("insufficient funds")
	ErrWalletNotFound       = errors.New("wallet not found")
	ErrTransferToSameUser   = errors.New("cannot transfer money to yourself")
)
