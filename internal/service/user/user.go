package user

import (
	"context"
	"errors"
	"strings"
	"userwalletservice/internal/infrastructure/jwt"
	userModel "userwalletservice/internal/model/user"
	"userwalletservice/internal/model/wallet"
	walletModel "userwalletservice/internal/model/wallet"
	userRepo "userwalletservice/internal/repository/user"
	walletRepo "userwalletservice/internal/repository/wallet"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrEmailRequired      = errors.New("email is required")
	ErrPasswordRequired   = errors.New("password is required")
	ErrPasswordTooShort   = errors.New("password is too short (min 3 chars)")
)

type UserService struct {
	repo       *userRepo.UserRepository
	jwtManager *jwt.JWTManager
	walletRepo walletRepo.Repository
}

func New(repo *userRepo.UserRepository, jwtManager *jwt.JWTManager, wRepo walletRepo.Repository) *UserService {
	return &UserService{
		repo:       repo,
		jwtManager: jwtManager,
		walletRepo: wRepo,
	}
}

func (s *UserService) Register(ctx context.Context, email, password string) error {
	email = strings.ToLower(strings.TrimSpace(email))

	if err := validateRegisterInput(email, password); err != nil {
		return err
	}

	hashedPassword, err := s.hashPassword(password)
	if err != nil {
		return err
	}

	user := userModel.User{
		Email:        email,
		PasswordHash: hashedPassword,
	}

	// Сначала регистрируем пользователя в базе
	if err := s.repo.Register(ctx, &user); err != nil {
		return err
	}

	if err := s.walletRepo.Create(ctx, user.ID); err != nil {
		// Если кошелек не создался, возвращаем ошибку
		return err
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.jwtManager.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*userModel.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func validateRegisterInput(email, password string) error {
	if email == "" {
		return ErrEmailRequired
	}

	if password == "" {
		return ErrPasswordRequired
	}

	if len(password) < 3 {
		return ErrPasswordTooShort
	}

	if !strings.Contains(email, "@") {
		return ErrInvalidEmail
	}

	return nil
}

func (s *UserService) GetWalletByUserID(ctx context.Context, userID int) (*walletModel.Wallet, error) {
	// Здесь нам тоже нужен импорт моделей кошелька!
	// Добавь в импорты этого файла (user.go) строчку:
	// walletModel "userwalletservice/internal/model/wallet"

	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *UserService) Deposit(ctx context.Context, userID int, amount float64) (*wallet.Wallet, error) {
	// Валидация: фронтенд не должен прислать отрицательную сумму или ноль
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	// Вызываем репозиторий.
	// Посмотри, как в твоем сервисе называется поле репозитория кошелька
	// (например, s.walletRepo или s.userRepo, если метод Deposit лежит там же)
	return s.walletRepo.Deposit(ctx, userID, amount)
}

func (s *UserService) Withdraw(ctx context.Context, userID int, amount float64) (*wallet.Wallet, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	return s.walletRepo.Withdraw(ctx, userID, amount)
}

func (s *UserService) Transfer(ctx context.Context, fromUserID int, toUserID int, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	if fromUserID == toUserID {
		return errors.New("cannot transfer money to yourself")
	}

	return s.walletRepo.Transfer(ctx, fromUserID, toUserID, amount)
}
