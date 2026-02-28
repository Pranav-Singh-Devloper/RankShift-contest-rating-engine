package config

import (
	"log"
	"os"
	"time"

	"contest-backend/prisma/db"
	"github.com/joho/godotenv"
)

// ConnectDB initializes and returns the Prisma client instance with retry logic
func ConnectDB() *db.PrismaClient {
	if os.Getenv("APP_ENV") != "production" {
		_ = godotenv.Load()
	}

	client := db.NewClient()
	
	// Retry logic for Serverless/Neon environments (Fixes 57P01)
	maxRetries := 5
	for i := 1; i <= maxRetries; i++ {
		err := client.Prisma.Connect()
		if err == nil {
			log.Printf("Database connected successfully (Attempt %d)", i)
			return client
		}

		log.Printf("Database connection attempt %d failed: %v. Retrying in 2s...", i, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("Database connection failed after %d attempts. Exiting.", maxRetries)
	return nil
}

// DisconnectDB cleanly closes the connection during server shutdown
func DisconnectDB(client *db.PrismaClient) {
	if err := client.Prisma.Disconnect(); err != nil {
		log.Printf("Failed to disconnect database: %v", err)
	}
}