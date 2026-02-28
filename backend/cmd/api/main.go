package main 
import (
	"log"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"contest-backend/internal/config"
	"contest-backend/internal/handlers"
	"contest-backend/internal/repository"
	"contest-backend/internal/routes"
	"contest-backend/internal/services"
)


func main(){
	dbClient := config.ConnectDB()
	defer config.DisconnectDB(dbClient)

	ratingRepo := repository.NewRatingRepository(dbClient)
	ratingService := services.NewRatingService(ratingRepo)
	ratingHandler := handlers.NewRatingHandler(ratingService)

	app := fiber.New(fiber.Config{
		AppName: "Contest Rating Engine v1.0",
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, OPTIONS",
	}))
	app.Use(logger.New())
	app.Use(recover.New())

	routes.SetupRoutes(app, ratingHandler)

	port := os.Getenv("PORT")
	if port == ""{
		port = "8080"
	}

	log.Printf("Starting Clean Architecture server on port %s...",port)
	if err := app.Listen(":"+port);err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}