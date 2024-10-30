package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewDatabase(connectionString string) (*sql.DB, error) {
	log.Println("Connecting to the database")
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Println(err)
		return db, err
	}

	log.Println("Pinging the database")
	if err := db.Ping(); err != nil {
		log.Println(err)
		return db, err
	}
	log.Println("Successfully connected to the database")
	return db, nil
}
