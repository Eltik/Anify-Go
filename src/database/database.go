package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var (
	// DB is the connection pool to the database
	DB *pgx.Conn
)

func Connect() (*pgx.Conn, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	println("Connecting to database at", os.Getenv("DATABASE_URL"))

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
		return nil, err
	}

	DB = conn

	err = conn.Ping(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	println("Connected to database")

	return conn, nil
}
