package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	connectEnv := os.Getenv("DB_MODE")
	slog.Info("database connection requested mode", "mode", connectEnv)

	// Get the connection string based on the mode
	var connectionStr string

	// Switch case to get the connection string based on the mode
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

	if connectionStr == "" {
		slog.Error("database connection failed no connection string found")
		return nil, fmt.Errorf("no connection string found")
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

	// Return the database connection
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

	// Format the connection string
	connectString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		dbName,
		sslMode,
	)

	// Return the connection string
	return connectString
}

func connectToDevDB() string {
	// Return the development connection string
	return os.Getenv("DB_DEV")
}

func connectToProdDB() string {
	// Return the production connection string
	return os.Getenv("DB_PROD")
}
