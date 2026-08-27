package wallet

import (
	"time"
)

const (
	TypeDeposit  = "deposit"
	TypeWithdraw = "withdraw"
	TypeTransfer = "transfer"
)

type Transaction struct {
	ID               int        `json:"id"`
	SenderWalletID   *int       `json:"sender_wallet_id,omitempty"`   
	ReceiverWalletID *int       `json:"receiver_wallet_id,omitempty"` 
	Amount           float64    `json:"amount"`
	Type             string     `json:"type"` 
	CreatedAt        time.Time  `json:"created_at"`
}
