package handlers

import (
	"database/sql"
	"os"
	"testing"

	"github.com/tonitran/dischord/store"
)

func testStore(t *testing.T) *store.Store {
	t.Helper()
	connStr := os.Getenv("TEST_DATABASE_URL")
	if connStr == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}
	if err := store.ApplySchema(db); err != nil {
		t.Fatal(err)
	}
	if err := store.TruncateAll(db); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		store.TruncateAll(db)
		db.Close()
	})
	return store.New(db)
}
