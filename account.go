package trade

import "time"

// Account represents a trading account.
type Account struct {
	ID                string             `json:"id"`
	Owner             string             `json:"owner"`
	Balances          map[string]float64 `json:"balances"`
	Reputation        int                `json:"reputation"`
	CreationTimestamp time.Time          `json:"creationTimestamp"`
}
