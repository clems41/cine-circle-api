package model

type User struct {
	ID             uint   `json:"id"`
	Username       string `json:"username"`
	LastName       string `json:"last_name"`
	FirstName      string `json:"first_name"`
	Email          string `json:"email"`
	Active         bool   `json:"active"`
	EmailConfirmed bool   `json:"email_confirmed"`
}
