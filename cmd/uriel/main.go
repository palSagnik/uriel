package main

import (
	"log"

	"github.com/palSagnik/uriel/internal/database"
)


func main() {

	// connect to database
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("database is not live: %v", err)
	}
	log.Println("connected to database")
}