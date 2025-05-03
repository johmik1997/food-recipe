package database

import (
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect() *sqlx.DB {
	// // Get configuration from environment variables
	// host := getEnv("DB_HOST", "localhost")
	// port := getEnv("DB_PORT", "5432")
	// user := getEnv("DB_USER", "postgres")
	// password := getEnv("DB_PASSWORD", "199497")
	// dbname := getEnv("DB_NAME", "mydatabase")

	// connStr := fmt.Sprintf(
	// 	"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	// 	host, port, user, password, dbname,
	// )
	connStr := getEnv("DB_CONNECTION", "postgres://postgres:postgrespassword@postgres:5432/postgres?sslmode=disable")
	
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	log.Println("Successfully connected to PostgreSQL")
	return db
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}