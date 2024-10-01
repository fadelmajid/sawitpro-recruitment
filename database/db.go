package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB // Global variable to hold the database connection

// InitDB initializes the database connection
func InitDB() {
	// Change this connection string according to your PostgreSQL setup
	connStr := "user=postgres dbname=plantation sslmode=disable password=yourpassword"
	var err error

	// Open the database connection
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Ping the database to verify the connection
	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping the database:", err)
	}

	log.Println("Database connection established")
}
