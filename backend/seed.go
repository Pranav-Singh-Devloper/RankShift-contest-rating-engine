// backend/seed.go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	// IMPORTANT: Replace "contest-backend" with your actual go.mod module name
	"contest-backend/internal/config"
	"contest-backend/prisma/db"
)

func main() {
	ctx := context.Background()

	// 1. Connect to the database
	client := config.ConnectDB()
	defer config.DisconnectDB(client)

	fmt.Println("🌱 Starting database seed...")

	// 2. Create the Primary User
	// Simulating a profile with a 1306 current rating and a 1462 peak.
	user, err := client.User.CreateOne(
		db.User.Name.Set("Pranav Singh"),
		db.User.CurrentRating.Set(1306),
		db.User.MaxRating.Set(1462),
		db.User.ContestsPlayed.Set(6),
		db.User.Tier.Set("Gold"),
	).Exec(ctx)

	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	fmt.Printf("✅ Created User: %s (ID: %s)\n", user.Name, user.ID)

	// 3. Define a realistic chronological history of contests
	type ContestSeed struct {
		Name              string
		Date              time.Time
		TotalParticipants int
		Rank              int
		OldRating         int
		NewRating         int
		PerfRating        int
		Percentile        float64
		RatingChange      int
	}

	history := []ContestSeed{
		{
			Name:              "Visual Vortex 2.0 Hackathon",
			Date:              time.Date(2025, 1, 10, 10, 0, 0, 0, time.UTC),
			TotalParticipants: 100, Rank: 10, OldRating: 1000, NewRating: 1100, PerfRating: 1200, Percentile: 0.90, RatingChange: 100,
		},
		{
			Name:              "DECODE 2025",
			Date:              time.Date(2025, 10, 5, 10, 0, 0, 0, time.UTC),
			TotalParticipants: 100, Rank: 25, OldRating: 1100, NewRating: 1100, PerfRating: 1100, Percentile: 0.75, RatingChange: 0,
		},
		{
			Name:              "HackX 3.0 Jaipur",
			Date:              time.Date(2025, 11, 12, 10, 0, 0, 0, time.UTC),
			TotalParticipants: 200, Rank: 10, OldRating: 1100, NewRating: 1250, PerfRating: 1400, Percentile: 0.95, RatingChange: 150,
		},
		{
			Name:              "Mumbai Hacks 2025 Agentic AI",
			Date:              time.Date(2025, 11, 28, 10, 0, 0, 0, time.UTC),
			TotalParticipants: 150, Rank: 60, OldRating: 1250, NewRating: 1125, PerfRating: 1000, Percentile: 0.60, RatingChange: -125, // A dip in the chart!
		},
		{
			Name:              "Codeforces Global Round",
			Date:              time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC),
			TotalParticipants: 500, Rank: 5, OldRating: 1125, NewRating: 1462, PerfRating: 1800, Percentile: 0.99, RatingChange: 337, // Massive spike!
		},
		{
			Name:              "Educational Codeforces Round",
			Date:              time.Date(2026, 2, 20, 10, 0, 0, 0, time.UTC),
			TotalParticipants: 300, Rank: 45, OldRating: 1462, NewRating: 1306, PerfRating: 1150, Percentile: 0.85, RatingChange: -156,
		},
	}

	// 4. Insert Contests and Rating History
	for _, data := range history {
		// Create the contest and set its date
		contest, err := client.Contest.CreateOne(
			db.Contest.Name.Set(data.Name),
			db.Contest.TotalParticipants.Set(data.TotalParticipants),
			db.Contest.Date.Set(data.Date),
		).Exec(ctx)

		if err != nil {
			log.Fatalf("Failed to create contest %s: %v", data.Name, err)
		}

		// Create the associated rating history ledger entry
		_, err = client.RatingHistory.CreateOne(
			db.RatingHistory.OldRating.Set(data.OldRating),
			db.RatingHistory.NewRating.Set(data.NewRating),
			db.RatingHistory.PerformanceRating.Set(data.PerfRating),
			db.RatingHistory.Rank.Set(data.Rank),
			db.RatingHistory.Percentile.Set(data.Percentile),
			db.RatingHistory.RatingChange.Set(data.RatingChange),
			db.RatingHistory.User.Link(db.User.ID.Equals(user.ID)),
			db.RatingHistory.Contest.Link(db.Contest.ID.Equals(contest.ID)),
		).Exec(ctx)

		if err != nil {
			log.Fatalf("Failed to create history for %s: %v", data.Name, err)
		}
	}

	fmt.Println("✅ Successfully seeded 6 historical contest records.")
	fmt.Println("--------------------------------------------------")
	fmt.Printf("🔥 TEST URL: http://localhost:3000/profile/%s\n", user.ID)
	fmt.Println("--------------------------------------------------")
}