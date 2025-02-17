package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseConnection(t *testing.T) {
	db := NewDatabase()

	sqlDB, err := db.DB()
	assert.Nil(t, err)
	assert.Nil(t, sqlDB.Ping())
}
