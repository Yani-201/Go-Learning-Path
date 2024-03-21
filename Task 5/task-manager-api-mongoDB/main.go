package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}

// var tasks = []Task{
// 	{ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"},
// 	{ID: "2", Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1), Status: "In Progress"},
// 	{ID: "3", Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2), Status: "Completed"},
// }

func getTasks(ctx *gin.Context) {
	db := GetClient()
	cursor, err := db.Collection("Task").Find(context.TODO(), bson.D{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}
	defer cursor.Close(context.TODO())

	var tasks []Task
	for cursor.Next(context.Background()) {
		var task Task
		err := cursor.Decode(&task)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode task"})
			return
		}
		tasks = append(tasks, task)
	}

	if err := cursor.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error"})
		return
	}

	if len(tasks) == 0 {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "No tasks found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"Tasks": tasks})
}

func getTasksByID(ctx *gin.Context) {
	db := GetClient()
	id := ctx.Param("id")

	var task Task
	filter := bson.D{{Key: "id", Value: id}}
	err := db.Collection("Task").FindOne(context.TODO(), filter).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, task)
}

func updateTask(ctx *gin.Context) {
	db := GetClient()
	id := ctx.Param("id")
	updateTask := Task{}

	if err := ctx.BindJSON(&updateTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := bson.D{{Key: "id", Value: id}}
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	if result.ModifiedCount == 0 {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Task Updated"})

}

func createTask(ctx *gin.Context) {
	db := GetClient()
	newTask := Task{}

	if err := ctx.BindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Collection("Task").InsertOne(context.TODO(), newTask)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	// tasks = append(tasks, newTask)
	ctx.IndentedJSON(http.StatusCreated, gin.H{"message": "Task created"})

}

func deleteTask(ctx *gin.Context) {
	db := GetClient()
	id := ctx.Param("id")

	filter := bson.D{{Key: "id", Value: id}}
	result, err := db.Collection("Task").DeleteOne(context.TODO(), filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	if result.DeletedCount == 0 {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Task removed"})
}

func main() {
	fmt.Println("Task Manger API")
	router := gin.Default()
	router.GET("/tasks", getTasks)
	router.GET("/tasks/:id", getTasksByID)
	router.PUT("/tasks/:id/", updateTask)
	router.POST("/tasks", createTask)
	router.DELETE("/tasks/:id", deleteTask)
	// Listen and Server in :8080

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}

}
