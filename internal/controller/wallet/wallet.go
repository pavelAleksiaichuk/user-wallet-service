package wallet

import (
	"encoding/json"
	"net/http"
	contextModel "userwalletservice/internal/model/context"
	walletServ "userwalletservice/internal/service/wallet"
)

type Controller struct {
	// Нам нужен доступ к сервису пользователей, чтобы вызывать методы кошелька
	walletService *walletServ.Service
}

type WithdrawRequest struct {
	Amount float64 `json:"amount"`
}

type TransferRequest struct {
	ToUserID int     `json:"to_user_id"`
	Amount   float64 `json:"amount"`
}

func New(walletService *walletServ.Service) *Controller {
	return &Controller{
		walletService: walletService,
	}
}

// GetBalance обрабатывает запрос на получение баланса кошелька
func (c *Controller) GetBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Извлекаем userID, который наш мидлвар (AuthMiddleware) бережно достал из JWT
	// и сохранил внутри контекста запроса.
	userID, ok := r.Context().Value(contextModel.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	// 2. Внимание! Нам нужен метод в сервисе, который умеет получать кошелек по userID.
	// Прямо СЕЙЧАС у нас такого метода в walletService еще нет. Мы напишем его следующим шагом!
	wallet, err := c.walletService.GetWalletByUserID(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	timeLayout := "2006-01-02 15:04:05"

	// 3. Формируем красивый ответ для фронтенда
	response := WalletResponse{
		ID:        wallet.ID,
		UserID:    wallet.UserID,
		Balance:   wallet.Balance,
		CreatedAt: wallet.CreatedAt.Format(timeLayout),
		UpdatedAt: wallet.UpdatedAt.Format(timeLayout),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *Controller) Deposit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Достаем userID из контекста (наш проверенный ключ)
	userID, ok := r.Context().Value(contextModel.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	// 2. Парсим сумму из тела запроса
	var req DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	// 3. Вызываем сервис
	wallet, err := c.walletService.Deposit(r.Context(), userID, req.Amount) // Сделай вызов под свою структуру сервиса
	if err != nil {
		// Если ошибка валидации — отдаем 400, если базы — 500
		if err.Error() == "amount must be greater than zero" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// 4. Форматируем красивый ответ
	timeLayout := "2006-01-02 15:04:05"
	response := WalletResponse{
		ID:        wallet.ID,
		Balance:   wallet.Balance,
		CreatedAt: wallet.CreatedAt.Format(timeLayout),
		UpdatedAt: wallet.UpdatedAt.Format(timeLayout),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *Controller) Withdraw(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value(contextModel.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	var req WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	wallet, err := c.walletService.Withdraw(r.Context(), userID, req.Amount)
	if err != nil {
		if err.Error() == "amount must be greater than zero" || err.Error() == "insufficient funds or wallet not found" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	timeLayout := "2006-01-02 15:04:05"
	response := WalletResponse{
		ID:        wallet.ID,
		UserID:    wallet.UserID,
		Balance:   wallet.Balance,
		CreatedAt: wallet.CreatedAt.Format(timeLayout),
		UpdatedAt: wallet.UpdatedAt.Format(timeLayout),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (c *Controller) Transfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fromUserID, ok := r.Context().Value(contextModel.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	err := c.walletService.Transfer(r.Context(), fromUserID, req.ToUserID, req.Amount)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "transfer completed successfully"})
}
