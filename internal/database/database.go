package database

import (
	"fmt"
	"os"

	"github.com/FaizanAC/Go-Banking/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", os.Getenv("HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Cannot connect to the DB")
	}

	return db
}

func MigrateDB(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{},
		&models.BankAccount{},
		&models.Transaction{},
		&models.Transfer{},
	)
}
