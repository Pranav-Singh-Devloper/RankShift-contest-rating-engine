package services

import (
	"contest-backend/internal/models"
	"contest-backend/internal/repository"
	"contest-backend/prisma/db"
	"context"
	"fmt"
	"math"

)

// ratingService defines the interface for core business logic
type RatingService interface {
	ProcessContestResults(ctx context.Context, payload models.ContestEndPayload) error
	GetUserProfile(ctx context.Context, userID string) (*db.UserModel, error)
}

type ratingService struct{
	repo repository.RatingRepository
}

// newRatingService injects the repository dependency 
func NewRatingService(repo repository.RatingRepository) RatingService{
	return &ratingService{
		repo: repo,
	}
}

// MATH: Expected Upset Probability (Logistic Curve)
// Calculates the statistical probability of the user achieving this performance
// using the Elo logistic cumulative distribution function (CDF).
// Formula: E = 1 / (1 + 10^((Perf - Rating)/400))

func calculateUpsetProbability(currentRating,performanceRating int) float64{
	exponent := float64(performanceRating-currentRating) / 400.0 
	probability := 1.0/(1.0 + math.Pow(10,exponent))
	return probability
}

func getPerformanceRating(percentile float64) int {
	if percentile >= 0.99 {return 1800}
	if percentile >= 0.95 {return 1400}
	if percentile >= 0.90 {return 1200}
	if percentile >= 0.80 {return 1150}
	if percentile >= 0.70 {return 1100}
	if percentile >= 0.50 {return 1000}
	return 800 // for bottom 50%
}

func getTier(rating int) string{
	if rating >= 1800 {return "Diamond"}
	if rating >= 1400 {return "Platinum"}
	if rating >= 1200 {return "Gold"}
	if rating >= 1000 {return "Silver"}
	return "Bronze"
}

// processContestResults 
func (s *ratingService) ProcessContestResults(ctx context.Context,payload models.ContestEndPayload)error{
	contest, err := s.repo.GetContest(ctx,payload.ContestID)
	if err != nil {
		return fmt.Errorf("invalid contest :%w",err)
	}

	totalParticipants := float64(contest.TotalParticipants)

	for _, result := range payload.Results{
		user, err := s.repo.GetUser(ctx, result.UserID)
		if err != nil{
			fmt.Printf("User %s not found, skipping...\n",result.UserID)
			continue
		}
		// step-1
		beaten := contest.TotalParticipants - result.Rank 

		// step-2
		percentile := float64(beaten) / totalParticipants

		// step-3 & 4 (Bracket & Standard Performance)
		perfomanceRating := getPerformanceRating(percentile)

		// maths
		upsetProb := calculateUpsetProbability(user.CurrentRating,perfomanceRating)
		fmt.Printf("User %s | Perf: %d | Anomaly Prob: %.2f%%\n",user.ID,perfomanceRating,upsetProb*100)

		// step-5 (A fixed-gain Kalman Filter (K=0.5))
		ratingChange := (perfomanceRating - user.CurrentRating)/2

		// step-6 
		newRating := user.CurrentRating + ratingChange 

		// setp - 7 
		newTier := getTier(newRating)

		maxRating := user.MaxRating
		if newRating > maxRating{
			maxRating = newRating
		}

		// packaging the  calculated data 
		params := repository.UpdateParams{
			UserID: user.ID,
			ContestID: contest.ID,
			NewRating: newRating,
			PerformanceRating: perfomanceRating,
			Rank: result.Rank,
			Percentile: percentile,
			RatingChange: ratingChange,
			NewTier: newTier, 
			MaxRating: maxRating,
		}
		
		// executing ACID transaction via the repo
		if err := s.repo.SaveRatingUpdate(ctx,params);err != nil{
			return fmt.Errorf("failed to save upadtes for user %s: %w",user.ID,err)
		}
	}
	return nil
}

func (s * ratingService) GetUserProfile(ctx context.Context, userID string) (*db.UserModel, error){
	return s.repo.GetUserProfile(ctx,userID)
}

