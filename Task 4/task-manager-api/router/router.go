package router

import (
	"task-manager-api/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/tasks", controller.GetTasks)
	router.GET("/tasks/:id", controller.GetTaskByID)
	router.PATCH("/tasks/:id", controller.UpdateTask)
	router.POST("/tasks", controller.CreateTask)
	router.DELETE("/tasks/:id", controller.DeleteTask)

	return router
}
