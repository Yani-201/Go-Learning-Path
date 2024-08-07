package data

import (
	"errors"
	"task-manager-api/model"
	"time"
)

var tasks = []model.Task{
	{ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"},
	{ID: "2", Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1), Status: "In Progress"},
	{ID: "3", Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2), Status: "Completed"},
}

func GetAllTasks() []model.Task {
	return tasks
}

func GetTaskByID(id string) (model.Task, error) {
	for _, item := range tasks {
		if item.ID == id {
			return item, nil
		}
	}
	return model.Task{}, errors.New("task not found")
}

func UpdateTask(id string, updateTask model.Task) error {
	for k, task := range tasks {
		if task.ID == id {
			if updateTask.Title != "" {
				tasks[k].Title = updateTask.Title
			}
			if !updateTask.DueDate.IsZero() {
				tasks[k].DueDate = updateTask.DueDate
			}
			if updateTask.Status != "" {
				tasks[k].Status = updateTask.Status
			}
			if updateTask.Description != "" {
				tasks[k].Description = updateTask.Description
			}
			return nil
		}
	}
	return errors.New("task not found")
}

func CreateTask(newTask model.Task) {
	tasks = append(tasks, newTask)
}

func DeleteTask(id string) error {
	for k, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:k], tasks[k+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}
