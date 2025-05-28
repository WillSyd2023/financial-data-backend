package main

import (
	"Backend/db"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Setup database
	db, err := db.InitDB()
	if err != nil {
		log.Fatalf("error connect DB: %s", err)
	}
	defer db.Close()
}
