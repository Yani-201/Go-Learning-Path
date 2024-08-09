package data

import (
	"context"
	"errors"
	"task-manager-api-Auth/model"
	"go.mongodb.org/mongo-driver/bson"
)

var taskCollection = GetClient().Collection("Task")

func CreateTask(ctx context.Context, task *model.Task) error {
	_, err := taskCollection.InsertOne(ctx, task)
	return err
}

func GetTasks(ctx context.Context) ([]model.Task, error) {
	cursor, err := taskCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []model.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTaskByID(ctx context.Context, id string) (*model.Task, error) {

	var task model.Task
	if err := taskCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&task); err != nil {
		return nil, err
	}
	return &task, nil
}

func UpdateTask(ctx context.Context, id string, updateTask *model.Task) (*model.Task, error) {

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
	
	if len(setFields) >= 0 {
		update["$set"] = setFields
	}

    result, err := taskCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
    if err != nil {
        return nil, err
    }

    if result.ModifiedCount == 0 {
        return nil, errors.New("task not updated, no new information is provided")
    }

    var updated model.Task
    if err = taskCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&updated); err != nil {
        return nil, err
    }

    return &updated, nil
}

func DeleteTask(ctx context.Context, id string) error {

	_, err := taskCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
        return err
    }

    return nil
}

