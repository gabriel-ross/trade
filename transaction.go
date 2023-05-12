package trade

import "time"

// Transaction represents a transaction between two accounts.
type Transaction struct {
	Sender     string             `json:"_from"`
	Recipient  string             `json:"_to"`
	ID         string             `json:"_id"`
	Quantities map[string]float64 `json:"quantities"`
	Timestamp  time.Time          `json:"timestamp"`
}
