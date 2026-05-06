package repository

import (
	"userwalletservice/internal/infrastructure/database"
	"userwalletservice/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.DB,
	}
}

// CREATE
func (r *UserRepository) Create(user model.User) error {
	query := `INSERT INTO users (email) VALUES ($1)`
	_, err := r.db.Exec(query, user.Email)
	return err
}

// GET BY ID
func (r *UserRepository) GetByID(id int) (model.User, error) {
	var user model.User

	query := `SELECT id, email FROM users WHERE id=$1`
	err := r.db.Get(&user, query, id)

	return user, err
}