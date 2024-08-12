package controller

import (
	"net/http"
	"go.mongodb.org/mongo-driver/mongo"
	"task-manager-api-clean/domain"
	"task-manager-api-clean/utils"
	"github.com/gin-gonic/gin"
)

type TaskController struct {
	taskUseCase domain.TaskUseCase
}

func NewTaskController(taskUseCase domain.TaskUseCase) *TaskController {
	return &TaskController{
		taskUseCase: taskUseCase,
	}
}

func (tc *TaskController) CreateTask(ctx *gin.Context) {
	var newTask domain.TaskInput

	// Get authenticated user from gin context
	user, err := utils.CheckUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if user.Role != "admin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Only the admin can create a task"})
		return
	}
	// Bind JSON data to newTask
	if err := ctx.BindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	task, err := tc.taskUseCase.Create(ctx, &newTask)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Task created", "task": task})
}

func (tc *TaskController) GetTasks(ctx *gin.Context) {
	// Get authenticated user from gin context
	_, err := utils.CheckUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	tasks, err := tc.taskUseCase.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	if len(*tasks) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No tasks found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (tc *TaskController) GetTaskByID(ctx *gin.Context) {
	id := ctx.Param("id")

	// Get authenticated user from gin context
	_, err := utils.CheckUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}


	task, err := tc.taskUseCase.GetById(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		return
	}

	ctx.JSON(http.StatusOK, task)
}

func (tc *TaskController) UpdateTask(ctx *gin.Context) {
	id := ctx.Param("id")
	var updateTask domain.TaskInput

	// Get authenticated user from gin context
	user, errr := utils.CheckUser(ctx)
	if errr != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errr.Error()})
		return
	}

	if user.Role != "admin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Only the admin can update this task"})
		return
	}

	_, err := tc.taskUseCase.GetById(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		return
	}

	if err := ctx.BindJSON(&updateTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	task, err := tc.taskUseCase.Update(ctx, id, &updateTask)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, task)
}

func (tc *TaskController) DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")

	// Get authenticated user from gin context
	user, errr := utils.CheckUser(ctx)
	if errr != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errr.Error()})
		return
	}


	if user.Role != "admin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Only the admin can delete this task"})
		return
	}

	_, err := tc.taskUseCase.GetById(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		return
	}

	if err := tc.taskUseCase.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task removed"})
}
