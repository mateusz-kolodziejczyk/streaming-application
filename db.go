package main

import (
	"database/sql"
	"fmt"
	"os"
)

// connectToDB Connects to PostgreSQL database using environmental variables
func connectToDB() (*sql.DB, error) {
	// Connect to database
	// Connection string is made up of the various environemntal variables.
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_ADDRESS"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)
	return sql.Open("pgx", connectionString)

}