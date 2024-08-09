package domain

import (
	"context"
	"time"
)

type Task struct {
	TaskID      string    `json:"task_id" bson:"_id"`
	UserID      string    `json:"user_id" bson:"user_id"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	DueDate     time.Time `json:"due_date" bson:"due_date"`
	Status      string    `json:"status" bson:"status"`
	CreatedAt   time.Time `json:"createtimestamp" bson:"createtimestamp"`
	UpdatedAt   time.Time `json:"updatetimestamp" bson:"updatetimestamp"`
}

type TaskRepository interface {
	GetAll(c context.Context, param string) (*[]*Task, error)
	GetByID(c context.Context, taskID string) (*Task, error)
	GetByUserId(c context.Context, userID string) (*[]*Task, error)
	Create(c context.Context, task *Task) (*Task, error)
	Update(c context.Context, task *Task) (*Task, error)
	Delete(c context.Context, taskID string) (*Task, error)

}
