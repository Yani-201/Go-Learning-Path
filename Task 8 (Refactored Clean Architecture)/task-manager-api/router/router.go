package router



import (
	"task-manager-api/config"
	"task-manager-api/controller"
	"task-manager-api/middleware"
	"task-manager-api/repository"
	"task-manager-api/usecase"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(env *config.Environment, timeout time.Duration, db *mongo.Database, gin *gin.Engine) {
	ctx := context.TODO()
	userRepository := repository.NewUserRepository(db, "users")
	taskRepository := repository.NewTaskRepository(db, "tasks")

	userUseCase := usecase.NewUserUseCase(&ctx, env, &userRepository )
	taskUsecase := usecase.NewTaskUseCase(&ctx, env, &taskRepository)

	userController := controller.NewUserController(env, &userUseCase)
	taskController :=  controller.NewTaskController(env, &taskUsecase)

	publicRouter := gin.Group("auth")
	publicRouter.POST("/register", userController.Register)
	publicRouter.POST("/login", userController.Login)
	publicRouter.POST("/adminRegister", userController.AdminRegister)

	taskRouter := gin.Group("task")
	taskRouter.GET("/", middleware.AuthMiddleware(env.JwtSecret), taskController.getTasks)
	taskRouter.GET("/:id", middleware.AuthMiddleware(env.JwtSecret), taskController.getTasksByID)
	taskRouter.PATCH("/:id", middleware.AuthMiddleware(env.JwtSecret), taskController.updateTask)
	taskRouter.POST("/", middleware.AuthMiddleware(env.JwtSecret), taskController.createTask)
	taskRouter.DELETE(":id", middleware.AuthMiddleware(env.JwtSecret), taskController.deleteTask)

	userRouter := gin.Group("user")
	userRouter.GET("/", middleware.AuthMiddleware(env.JwtSecret), userController.getUsers)
	userRouter.GET("/:id", middleware.AuthMiddleware(env.JwtSecret), userController.getUserByID)
	userRouter.PATCH("/updateuser", middleware.AuthMiddleware(env.JwtSecret), userController.updateUserInfo)
	userRouter.PATCH("/updatepassword", middleware.AuthMiddleware(env.JwtSecret), userController.updatePassword)
	userRouter.DELETE("/:id", middleware.AuthMiddleware(env.JwtSecret), userController.deleteUser)




}
