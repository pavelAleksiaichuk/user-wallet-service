package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func Connect() {

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var db *sqlx.DB
	var err error

	// 🔥 retry logic
	for i := 0; i < 10; i++ {

		db, err = sqlx.Connect("postgres", dsn)
		if err == nil {

			err = db.Ping()
			if err == nil {
				DB = db
				fmt.Println("✅ PostgreSQL connected")
				return
			}
		}

		fmt.Println("⏳ DB not ready, retrying... attempt:", i+1)
		time.Sleep(2 * time.Second)
	}

	log.Fatal("❌ DB connection failed after retries:", err)
}