package trade

import "time"

// Transaction represents a transaction between two accounts.
type Transaction struct {
	ID          string
	Quantities  map[string]int
	SenderID    string
	RecipientID string
	Timestamp   time.Time
}
