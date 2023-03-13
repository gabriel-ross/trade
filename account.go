package trade

import "time"

// Account represents a trading account.
type Account struct {
	ID                string
	Owner             string
	Balances          map[string]int
	Reputation        int
	CreationTimestamp time.Time
}
