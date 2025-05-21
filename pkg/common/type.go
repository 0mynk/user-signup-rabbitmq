package common

// Add UserEvent struct definition
type UserEvent struct {
	Email     string        `json:"email"`
	FirstName string        `json:"first_name"`
	LastName  string        `json:"last_name"`
}

type InvoiceData = UserEvent

