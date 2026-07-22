package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {

	// Load env variables from .env file
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	// Initialize string parameters
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSLMODE")

	// connectEnv := os.Getenv("DB_MODE")
	connectString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		dbName,
		sslMode,
	)

	// Open the connection
	db, err := sql.Open("postgres", connectString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	// Ping the db
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return db, nil

}
