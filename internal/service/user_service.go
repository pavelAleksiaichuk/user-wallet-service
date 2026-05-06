package service

import (
	"errors"
	"strings"

	"userwalletservice/internal/model"
	"userwalletservice/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CREATE USER
func (s *UserService) CreateUser(email string) error {

	// 1. простая бизнес-валидация
	if email == "" {
		return errors.New("email is required")
	}

	if !strings.Contains(email, "@") {
		return errors.New("invalid email format")
	}

	user := model.User{
		Email: email,
	}

	// 2. вызываем repository
	return s.repo.Create(user)
}

// GET USER
func (s *UserService) GetUser(id int) (model.User, error) {
	return s.repo.GetByID(id)
}