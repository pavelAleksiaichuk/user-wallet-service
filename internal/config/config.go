package config

import "os"

// DB — структура, которую запросил ментор
type DB struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

// Config — общая обёртка для настроек приложения
type Config struct {
	DB        DB
	JWTSecret string
}

// New собирает конфигурацию вручную с помощью стандартного пакета os
func New() *Config {
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432" // Дефолтный порт PostgreSQL, если он не задан в .env
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super-secret-wallet-key-2026"
	}

	return &Config{
		DB: DB{
			Host:     os.Getenv("DB_HOST"),
			Port:     port,
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		},
		JWTSecret: jwtSecret,
	}
}
