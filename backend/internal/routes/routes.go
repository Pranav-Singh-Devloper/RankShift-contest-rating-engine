package routes 

import (
	"github.com/gofiber/fiber/v2"
	"contest-backend/internal/handlers"
)

func SetupRoutes(app *fiber.App, ratingHandler handlers.RatingHandler){
	app.Get("/health",func(c *fiber.Ctx) error{
		return c.SendString("Contest Rating Engine is live")
	})

	//  api grouping
	api := app.Group("/api")

	// Contest Routes 
	api.Post("/contests/end", ratingHandler.HandleContestEnd)
}
