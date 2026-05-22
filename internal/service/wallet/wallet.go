package wallet

import (
	"context"
	"errors"
	walletModel "userwalletservice/internal/model/wallet"
	"userwalletservice/internal/repository"

	"github.com/shopspring/decimal"
)

type Service struct {
	walletRepository repository.Wallet
}

// New — конструктор сервиса кошельков
func New(walletRepository repository.Wallet) *Service {
	return &Service{
		walletRepository: walletRepository,
	}
}

// GetWalletByUserID возвращает кошелек по ID пользователя
func (s *Service) GetWalletByUserID(ctx context.Context, userID int) (*walletModel.Wallet, error) {
	return s.walletRepository.GetByUserID(ctx, userID)
}

// Deposit пополняет баланс
func (s *Service) Deposit(ctx context.Context, userID int, amount decimal.Decimal) (*walletModel.Wallet, error) {
	if !amount.IsPositive() {
		return nil, errors.New("amount must be greater than zero")
	}
	return s.walletRepository.Deposit(ctx, userID, amount)
}

// Withdraw списывает средства
func (s *Service) Withdraw(ctx context.Context, userID int, amount decimal.Decimal) (*walletModel.Wallet, error) {
	if !amount.IsPositive() {
		return nil, errors.New("amount must be greater than zero")
	}
	return s.walletRepository.Withdraw(ctx, userID, amount)
}

// Transfer переводит деньги от одного юзера другому и возвращает новый баланс отправителя
func (s *Service) Transfer(ctx context.Context, fromUserID int, toUserID int, amount decimal.Decimal) (decimal.Decimal, error) {
	if !amount.IsPositive() {
		return decimal.Zero, errors.New("amount must be greater than zero")
	}
	if fromUserID == toUserID {
		return decimal.Zero, errors.New("cannot transfer money to yourself")
	}

	// Ловим структуру *walletModel.Wallet, которую отдаёт обновленный репозиторий
	senderWallet, err := s.walletRepository.Transfer(ctx, fromUserID, toUserID, amount)
	if err != nil {
		return decimal.Zero, err
	}

	// Возвращаем конкретно поле Balance типа decimal.Decimal
	return senderWallet.Balance, nil
}
