//  env and db setup

package config 

import (
	"log"
	"os"
	"contest-backend/prisma/db"
	"github.com/joho/godotenv"
)

// connectdb initializes and returns the Prismaclient instance
func ConnectDB() *db.PrismaClient{
	if os.Getenv("APP_ENV") != "production"{
		_ = godotenv.Load()
	}
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil{
		log.Fatalf("Database connection failed: %v",err)
	}
	
	log.Println("Database connected successfully via config layer")
	return client 
}

// disconnectdb cleanly closes the connection during server shutdown
func DisconnectDB(client *db.PrismaClient){
	if err := client.Prisma.Disconnect(); err != nil{
		log.Printf("Failed to disconnect database: %v",err)
	}
}