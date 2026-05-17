package user

import (
	"context"
	"errors"
	"strings"
	"userwalletservice/internal/infrastructure/jwt"
	userModel "userwalletservice/internal/model/user"
	userRepo "userwalletservice/internal/repository/user"
	
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
	repo 				*userRepo.UserRepository
	jwtManager 	*jwt.JWTManager
}

func New(repo *userRepo.UserRepository, jwtManager *jwt.JWTManager) *UserService {
	return &UserService{
		repo: repo,
		jwtManager: jwtManager,
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

	return s.repo.Register(ctx, &user)
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
