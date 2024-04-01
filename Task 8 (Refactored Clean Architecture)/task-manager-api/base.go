package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
)

// Task represents a task entity
type Task struct {
	TaskID      string    `json:"task_id" bson:"_id"`
	UserID      string    `json:"user_id" bson:"user_id"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	DueDate     time.Time `json:"due_date" bson:"due_date"`
	Status      string    `json:"status" bson:"status"`
}

// User represents a user entity
type User struct {
	UserID   string `json:"user_id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Role     string `json:"role" bson:"role"`
}

// TaskRepository provides an interface for interacting with Task data storage
type TaskRepository interface {
	Create(task *Task) error
	GetAll() ([]Task, error)
	GetByID(id string) (*Task, error)
	Update(id string, task *Task) error
	Delete(id string) error
}

// UserRepository provides an interface for interacting with User data storage
type UserRepository interface {
	Create(user *User) error
	GetByUsername(username string) (*User, error)
}

// UseCase represents the business logic layer
type UseCase struct {
	TaskRepo TaskRepository
	UserRepo UserRepository
}

// TaskController handles HTTP requests related to tasks
type TaskController struct {
	uc *UseCase
}

// UserController handles HTTP requests related to users
type UserController struct {
	uc *UseCase
}

// NewTaskController creates a new TaskController instance
func NewTaskController(uc *UseCase) *TaskController {
	return &TaskController{uc: uc}
}

// NewUserController creates a new UserController instance
func NewUserController(uc *UseCase) *UserController {
	return &UserController{uc: uc}
}

// CreateUser handles user registration
func (uc *UserController) CreateUser(ctx *gin.Context) {
	newUser := User{}

	if err := ctx.BindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if user already exists
	existingUser, err := uc.uc.UserRepo.GetByUsername(newUser.Username)
	if err != nil && err != mongo.ErrNoDocuments {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing user"})
		return
	}
	if existingUser != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	newUser.Password = string(hashedPassword)
	userID := primitive.NewObjectID()
	newUser.UserID = userID.Hex()
	newUser.Role = "user"

	err = uc.uc.UserRepo.Create(&newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// CreateAdmin handles admin registration
func (uc *UserController) CreateAdmin(ctx *gin.Context) {
	newAdmin := User{}

	if err := ctx.BindJSON(&newAdmin); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if admin already exists
	existingAdmin, err := uc.uc.UserRepo.GetByUsername(newAdmin.Username)
	if err != nil && err != mongo.ErrNoDocuments {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing admin"})
		return
	}
	if existingAdmin != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Admin already exists"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newAdmin.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	newAdmin.Password = string(hashedPassword)
	userID := primitive.NewObjectID()
	newAdmin.UserID = userID.Hex()
	newAdmin.Role = "admin"

	err = uc.uc.UserRepo.Create(&newAdmin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Admin registered successfully"})
}

// LoginUser handles user login
func (uc *UserController) LoginUser(ctx *gin.Context) {
	user := User{}

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the user from the database based on the provided username
	existingUser, err := uc.uc.UserRepo.GetByUsername(user.Username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		}
		return
	}

	// Verify the password using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  existingUser.UserID,
		"username": existingUser.Username,
		"role":     existingUser.Role,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // Token expires in 1 week
	})

	jwtToken, err := token.SignedString(jwtSecret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User logged in successfully", "token": jwtToken})
}

// CreateTask handles task creation
func (tc *TaskController) CreateTask(ctx *gin.Context) {
	newTask := Task{}

	// Get authenticated user from Gin context
	authenticatedUser, exists := ctx.Get("AuthenticatedUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user, ok := authenticatedUser.(*AuthenticatedUser)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get authenticated user"})
		return
	}

	// Bind JSON data to newTask
	if err := ctx.BindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	newTask.UserID = user.UserID
	taskID := primitive.NewObjectID()
	newTask.TaskID = taskID.Hex()

	err := tc.uc.TaskRepo.Create(&newTask)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	// Return success message
	ctx.JSON(http.StatusCreated, gin.H{"message": "Task created"})
}

// GetTasks retrieves all tasks
func (tc *TaskController) GetTasks(ctx *gin.Context) {
	tasks, err := tc.uc.TaskRepo.GetAll()
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

// GetTaskByID retrieves a task by ID
func (tc *TaskController) GetTaskByID(ctx *gin.Context) {
	id := ctx.Param("id")

	task, err := tc.uc.TaskRepo.GetByID(id)
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

// UpdateTask updates a task
func (tc *TaskController) UpdateTask(ctx *gin.Context) {
	id := ctx.Param("id")
	updateTask := Task{}

	// Get authenticated user from Gin context
	authenticatedUser, exists := ctx.Get("AuthenticatedUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user, ok := authenticatedUser.(*AuthenticatedUser)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get authenticated user"})
		return
	}

	// Retrieve the task from the database
	existingTask, err := tc.uc.TaskRepo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		}
		return
	}

	// Check if the authenticated user is the creator of the task or an admin
	if user.Role != "admin" && existingTask.UserID != user.UserID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Only the task creator or admin can update this task"})
		return
	}

	// Bind JSON data to updateTask
	if err := ctx.BindJSON(&updateTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the task in the database
	err = tc.uc.TaskRepo.Update(id, &updateTask)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Task Updated"})
}

// DeleteTask deletes a task
func (tc *TaskController) DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")

	// Get authenticated user from Gin context
	authenticatedUser, exists := ctx.Get("AuthenticatedUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user, ok := authenticatedUser.(*AuthenticatedUser)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get authenticated user"})
		return
	}

	// Retrieve the task from the database
	existingTask, err := tc.uc.TaskRepo.GetByID(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		}
		return
	}

	// Check if the authenticated user is the creator of the task or an admin
	if user.Role != "admin" && existingTask.UserID != user.UserID {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Only the task creator or admin can delete this task"})
		return
	}

	// Delete the task from the database
	err = tc.uc.TaskRepo.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Task removed"})
}

func main() {
	fmt.Println("Task Manager API")

	// Initialize MongoDB client and repositories
	mongoClient := GetClient()
	taskRepo := NewMongoTaskRepository(mongoClient.Database("task_manager").Collection("Task"))
	userRepo := NewMongoUserRepository(mongoClient.Database("task_manager").Collection("User"))

	// Initialize use case
	uc := &UseCase{TaskRepo: taskRepo, UserRepo: userRepo}

	// Initialize controllers
	taskController := NewTaskController(uc)
	userController := NewUserController(uc)

	router := gin.Default()
	router.POST("/register", userController.CreateUser)
	router.POST("/register/admin", userController.CreateAdmin)
	router.POST("/login", userController.LoginUser)

	authMiddleware := AuthMiddleware()

	taskRoutes := router.Group("/tasks")
	{
		taskRoutes.Use(authMiddleware)
		taskRoutes.POST("/", taskController.CreateTask)
		taskRoutes.GET("/", taskController.GetTasks)
		taskRoutes.GET("/:id", taskController.GetTaskByID)
		taskRoutes.PATCH("/:id/", taskController.UpdateTask)
		taskRoutes.DELETE("/:id", taskController.DeleteTask)
	}

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
