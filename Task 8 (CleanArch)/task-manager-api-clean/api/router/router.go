package router

import (
	"task-manager-api-clean/config"
	"task-manager-api-clean/api/controller"
	"task-manager-api-clean/api/middleware"
	"task-manager-api-clean/repository"
	"task-manager-api-clean/usecase"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(env *config.Environment, db *mongo.Database, gin *gin.Engine) {
	// Initialize repositories
	userRepository := repository.NewUserRepository(db, "users")
	taskRepository := repository.NewTaskRepository(db, "tasks")

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepository, env)
	taskUseCase := usecase.NewTaskUseCase(taskRepository, userRepository, env)

	// Initialize controllers
	userController := controller.NewUserController(userUseCase)
	taskController := controller.NewTaskController(taskUseCase)

	// Middleware
	authMiddleware := middleware.AuthMiddleware(env.JwtSecret)

	// User routes
	userRouter := gin.Group("")
	{
		userRouter.POST("/register", userController.CreateUser)
		userRouter.POST("/login", userController.LoginUser)
		userRouter.POST("/promote/:username", authMiddleware, userController.PromoteUser)
	}

	// Task routes
	taskRouter := gin.Group("tasks")
	taskRouter.Use(authMiddleware)
	{
		taskRouter.GET("/", taskController.GetTasks)
		taskRouter.GET("/:id", taskController.GetTaskByID)
		taskRouter.PATCH("/:id", taskController.UpdateTask)
		taskRouter.POST("/", taskController.CreateTask)
		taskRouter.DELETE("/:id", taskController.DeleteTask)
	}
}

