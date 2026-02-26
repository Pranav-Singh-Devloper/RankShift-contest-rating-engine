package main 

import (
	"context"
	"fmt"
	"log"
	"time"
	"contest-backend/prisma/db"
	"github.com/joho/godotenv"
	"os"
)
func main() {
// 1. Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is empty! Check your .env file.")
	}
	fmt.Println("Environment loaded. Attempting to connect to the database...")

	// 2. Connect to Prisma
	client := db.NewClient()
	
	// Use %v and pass err directly - no .Error() needed here!
	if err := client.Prisma.Connect(); err != nil {
		log.Fatalf("Failed to connect to Prisma:")
	}
	defer client.Prisma.Disconnect()

	ctx := context.Background()
	fmt.Println("Successfully connected to NeonDB! Seeding data...")

	// dummy contest
	contest, err := client.Contest.CreateOne(
		db.Contest.Name.Set("Weekly Algorithm Challenge #1"),
		db.Contest.TotalParticipants.Set(100),
		db.Contest.Date.Set(time.Now()),
	).Exec(ctx)
	if err != nil {
		log.Fatalf("Failed to create contest: %v",err)
	}
	fmt.Printf("Created Contest: %s\n",contest.Name)

	// dummy users
	user1, _ := client.User.CreateOne(
		db.User.Name.Set("Alice Hacker"),
		db.User.CurrentRating.Set(1200),
		db.User.MaxRating.Set(1200),
		db.User.Tier.Set("Silver"),
		db.User.ContestsPlayed.Set(1),
	).Exec(ctx)

	user2, _ := client.User.CreateOne(
		db.User.Name.Set("Bob Builder"),
		db.User.CurrentRating.Set(1000),
		db.User.MaxRating.Set(1000),
		db.User.Tier.Set("Bronze"),
		db.User.ContestsPlayed.Set(1),
	).Exec(ctx)	

	fmt.Println("Created Users: Alice Hacker & Bob Builder")

	// dummy rating history 
	_, err = client.RatingHistory.CreateOne(
		db.RatingHistory.OldRating.Set(1000),
		db.RatingHistory.NewRating.Set(1200),
		db.RatingHistory.PerformanceRating.Set(1400),
		db.RatingHistory.Rank.Set(5),
		db.RatingHistory.Percentile.Set(0.95),
		db.RatingHistory.RatingChange.Set(200),
		db.RatingHistory.User.Link(db.User.ID.Equals(user1.ID)),
		db.RatingHistory.Contest.Link(db.Contest.ID.Equals(contest.ID)),
	).Exec(ctx)

	_, err = client.RatingHistory.CreateOne(
		db.RatingHistory.OldRating.Set(1000),
		db.RatingHistory.NewRating.Set(1000),
		db.RatingHistory.PerformanceRating.Set(1000),
		db.RatingHistory.Rank.Set(40),
		db.RatingHistory.Percentile.Set(0.60),
		db.RatingHistory.RatingChange.Set(0),
		db.RatingHistory.User.Link(db.User.ID.Equals(user2.ID)),
		db.RatingHistory.Contest.Link(db.Contest.ID.Equals(contest.ID)),
	).Exec(ctx)


	if err != nil{
		log.Fatalf("Failed to create rating history: %v",err)
	}
	fmt.Println("Databas successfully seeded with dummy data!")
}