package main 
import (
	"log"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"contest-backend/prisma/db"
)

var PrismaClient *db.PrismaClient

func main(){
	if os.Getenv("APP_ENV") != "production"{
		_ = godotenv.Load()
	}

	// database connection logic 
	PrismaClient = db.NewClient()
	if err := PrismaClient.Prisma.Connect(); err != nil {
		log.Println("Database connection failed:",err)
		os.Exit(1)
	}
	defer PrismaClient.Prisma.Disconnect()
	
	log.Println("Database connected successfull")

	// initialize the go fiber server 
	app := fiber.New()

	// simple health check
	app.Get("/health", func(c *fiber.Ctx) error{
		return c.SendString("Contest Rating Engine is live!")
	})

	port := os.Getenv("PORT")
	if port == ""{
		port = "8080"
	}

	log.Printf("String server on port %s",port)
	if err := app.Listen(":"+port); err != nil{
		log.Fatalf("Server failde to start: %v",err)
	}


}