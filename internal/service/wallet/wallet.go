package wallet

import (
	"context"
	"errors"
	walletModel "userwalletservice/internal/model/wallet"
	walletRepository "userwalletservice/internal/repository/wallet"
)

type Service struct {
	walletRepository walletRepository.Repository
}

// New — конструктор сервиса кошельков
func New(walletRepository walletRepository.Repository) *Service {
	return &Service{
		walletRepository: walletRepository,
	}
}

// GetWalletByUserID возвращает кошелек по ID пользователя
func (s *Service) GetWalletByUserID(ctx context.Context, userID int) (*walletModel.Wallet, error) {
	return s.walletRepository.GetByUserID(ctx, userID)
}

// Deposit пополняет баланс
func (s *Service) Deposit(ctx context.Context, userID int, amount float64) (*walletModel.Wallet, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	return s.walletRepository.Deposit(ctx, userID, amount)
}

// Withdraw списывает средства
func (s *Service) Withdraw(ctx context.Context, userID int, amount float64) (*walletModel.Wallet, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	return s.walletRepository.Withdraw(ctx, userID, amount)
}

// Transfer переводит деньги от одного юзера другому
func (s *Service) Transfer(ctx context.Context, fromUserID int, toUserID int, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if fromUserID == toUserID {
		return errors.New("cannot transfer money to yourself")
	}
	return s.walletRepository.Transfer(ctx, fromUserID, toUserID, amount)
}
