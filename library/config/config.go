package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DB_USER    string
	DB_PASS    string
	DB_CLUSTER string
	DB_NAME    string
	JWT_SECRET string
	PORT       string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	DB_USER = os.Getenv("DB_USER")
	DB_PASS = os.Getenv("DB_PASS")
	DB_CLUSTER = os.Getenv("DB_CLUSTER")
	DB_NAME = os.Getenv("DB_NAME")
	JWT_SECRET = os.Getenv("JWT_SECRET")
	PORT = port
}
