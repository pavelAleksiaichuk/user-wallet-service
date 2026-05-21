package user

import (
	"context"
	"errors"
	"strings"

	"userwalletservice/internal/infrastructure/jwt"
	userModel "userwalletservice/internal/model/user"
	userRepo "userwalletservice/internal/repository/user"
	walletRepo "userwalletservice/internal/repository/wallet"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type Service struct {
	repo       *userRepo.Repository
	jwtManager *jwt.JWTManager
	walletRepo walletRepo.Repository // Оставляем только для транзакции в Register
}

func New(repo *userRepo.Repository, jwtManager *jwt.JWTManager, wRepo walletRepo.Repository) *Service {
	return &Service{
		repo:       repo,
		jwtManager: jwtManager,
		walletRepo: wRepo,
	}
}

func (s *Service) Register(ctx context.Context, email, password string) error {
	email = strings.ToLower(strings.TrimSpace(email))

	user := userModel.User{
		Email: email,
	}

	if err := user.Validate(password); err != nil {
		return err
	}

	hashedPassword, err := s.hashPassword(password)
	if err != nil {
		return err
	}
	user.PasswordHash = hashedPassword

	if err := s.repo.Register(ctx, &user, s.walletRepo); err != nil {
		return err
	}

	return nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
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

func (s *Service) GetUserByID(ctx context.Context, id int) (*userModel.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
