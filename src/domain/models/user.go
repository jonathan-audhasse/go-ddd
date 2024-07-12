package models

// User domain model
type User struct {
	// ID
	Id string `json:"id" db:"id"`
	// name
	Username string `json:"username" db:"username"`
	// email
	Password string `json:"password" db:"password"`
}
