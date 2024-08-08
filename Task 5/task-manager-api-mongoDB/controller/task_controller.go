package controller

import (
	"errors"
	"net/http"
	"task-manager-api-mongoDB/data"
	"task-manager-api-mongoDB/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetTasks(ctx *gin.Context) {
	tasks, err := data.GetAllTasks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	if len(tasks) == 0 {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "No tasks found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"Tasks": tasks})
}

func GetTaskByID(ctx *gin.Context) {
	id := ctx.Param("id")
	task, err := data.GetTaskByID(id)
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

func UpdateTask(ctx *gin.Context) {
	id := ctx.Param("id")
	var updateTask model.Task

	if err := ctx.BindJSON(&updateTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := data.UpdateTask(id, updateTask)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Task Updated"})
}

func CreateTask(ctx *gin.Context) {
	var newTask model.Task

	if err := ctx.BindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	taskID := primitive.NewObjectID()
    newTask.ID = taskID.Hex()

	err := data.CreateTask(newTask)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, gin.H{"message": "Task created"})
}

func DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")

	err := data.DeleteTask(id)
	if err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) || errors.Is(err, data.ErrTaskNotFound) {
            ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Task not found"})
            return
        }
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Task removed"})
}
