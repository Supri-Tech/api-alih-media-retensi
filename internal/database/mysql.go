package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() *sql.DB {
	dsn := os.Getenv("DB_DSN")
	log.Println("Connecting to DB with DSN:", dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping DB:", err)
	}

	log.Println("Successfully connected to the database")
	return db
}

// func InitDB() *sql.DB {
// 	dsn := os.Getenv("DB_DSN")
// 	// dbClient := os.Getenv("DB_CLIENT")

// 	db, err := sql.Open("mysql", dsn)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if err := db.Ping(); err != nil {
// 		log.Fatal(err)
// 	}

// 	return db
// }
