package user

import (
	"context"
	"strings"

	"userwalletservice/internal/infrastructure/jwt"
	userModel "userwalletservice/internal/model/user"
	"userwalletservice/internal/repository"
	userRepository "userwalletservice/internal/repository/user"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepository   *userRepository.Repository
	jwtManager       *jwt.JWTManager
	walletRepository repository.Wallet // Оставляем только для транзакции в Register
}

func New(userRepository *userRepository.Repository, jwtManager *jwt.JWTManager, walletRepository repository.Wallet) *Service {
	return &Service{
		userRepository:   userRepository,
		jwtManager:       jwtManager,
		walletRepository: walletRepository,
	}
}

func (s *Service) Register(ctx context.Context, email, password string) (userModel.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	user := userModel.User{
		Email: email,
	}

	if err := user.Validate(password); err != nil {
		return user, err
	}

	hashedPassword, err := s.hashPassword(password)
	if err != nil {
		return user, err
	}
	user.PasswordHash = hashedPassword

	// 🔥 1. Убрали & перед user.
	// 🔥 2. Теперь репозиторий возвращает ID, и мы присваиваем его локальной копии перед return
	userID, err := s.userRepository.Register(ctx, user, s.walletRepository)
	if err != nil {
		return user, err
	}

	user.ID = userID // Теперь сервис честно возвращает заполненный объект

	return user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return "", userModel.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", userModel.ErrInvalidCredentials
	}

	token, err := s.jwtManager.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) GetUserByID(ctx context.Context, id int) (*userModel.User, error) {
	return s.userRepository.GetByID(ctx, id)
}

func (s *Service) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
