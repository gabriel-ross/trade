package trade

// User represents a user.
type User struct {
	ID          string `json:"_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}
