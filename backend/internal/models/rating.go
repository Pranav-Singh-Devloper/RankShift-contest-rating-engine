// struchs and DTOs

package models

// UserResult represents a single user's rank in the contest 
type UserResult struct {
	UserID string `json:"user_id"`
	Rank int `json:"rank`
}

// ContestEndPayload is the exact JSON structure the backen expects 
type ContestEndPayload struct {
	ContestID string `json:"contest_id"`
	Results []UserResult `json:"resutls"`
}