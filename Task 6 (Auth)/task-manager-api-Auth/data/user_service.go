package data

import (
	"context"
	"errors"

	"task-manager-api-Auth/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var userCollection  = GetClient().Collection("User")

func CreateUser(ctx context.Context, user *model.User) error {

	
	// Check if this is the first user in the database
	users, _ := userCollection.CountDocuments(context.Background(), bson.M{})
	if users == 0 {
		user.Role = "admin" // Automatically make the first user an admin
	} else {
		user.Role = "user"
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.UserID = primitive.NewObjectID().Hex()

	_, err = userCollection.InsertOne(ctx, user)
	return err
}


func GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user); err != nil {
		return nil, errors.New("invalid username or password")
	}
	return &user, nil
}

func GetUserByID(ctx context.Context, id string) (*model.User, error) {

    // Define a variable to hold the user data
    var user model.User
    if err := userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
        return nil, errors.New("user not found")
    }

    // Return the user object
    return &user, nil
}


func PromoteUser(ctx context.Context, username string) error {

	filter := bson.M{"username": username}

	var user model.User
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return errors.New("user not found")
	}

	update := bson.M{
		"$set": bson.M{
			"role": "admin",
		},
	}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("failed to promote user")
	}

	return nil
}
