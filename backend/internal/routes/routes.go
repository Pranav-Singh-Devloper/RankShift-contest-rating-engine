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

	// Histroy Route
	api.Get("/users/:id",ratingHandler.GetProfile)

	// all contests 
	api.Get("/contests", ratingHandler.GetAllContests)

	// create contests
	api.Post("/contests", ratingHandler.CreateContest)
}
