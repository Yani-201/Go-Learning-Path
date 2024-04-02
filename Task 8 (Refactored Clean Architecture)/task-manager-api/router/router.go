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
	taskUseCase := usecase.NewTaskUseCase(&ctx, env, &taskRepository, &userRepository)
	authUseCase := usecase.NewAuthUseCase(&ctx, env, &userRepository)

	userController := controller.NewUserController(env, &userUseCase)
	taskController :=  controller.NewTaskController(env, &taskUseCase)
	authController := controller.NewAuthController(env, &authUseCase)

	publicRouter := gin.Group("auth")
	publicRouter.POST("/register", authController.RegisterUser)
	publicRouter.POST("/login", authController.Login)
	publicRouter.POST("/adminRegister", authController.RegisterAdmin)

	taskRouter := gin.Group("task")
	taskRouter.GET("/", middleware.AuthMiddleware(env.JwtSecret), taskController.getTasks)
	taskRouter.GET("/:id", middleware.AuthMiddleware(env.JwtSecret), taskController.getTasksByID)
	taskRouter.PATCH("/:id", middleware.AuthMiddleware(env.JwtSecret), taskController.updateTask)
	taskRouter.POST("/", middleware.AuthMiddleware(env.JwtSecret), taskController.createTask)
	taskRouter.DELETE(":id", middleware.AuthMiddleware(env.JwtSecret), taskController.deleteTask)

	userRouter := gin.Group("user")
	userRouter.GET("/", middleware.AuthMiddleware(env.JwtSecret), userController.getUsers)
	userRouter.GET("/:id", middleware.AuthMiddleware(env.JwtSecret), userController.getUserByID)
	userRouter.PATCH("/updateuser", middleware.AuthMiddleware(env.JwtSecret), userController.updateUser)
	userRouter.PATCH("/updatepassword", middleware.AuthMiddleware(env.JwtSecret), userController.updatePassword)
	userRouter.DELETE("/:id", middleware.AuthMiddleware(env.JwtSecret), userController.deleteUser)


	// blogRouter.GET("/", blogController.GetAllBlogs)
	// blogRouter.GET("/:blog_id", blogController.GetByBlogID)
	// blogRouter.POST("/",middleware.AuthMiddleware(env.JwtSecret), blogController.CreateBlog)
	// blogRouter.PUT("/:blog_id",middleware.AuthMiddleware(env.JwtSecret), blogController.UpdateBlog)
	// blogRouter.DELETE("/:blog_id",middleware.AuthMiddleware(env.JwtSecret), blogController.DeleteBlog)


}
