package wallet

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	contextModel "userwalletservice/internal/model/context"
	walletModel "userwalletservice/internal/model/wallet" // Импортируем наши модели
	walletService "userwalletservice/internal/service/wallet"
)

type Controller struct {
	walletService *walletService.Service
}

func New(walletService *walletService.Service) *Controller {
	return &Controller{
		walletService: walletService,
	}
}

const timeLayout = "2006-01-02 15:04:05"

func (c *Controller) GetBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Сюда пускает только AuthMiddleware. Если ID нет — это баг нашей разработки (500)
	userID, ok := r.Context().Value(contextModel.UserIDKey).(int)
	if !ok {
		log.Printf("[ERROR] AuthMiddleware failed to set UserIDKey for GetBalance")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	wallet, err := c.walletService.GetWalletByUserID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, walletModel.ErrWalletNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		log.Printf("failed to get wallet balance: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	response := walletModel.WalletResponse{
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

	userID, ok := r.Context().Value(contextModel.UserIDKey).(int)
	if !ok {
		log.Printf("[ERROR] AuthMiddleware failed to set UserIDKey for Deposit")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	var req walletModel.DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	wallet, err := c.walletService.Deposit(r.Context(), userID, req.Amount)
	if err != nil {
		// Проверяем через errors.Is конкретную доменную ошибку
		if errors.Is(err, walletModel.ErrAmountMustBePositive) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		log.Printf("failed to deposit funds: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	response := walletModel.WalletResponse{
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
		log.Printf("[ERROR] AuthMiddleware failed to set UserIDKey for Withdraw")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	var req walletModel.WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	wallet, err := c.walletService.Withdraw(r.Context(), userID, req.Amount)
	if err != nil {
		// Заменяем строковое сравнение на errors.Is
		if errors.Is(err, walletModel.ErrAmountMustBePositive) ||
			errors.Is(err, walletModel.ErrInsufficientFunds) ||
			errors.Is(err, walletModel.ErrWalletNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		log.Printf("failed to withdraw funds: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	response := walletModel.WalletResponse{
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
		log.Printf("[ERROR] AuthMiddleware failed to set UserIDKey for Transfer")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	var req walletModel.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	// Фиксим компилятор: ловим и новый баланс, и ошибку
	newBalance, err := c.walletService.Transfer(r.Context(), fromUserID, req.ToUserID, req.Amount)
	if err != nil {
		// Ошибки перевода (себе нельзя, не хватает денег, сумма <= 0) — это 400 Bad Request
		if errors.Is(err, walletModel.ErrAmountMustBePositive) ||
			errors.Is(err, walletModel.ErrInsufficientFunds) ||
			errors.Is(err, walletModel.ErrTransferToSameUser) ||
			errors.Is(err, walletModel.ErrWalletNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		log.Printf("failed to transfer funds from user %v to %v: %v", fromUserID, req.ToUserID, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
		return
	}

	// Создаем структуру для красивого ответа фронтенду
	// Можно объявить её прямо тут (анонимная структура), чтобы не засорять глобальную область
	response := struct {
		Status           string `json:"status"`
		SenderNewBalance string `json:"sender_new_balance"`
	}{
		Status:           "SUCCESS",
		SenderNewBalance: newBalance.String(), // Переводим decimal.Decimal в строку
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
