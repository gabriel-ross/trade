package trade

import "time"

// Transaction represents a transaction between two accounts.
type Transaction struct {
	ID          string             `json:"id"`
	Quantities  map[string]float64 `json:"quantities"`
	SenderID    string             `json:"senderID"`
	RecipientID string             `json:"recipientID"`
	Timestamp   time.Time          `json:"timestamp"`
}
