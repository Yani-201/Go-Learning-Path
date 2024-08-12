package repository

import (
	"context"
	"errors"

	"task-manager-api-clean/utils"
	"task-manager-api-clean/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

)

type UserRepository struct {
	database   *mongo.Database
	collection string
}

func NewUserRepository(db *mongo.Database, collection string) domain.UserRepository {
	return &UserRepository{
		database:   db,
		collection: collection,
		}
}

func (ur *UserRepository) Create(c context.Context, user *domain.User) (*domain.User, error) {
	// Check if this is the first user in the database
	users, _ := ur.database.Collection(ur.collection).CountDocuments(context.Background(), bson.M{})
	if users == 0 {
		user.Role = "admin" // Automatically make the first user an admin
	} else {
		user.Role = "user"
	}

	// hash the password
	hashedPassword, err := utils.EncryptPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword
	user.UserID = primitive.NewObjectID().Hex()

	_, err = ur.database.Collection(ur.collection).InsertOne(c, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) GetByUsername(c context.Context, username string) (*domain.User, error) {
	var user domain.User
	if err := ur.database.Collection(ur.collection).FindOne(c, bson.M{"username": username}).Decode(&user); err != nil {
		return nil, errors.New("invalid username or password")
	}
	return &user, nil
}

func (ur *UserRepository) UpdateRole(c context.Context, userID string, role string) (*domain.User, error) {
	filter := bson.M{"username": userID}

	var user domain.User
	err := ur.database.Collection(ur.collection).FindOne(c, filter).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}

	update := bson.M{
		"$set": bson.M{
			"role": role,
		},
	}

	_, err = ur.database.Collection(ur.collection).UpdateOne(c, filter, update)
	if err != nil {
		return nil, errors.New("failed to update user role")
	}

	return &user, nil
}
