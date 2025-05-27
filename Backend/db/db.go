package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func InitDB() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	// connecting to DB
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	name := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DBNAME")
	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, name, password, dbname)

	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
