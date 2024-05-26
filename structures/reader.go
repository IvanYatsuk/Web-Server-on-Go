package structures

import "time"

type Reader struct {
	ID               int       `json:"id"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	PhoneNumber      string    `json:"phone_number"`
	RegistrationDate time.Time `json:"registration_date"`
	Notes            string    `json:"notes"`
}
