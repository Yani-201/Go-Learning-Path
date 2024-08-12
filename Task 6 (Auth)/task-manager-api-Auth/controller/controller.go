package controller

import (
	"net/http"
	"time"

	"task-manager-api-Auth/model"
	"task-manager-api-Auth/data"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("your_jwt_secret_is_Yanet")

func CreateUser(ctx *gin.Context) {
	newUser := model.User{}
// Bind JSON to new user
	if err := ctx.BindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	//check if user exists
	_, err := data.GetUserByUsername(ctx, newUser.Username)
	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	if err := data.CreateUser(ctx, &newUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}


func LoginUser(ctx *gin.Context) {
	user := model.User{}

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Retrive use from database
	retrievedUser, err := data.GetUserByUsername(ctx, user.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// verify the password
	if err := bcrypt.CompareHashAndPassword([]byte(retrievedUser.Password), []byte(user.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	//Generate JWT Token 
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  retrievedUser.UserID,
		"username": retrievedUser.Username,
		"role":     retrievedUser.Role,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	jwtToken, err := token.SignedString(jwtSecret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User logged in successfully", "token": jwtToken})
}


func PromoteUser(ctx *gin.Context) {
	//Get authenticated user from gin context
	authenticatedUser, exists := ctx.Get("AuthenticatedUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user, ok := authenticatedUser.(*model.AuthenticatedUser)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get authenticated user"})
		return
	}

	//Retrive use from database
	_, er := data.GetUserByID(ctx, user.UserID)
	if er != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	if user.Role != "admin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Only the admin can promote a user"})
		return
	}
	username := ctx.Param("username")
	

	err := data.PromoteUser(ctx, username)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User promoted to admin"})
}


func CreateTask(ctx *gin.Context) {
	newTask := model.Task{}

	//Get authenticated user from gin context
	authenticatedUser, exists := ctx.Get("AuthenticatedUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user, ok := authenticatedUser.(*model.AuthenticatedUser)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get authenticated user"})
		return
	}

	//Retrive use from database
	_, er := data.GetUserByID(ctx, user.UserID)
	if er != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	if user.Role != "admin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Only the admin can create a task"})
		return
	}
	//Bind JSON data to newTask
	if err := ctx.BindJSON(&newTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	newTask.UserID = user.UserID
	taskID := primitive.NewObjectID()
	newTask.TaskID = taskID.Hex()

	if err := data.CreateTask(ctx, &newTask); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Task created"})
}

func GetTasks(ctx *gin.Context) {
	authenticatedUser, exists := ctx.Get("AuthenticatedUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, ok := authenticatedUser.(*model.AuthenticatedUser)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get authenticated user"})
		return
	}


	//Retrive use from database
	_, er := data.GetUserByID(ctx, user.UserID)
	if er != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	tasks, err := data.GetTasks(ctx)
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

	authenticatedUser, exists := ctx.Get("AuthenticatedUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, ok := authenticatedUser.(*model.AuthenticatedUser)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get authenticated user"})
		return
	}


	//Retrive use from database
	_, er := data.GetUserByID(ctx, user.UserID)
	if er != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	task, err := data.GetTaskByID(ctx, id)
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
	updateTask := model.Task{}
//Get authenticated user fron gin context
	authenticatedUser, exists := ctx.Get("AuthenticatedUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user, ok := authenticatedUser.(*model.AuthenticatedUser)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get authenticated user"})
		return
	}
	//Retrive use from database
	_, er := data.GetUserByID(ctx, user.UserID)
	if er != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	if user.Role != "admin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Only the admin can update this task"})
		return
	}
	
	_, err := data.GetTaskByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		return
	}

	if err := ctx.BindJSON(&updateTask); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	task, err := data.UpdateTask(ctx, id, &updateTask)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
	ctx.IndentedJSON(http.StatusOK, task)
}

func DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")

	//Get authenticated user from gin context
	authenticatedUser, exists := ctx.Get("AuthenticatedUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user, ok := authenticatedUser.(*model.AuthenticatedUser)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get authenticated user"})
		return
	}

	//Retrive use from database
	_, er := data.GetUserByID(ctx, user.UserID)
	if er != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	if user.Role != "admin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Only the admin can delete this task"})
		return
	}

	_, err := data.GetTaskByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		return
	}

	if err := data.DeleteTask(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Task removed"})
}
