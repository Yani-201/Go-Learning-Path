package controller

import (
	"net/http"
	"task-manager-api/data"
	"task-manager-api/model"

	"github.com/gin-gonic/gin"
)

func GetTasks(ctx *gin.Context) {
	tasks := data.GetAllTasks()
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
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Task Not Found"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, task)
}

func UpdateTask(ctx *gin.Context) {
	id := ctx.Param("id")
	updateTask := model.Task{}

	if err := ctx.BindJSON(&updateTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := data.UpdateTask(id, updateTask)
	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Task not Found"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Task Updated"})
}

func CreateTask(ctx *gin.Context) {
	newTask := model.Task{}

	if err := ctx.BindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data.CreateTask(newTask)
	ctx.IndentedJSON(http.StatusCreated, gin.H{"message": "Task created"})
}

func DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")

	err := data.DeleteTask(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Task not Found"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Task removed"})
}

