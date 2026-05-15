package database

import (
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect() (*sqlx.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		name,
	)

	var db *sqlx.DB
	var err error

	// 🔥 retry logic
	for i := range 10 {
		db, err = sqlx.Connect("postgres", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				fmt.Println("✅ PostgreSQL connected")
				db.SetMaxOpenConns(25)
				db.SetMaxIdleConns(25)

				return db, nil
			}
		}

		fmt.Printf("⏳ DB not ready (%v), retrying... attempt: %d\n", err, i+1)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("❌ DB connection failed after 10 retries: %w", err)
}
