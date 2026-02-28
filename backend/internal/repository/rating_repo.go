// db querires (Prisma specific logic)

package repository 

import (
	"context"
	"fmt"
	"contest-backend/prisma/db"
)

// update_params safely packages all teh data needed to apply a rating change 
type UpdateParams struct{
	UserID string 
	ContestID string 
	OldRating int 
	NewRating int 
	PerformanceRating int
	Rank int 
	Percentile float64 
	RatingChange int 
	NewTier string 
	MaxRating int 
}

// ratign_repository defines the strict contract for out database ops 
type RatingRepository interface{
	GetContest(ctx context.Context, contestID string) (*db.ContestModel,error)
	GetUser(ctx context.Context, userID string) (*db.UserModel,error)
	SaveRatingUpdate(ctx context.Context, params UpdateParams) error
}

type ratingRepo struct {
	client *db.PrismaClient
}

// newRatingRepository is the constructor injecting the prisma client 
func NewRatingRespository(client *db.PrismaClient) RatingRepository{
	return &ratingRepo{
		client: client,
	}
}

// getcontest fetches the contest to retrieve total participants
func (r *ratingRepo) GetContest(ctx context.Context, contestID string) (*db.ContestModel,error){
	contest,err := r.client.Contest.FindUnique(
		db.Contest.ID.Equals(contestID),
	).Exec(ctx)

	if err != nil{
		return nil,fmt.Errorf("failed to fetch contest :%w",err)
	}
	return contest, nil
}

// getUser fetchs the user's current state
func (r *ratingRepo) GetUser(ctx context.Context, userID string) (*db.UserModel,error){
	user, err := r.client.User.FindUnique(
		db.User.ID.Equals(userID),
	).Exec(ctx)
	if err != nil{
		return nil, fmt.Errorf("failed to fetch user: %w",err)
	}
	return user,nil
}

// saveratingupdate uses Prisma Transactions to ensure all-or-nothing data integrity
func (r *ratingRepo) SaveRatingUpdate(ctx context.Context, params UpdateParams)error{

	// prepare the rating history insertion 
	createHistory := r.client.RatingHistory.CreateOne(
		db.RatingHistory.OldRating.Set(params.OldRating),
		db.RatingHistory.NewRating.Set(params.NewRating),
		db.RatingHistory.PerformanceRating.Set(params.PerformanceRating),
		db.RatingHistory.Rank.Set(params.Rank),
		db.RatingHistory.Percentile.Set(params.Percentile),
		db.RatingHistory.RatingChange.Set(params.RatingChange),
		db.RatingHistory.User.Link(db.User.ID.Equals(params.UserID)),
		db.RatingHistory.Contest.Link(db.Contest.ID.Equals(params.ContestID)),
	).Tx()

	// prepare the user state update 
	updateUser := r.client.User.FindUnique(
		db.User.ID.Equals(params.UserID),
	).Update(
		db.User.CurrentRating.Set(params.NewRating),
		db.User.MaxRating.Set(params.MaxRating),
		db.User.Tier.Set(params.NewTier),
		db.User.ContestsPlayed.Increment(1),
	).Tx()

	// execute both operations atomically 
	if err := r.client.Prisma.Transaction(createHistory,updateUser).Exec(ctx);err != nil{
		return fmt.Errorf("transaction failed :%w",err)
	}
	return nil
}