package model

// Login contains user's login information.
type Login struct {
	Username string `json:"username" validate:"required"`
}
