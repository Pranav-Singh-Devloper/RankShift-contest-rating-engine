package handlers 

import (
	"context"
	"time"
	"github.com/gofiber/fiber/v2"
	"contest-backend/internal/models"
	"contest-backend/internal/services"
)

// ratingHandler 
type RatingHandler interface {
	HandleContestEnd(c *fiber.Ctx) error
	GetProfile(c *fiber.Ctx) error
}

type ratingHandler struct {
	service services.RatingService
}

// constructor for injecting the service dependency
func NewRatingHandler(service services.RatingService) RatingHandler{
	return &ratingHandler{
		service :service,
	}
}

// process the POST request when a contest finishes
func (h *ratingHandler) HandleContestEnd(c * fiber.Ctx) error {
	var payload models.ContestEndPayload 

	// parsing the incoming JSON body
	if err := c.BodyParser(&payload); err != nil {
    	// fmt.Printf("BODY PARSE ERROR: %v", err) only used for testing 
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"Invalid request payload format",
		})
	}

	// basic paylad validation 
	if payload.ContestID == "" || len(payload.Results)==0{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"contest_id and a pupulated results array are required",
		})
	}

	// preenting hanging db queries by creating a timeout context 
	ctx, cancle := context.WithTimeout(context.Background(),10*time.Second)
	defer cancle()

	// passing the validated data to the service layer for mathematical processing
	if err := h.service.ProcessContestResults(ctx,payload);err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":"Failed to process contest results",
			"details": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Message":"Contest ratings calculated and updated successfully",
	})

}

func (h *ratingHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	user,err := h.service.GetUserProfile(ctx,userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":"User not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}