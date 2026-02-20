package main

import (
	"log"
	"net/http"
	"os"

	"github.com/tonitran/dischord/router"
	"github.com/tonitran/dischord/store"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://localhost/dischord?sslmode=disable"
	}
	s, err := store.Open(connStr)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	handler := router.New(s)
	log.Println("DisChord server starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
