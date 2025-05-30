package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB

func InitDB() {
	errENV := godotenv.Load()
	if errENV != nil {
		log.Println("No .env file available")
	}
	
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "yourHost"),
		getEnv("DB_PORT", "yourPort"),
		getEnv("DB_USER", "yourUser"),
		getEnv("DB_PASSWORD", "yourPassword"),
		getEnv("DB_NAME", "bitebox"),
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	log.Println("Database connection established")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

