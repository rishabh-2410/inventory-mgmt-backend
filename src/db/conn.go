package db

import (
	"database/sql"
	"fmt"
	"log/slog"
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

	connectEnv := os.Getenv("DB_MODE")
	slog.Info("database connection requested mode", "mode", connectEnv)

	var connectionStr string

	switch connectEnv {
	case "dev":
		connectionStr = connectToDevDB()
	case "prod":
		connectionStr = connectToProdDB()
	case "local":
		connectionStr = connectToLocalDB()
	default:
		slog.Error("database connection failed invalid mode", "mode", connectEnv)
		return nil, fmt.Errorf("invalid connection mode: %s", connectEnv)
	}

	// Open the connection
	db, err := sql.Open("postgres", connectionStr)
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

func connectToLocalDB() string {
	// Initialize string parameters
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSLMODE")

	connectString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		dbName,
		sslMode,
	)
	return connectString
}

func connectToDevDB() string {
	return os.Getenv("DB_DEV")
}

func connectToProdDB() string {
	return os.Getenv("DB_PROD")
}
