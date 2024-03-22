package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}

var tasks = []Task{
	{ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"},
	{ID: "2", Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1), Status: "In Progress"},
	{ID: "3", Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2), Status: "Completed"},
}

func getTasks(ctx *gin.Context) {
	if len(tasks) == 0 {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "No tasks found"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{"Tasks": tasks})
}

func getTasksbyID(ctx *gin.Context) {
	id := ctx.Param("id")
	for _, item := range tasks {
		if item.ID == id {
			ctx.IndentedJSON(http.StatusOK, item)
			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Task Not Found"})
}

func updateTask(ctx *gin.Context) {
	id := ctx.Param("id")
	updateTask := Task{}

	if err := ctx.BindJSON(&updateTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
			ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Task Updated"})
			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Task not Found"})

}

func createTask(ctx *gin.Context) {
	newTask := Task{}

	if err := ctx.BindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tasks = append(tasks, newTask)
	ctx.IndentedJSON(http.StatusCreated, gin.H{"message": "Task created"})

}

func deleteTask(ctx *gin.Context) {
	id := ctx.Param("id")

	for k, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:k], tasks[k+1:]...)
			ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Task removed"})
			return
		}
	}

	ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Task not Found"})
}

func main() {
	fmt.Println("Task Manger API")
	router := gin.Default()
	router.GET("/tasks", getTasks)
	router.GET("/tasks/:id", getTasksbyID)
	router.PATCH("/tasks/:id/", updateTask)
	// router.PUT("/tasks/:id/", updateTask)
	router.POST("/tasks", createTask)
	router.DELETE("/tasks/:id", deleteTask)
	// Listen and Server in :8080

	router.Run("localhost: 8080")

}
