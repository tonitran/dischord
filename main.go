package main

import (
	"log"
	"net/http"

	"github.com/tonitran/dischord/router"
	"github.com/tonitran/dischord/store"
)

func main() {
	s := store.New()
	handler := router.New(s)

	log.Println("DisChord server starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
