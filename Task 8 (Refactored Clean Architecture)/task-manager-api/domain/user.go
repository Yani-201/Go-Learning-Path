package domain

import (
	"context"
	"time"
)

type User struct {
	UserID    string    `json:"user_id" bson:"_id"`
	Username  string    `json:"username" bson:"username"`
	Email     string    `json:"email" bson:"email"`
	Name 	string    `json:"name" bson:"name"`
	Password  string    `json:"password" bson:"password"`
	Role      string    `json:"role" bson:"role"`
	CreatedAt time.Time `json:"timestamp" bson:"timestamp"`
}

type UserRepository interface {
	Create(c context.Context, user *User) (*User, error)
	Delete(c context.Context, user *User) (*User, error)
	UpdatePassword(c context.Context, user *User) (*User, error)
	UpdateUser(c context.Context, user *User) (*User, error)
	GetByUsername(c context.Context,username string) (*User, error)
	GetByEmail(c context.Context, email string) (*User, error)
	GetById(c context.Context, id string) (*User, error)
	GetAll(c context.Context, param string) (*[]*User, error)

}

