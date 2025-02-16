package database

import (
	"testing"
)

func TestDatabaseConnection(t *testing.T) {
	db := NewDatabase()

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Invalid DB")
	}
	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("No response from DB")
	}
}
