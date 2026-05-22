package database

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"userwalletservice/internal/config"
)

// ProvideDB занимается только созданием и настройкой подключения
func ProvideDB(cfg *config.DB) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
	)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия БД: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка ping БД: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

// ConnectWithRetry — обёртка, которая скрывает в себе логику повторных попыток
func ConnectWithRetry(cfg *config.DB) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	for i := range 10 {
		db, err = ProvideDB(cfg)
		if err == nil {
			return db, nil // Успешно подключились — возвращаем базу наружу
		}

		log.Printf("⏳ DB not ready (%v), retrying... attempt: %d/10\n", err, i+1)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("all connection attempts failed: %w", err)
}
