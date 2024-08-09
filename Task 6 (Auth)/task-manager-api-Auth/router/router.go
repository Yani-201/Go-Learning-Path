package router

import (
	"task-manager-api-Auth/controller"
	"task-manager-api-Auth/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	publicRoutes := r.Group("")
	{
		publicRoutes.POST("/register", controller.CreateUser)
		publicRoutes.POST("/login", controller.LoginUser)
	}

	protectedRoutes := r.Group("")
	protectedRoutes.Use(middleware.AuthMiddleware())
	{
		protectedRoutes.POST("/tasks", controller.CreateTask)
		protectedRoutes.GET("/tasks", controller.GetTasks)
		protectedRoutes.GET("/tasks/:id", controller.GetTaskByID)
		protectedRoutes.PATCH("/tasks/:id", controller.UpdateTask)
		protectedRoutes.DELETE("/tasks/:id", controller.DeleteTask)
		protectedRoutes.POST("/promote/:username", controller.PromoteUser) 
	}

	return r
}




