package trade

import "time"

// Transaction represents a transaction between two accounts.
type Transaction struct {
	ID         string             `json:"id"`
	Quantities map[string]float64 `json:"quantities"`
	Sender     string             `json:"sender"`
	Recipient  string             `json:"recipient"`
	Timestamp  time.Time          `json:"timestamp"`
}
