package domain

import (
	"context"
)

type User struct {
	UserID    string    `json:"user_id" bson:"_id"`
	Username  string    `json:"username" bson:"username"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"password" bson:"password"`
	Role      string    `json:"role" bson:"role"`
}

type UserInfo struct {
	UserId string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserCreate struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}


type AuthenticatedUser struct {
	UserID   string
	Username string
	Role     string
}


type UserRepository interface {
	Create(c context.Context, user *User) (*User, error)
	GetByUsername(c context.Context,username string) (*User, error)
	UpdateRole(c context.Context, userID string, role string) (*User, error)

}

type UserUseCase interface {
	RegisterUser(c context.Context, payload *UserCreate) (*UserInfo, error)
	Login(c context.Context, payload *UserLogin) (string, error)
	Promote(c context.Context, username string) (*UserInfo, error)
}





