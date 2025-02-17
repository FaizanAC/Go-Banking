package main

import (
	"log"
	"os"

	"github.com/FaizanAC/Go-Banking/internal/database"
	"github.com/FaizanAC/Go-Banking/internal/models"
	"github.com/FaizanAC/Go-Banking/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := database.NewDatabase()
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.BankAccount{})
	db.AutoMigrate(&models.Transaction{})
	db.AutoMigrate(&models.Transfer{})

	s := server.NewServer(
		db, os.Getenv("PORT"),
	)
	s.Start()
}
