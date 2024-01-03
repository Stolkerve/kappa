package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/Stolkerve/kappa/pgk/db"
	"github.com/Stolkerve/kappa/pgk/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db.SetupDB()
	server.NewServer(":" + os.Getenv("PORT"))
}
