package database

import (
    "database/sql"
    "github.com/sirupsen/logrus"
    _ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB // Global variable to hold the database connection

// InitDB initializes the database connection
func InitDB() {
    // Change this connection string according to your PostgreSQL setup
    connStr := "user=postgres password=yourpassword dbname=plantation host=db sslmode=disable"
    var err error

    // Open the database connection
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        logrus.Fatal("Failed to connect to the database:", err)
    }

    // Ping the database to verify the connection
    err = DB.Ping()
    if err != nil {
        logrus.Fatal("Failed to ping the database:", err)
    }

    logrus.Println("Database connection established")
}