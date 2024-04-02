package model

type AuthenticatedUser struct {
	UserID string
	Username string
	Email string
	Role string
}

type Token struct {
	Token string `json:"token"`
}