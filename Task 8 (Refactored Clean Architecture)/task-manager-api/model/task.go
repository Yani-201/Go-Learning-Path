package model

import "time"

type TaskInfo struct {
	TaskID      string    `json:"task_id" bson:"_id"`
	UserID      string    `json:"user_id" bson:"user_id"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	DueDate     time.Time `json:"due_date" bson:"due_date"`
	Status      string    `json:"status" bson:"status"`
}

type TaskCreate struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      bool      `json:"status"`
	DueDate     time.Time `json:"dueDate"`
}

type TaskUpdate struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      bool      `json:"status"`
	DueDate     time.Time `json:"dueDate"`
}