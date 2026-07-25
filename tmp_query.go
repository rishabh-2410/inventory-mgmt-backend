package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load()
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, name, sslmode)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, business_id, name, address, created_at FROM warehouses ORDER BY created_at`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	fmt.Println("warehouses:")
	for rows.Next() {
		var id, businessID, name, createdAt string
		var address sql.NullString
		rows.Scan(&id, &businessID, &name, &address, &createdAt)
		fmt.Printf("  id=%s business_id=%s name=%q created_at=%s\n", id, businessID, name, createdAt)
	}

	rows2, err := db.Query(`SELECT id, name FROM business ORDER BY created_at`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows2.Close()
	fmt.Println("\nbusinesses:")
	for rows2.Next() {
		var id, name string
		rows2.Scan(&id, &name)
		fmt.Printf("  id=%s name=%q\n", id, name)
	}
}
