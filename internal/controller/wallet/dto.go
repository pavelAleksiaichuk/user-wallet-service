package wallet

type DepositRequest struct {
	Amount float64 `json:"amount"`
}

// WalletResponse описывает то, что мы отдадим фронтенду в формате JSON
type WalletResponse struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	Balance   float64 `json:"balance"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
