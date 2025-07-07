package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func main() {
	
	// Gets the connection string from environment variables
	dsn := os.Getenv("DIRECT_URL")
	if dsn == "" {
		log.Fatal("DIRECT_URL not found in .env file")
	}

	// Connects to the database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Tests the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database:", err)
	}

	fmt.Println(" Successfully connected to Supabase database!")
}
