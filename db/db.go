package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func ConnectDatabase() {
	var err error
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	connectionString := fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s host=%s port=%s", dbUser, dbName, dbSSLMode, dbPassword, dbHost, dbPort)

	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
}

func GetDB() *sql.DB {
	return db
}
