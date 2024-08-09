package repository

import (
	"BlogApp/domain"
	"context"
	"time"

	// "time"

	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlogRepository struct {
	database   *mongo.Database
	collection string
}

func NewBlogRepository(db *mongo.Database, collection string) domain.BlogRepository {
	return &BlogRepository{
		database:   db,
		collection: collection,
	}
}

// Create implements domain.BlogRepository.
func (b *BlogRepository) Create(c context.Context, blog *domain.Blog) (*domain.Blog, error) {
	blog.BlogID = primitive.NewObjectID().Hex()
	_, err := b.database.Collection(b.collection).InsertOne(c, *blog)
	if err != nil{ 
	  return nil, err
	}
	return blog, nil
}

// GetAll implements domain.BlogRepository.
func (b *BlogRepository) GetAll(c context.Context, param string) (*[]*domain.Blog, error) {
	filter := bson.M{
        "$or": []bson.M{
            {"title": bson.M{"$regex": primitive.Regex{Pattern: param, Options: "i"}}},
            {"content": bson.M{"$regex": primitive.Regex{Pattern: param, Options: "i"}}},
        },
    }

    // Perform the find operation
    cursor, err := b.database.Collection(b.collection).Find(c, filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(c)

    // Iterate through the cursor and decode each document into a User struct
    var blogs []*domain.Blog
    for cursor.Next(c) {
        var blog domain.Blog
        if err := cursor.Decode(&blog); err != nil {
			cursor.Close(c)
            return nil, err
        }
        blogs = append(blogs, &blog)
    }

    // Check if any error occurred during cursor iteration
    if err := cursor.Err(); err != nil {
        return nil, err
    }

    return &blogs, nil
	
}

// GetByID implements domain.BlogRepository.
func (b *BlogRepository) GetByID(c context.Context, blogID string) (*domain.Blog, error) {
	filter := bson.M{"_id": blogID}
	var blog domain.Blog
	err := b.database.Collection(b.collection).FindOne(c, filter).Decode(&blog)
	if err != nil {
		return nil, err
	}

	return &blog, nil
}

// GetByUserId implements domain.BlogRepository.
func (b *BlogRepository) GetByUserId(c context.Context, userID string) (*[]*domain.Blog, error) {
	filter := bson.M{"user_id": userID}

	// Perform the find operation
	cursor, err := b.database.Collection(b.collection).Find(c, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c)

	// Iterate through the cursor and decode each document into a User struct
	var blogs []*domain.Blog
	for cursor.Next(c) {
		var blog domain.Blog
		if err := cursor.Decode(&blog); err != nil {
			cursor.Close(c)
			return nil, err
		}
		blogs = append(blogs, &blog)
	}

	// Check if any error occurred during cursor iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &blogs, nil
}

// Update implements domain.BlogRepository.
func (b *BlogRepository) Update(c context.Context, blog *domain.Blog) (*domain.Blog, error) {
	filter := bson.M{"_id": blog.BlogID}
	update := bson.M{"$set": bson.M{"title": blog.Title, "content": blog.Content, "updatetimestamp": time.Now()}}
	_, err := b.database.Collection(b.collection).UpdateOne(c, filter, update)
	if err != nil {
		return nil, err
	}
	return blog, nil
}


// Delete implements domain.BlogRepository.
func (b *BlogRepository) Delete(c context.Context, blogID string) (*domain.Blog, error) {
	filter := bson.M{"_id": blogID}

	var deletedBlog domain.Blog
	err := b.database.Collection(b.collection).FindOneAndDelete(c, filter).Decode(&deletedBlog)
	if err != nil {
		return nil, err
	}

	return &deletedBlog, nil
}

