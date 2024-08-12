package repository

import (
	"context"
	"errors"
	"task-manager-api-clean/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskRepository struct {
	database   *mongo.Database
	collection string
}

func NewTaskRepository(db *mongo.Database, collection string) domain.TaskRepository {
	return &TaskRepository{
		database:   db,
		collection: collection,
	}
}

func (repo *TaskRepository) Create(c context.Context, task *domain.Task) (*domain.Task, error) {
	task.Id = primitive.NewObjectID().Hex()
	_, err := repo.database.Collection(repo.collection).InsertOne(c, task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (repo *TaskRepository) Update(c context.Context, id string, updateTask *domain.Task) (*domain.Task, error) {


	update := bson.M{}
	setFields := bson.M{}
	
	if updateTask.Title != "" {
		setFields["title"] = updateTask.Title
	}
	if !updateTask.DueDate.IsZero() {
		setFields["due_date"] = updateTask.DueDate
	}
	if updateTask.Status != "" {
		setFields["status"] = updateTask.Status
	}
	if updateTask.Description != "" {
		setFields["description"] = updateTask.Description
	}
	
	if len(setFields) > 0 {
		update["$set"] = setFields
	}

    result, err := repo.database.Collection(repo.collection).UpdateOne(c, bson.M{"_id": id}, update)
    if err != nil {
        return nil, err
    }

    if result.ModifiedCount == 0 {
        return nil, errors.New("task not updated, no new information is provided")
    }

    var updated domain.Task
    if err = repo.database.Collection(repo.collection).FindOne(c, bson.M{"_id": id}).Decode(&updated); err != nil {
        return nil, err
    }

    return &updated, nil
}

func (repo *TaskRepository) Delete(c context.Context, id string) error {
	_, err := repo.database.Collection(repo.collection).DeleteOne(c, bson.M{"_id": id})
	if err != nil {
        return err
    }

    return nil
}

func (repo *TaskRepository) GetAll(c context.Context) (*[]*domain.Task, error) {
	cursor, err := repo.database.Collection(repo.collection).Find(c, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c)

	var tasks []*domain.Task
	if err = cursor.All(c, &tasks); err != nil {
		return nil, err
	}
	return &tasks, nil
}

func (repo *TaskRepository) GetById(c context.Context, id string) (*domain.Task, error) {

	var task domain.Task
	if err := repo.database.Collection(repo.collection).FindOne(c, bson.M{"_id": id}).Decode(&task); err != nil {
		return nil, err
	}
	return &task, nil
}

