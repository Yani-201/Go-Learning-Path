package model

type User struct {
	UserID   string `json:"user_id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Email string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Role     string `json:"role" bson:"role"`
}

type AuthenticatedUser struct {
	UserID   string
	Username string
	Role     string
}
