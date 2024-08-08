package data

import (
	"context"
	"errors"
	"task-manager-api-mongoDB/model"
	
	"go.mongodb.org/mongo-driver/bson"

)

func GetAllTasks() ([]model.Task, error) {
	db := GetClient()
	cursor, err := db.Collection("Task").Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var tasks []model.Task
	for cursor.Next(context.Background()) {
		var task model.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTaskByID(id string) (model.Task, error) {
	db := GetClient()
	var task model.Task
	filter := bson.D{{Key: "_id", Value: id}}
	err := db.Collection("Task").FindOne(context.TODO(), filter).Decode(&task)
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func UpdateTask(id string, updateTask model.Task) error {
	db := GetClient()
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{}

	if updateTask.Title != "" {
		update = append(update, bson.E{Key: "$set", Value: bson.D{{Key: "title", Value: updateTask.Title}}})
	}
	if !updateTask.DueDate.IsZero() {
		update = append(update, bson.E{Key: "$set", Value: bson.D{{Key: "due_date", Value: updateTask.DueDate}}})
	}
	if updateTask.Status != "" {
		update = append(update, bson.E{Key: "$set", Value: bson.D{{Key: "status", Value: updateTask.Status}}})
	}
	if updateTask.Description != "" {
		update = append(update, bson.E{Key: "$set", Value: bson.D{{Key: "description", Value: updateTask.Description}}})
	}

	result, err := db.Collection("Task").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("task not found")
	}

	return nil
}

func CreateTask(newTask model.Task) error {
	db := GetClient()
	_, err := db.Collection("Task").InsertOne(context.TODO(), newTask)
	if err != nil {
		return err
	}
	return nil
}

var ErrTaskNotFound = errors.New("task not found")
func DeleteTask(id string) error {
	db := GetClient()
	filter := bson.D{{Key: "_id", Value: id}}
	result, err := db.Collection("Task").DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrTaskNotFound
	}
	

	return nil
}
