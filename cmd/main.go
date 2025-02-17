package main

import (
	"log"
	"os"

	"github.com/FaizanAC/Go-Banking/internal/database"
	"github.com/FaizanAC/Go-Banking/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := database.NewDatabase()
	database.MigrateDB(db)

	s := server.NewServer(
		db, os.Getenv("PORT"),
	)
	s.Start()
}
