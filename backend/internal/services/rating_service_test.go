package services
import (
	"testing"
	"github.com/stretchr/testify/assert"
)
//  testing 
func TestAssessmentDocumentExample(t * testing.T){

	totalParticipants := 100
	rank := 10
	userCurrentRating := 1000

	// Step 1 Beaten
	beaten := totalParticipants - rank
	assert.Equal(t, 90, beaten, "Beaten should be 90")

	// Step 2 Percentile
	percentile := float64(beaten) / float64(totalParticipants)
	assert.Equal(t, 0.90, percentile, "Percentile should be 0.90")

	// Step 3 & 4 Performance Bracket
	performanceRating := getPerformanceRating(percentile)
	assert.Equal(t, 1200, performanceRating, "Standard Performance should be 1200")

	// Step 5: Rating Change
	ratingChange := (performanceRating - userCurrentRating) / 2
	assert.Equal(t, 100, ratingChange, "Rating Change should be 100")

	// Step 6: New Rating
	newRating := userCurrentRating + ratingChange
	assert.Equal(t, 1100, newRating, "New Rating should be 1100")

}

// performance brackets testing
func TestGetPerformanceRating(t *testing.T){
	tests := []struct {
		name string
		percentile float64
		expected int
	}{
		{"Top 1%", 0.995, 1800},
		{"Top 5%", 0.95, 1400},
		{"Top 10%", 0.90, 1200},
		{"Top 20%", 0.85, 1150},
		{"Top 30%", 0.75, 1100},
		{"Top 50%", 0.50, 1000},
		{"Bottom 50% (Floor)", 0.20, 800},
	}
	for _, tt := range tests{
		t.Run(tt.name, func(t *testing.T){
			result := getPerformanceRating(tt.percentile)
			assert.Equal(t,tt.expected,result)
		})
	}
}

// tiers testing
func TestGetTier(t *testing.T) {
	tests := []struct {
		name     string
		rating   int
		expected string
	}{
		{"Diamond", 1850, "Diamond"},
		{"Platinum", 1500, "Platinum"},
		{"Gold", 1250, "Gold"},
		{"Silver", 1050, "Silver"},
		{"Bronze", 900, "Bronze"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getTier(tt.rating)
			assert.Equal(t, tt.expected, result)
		})
	}
}
