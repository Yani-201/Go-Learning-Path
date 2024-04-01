package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"strings"

	"task-manager-api/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
)

type Task struct {
	TaskID  string    `json:"task_id" bson:"_id"`
	UserID 	string    `json:"user_id" bson:"user_id"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	DueDate     time.Time `json:"due_date" bson:"due_date"`
	Status      string    `json:"status" bson:"status"`
}

type User struct {
	UserID string `json:"user_id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Email string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Role string `json:"role" bson:"role"`
}

func createUser(ctx *gin.Context) {

		db := GetClient()
		newUser := User{}
	
		// Bind JSON data to newUser
		if err := ctx.BindJSON(&newUser); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
	
		// Check if user already exists
		filter := bson.D{{Key: "username", Value: newUser.Username}}
		existingUser := db.Collection("User").FindOne(context.TODO(), filter)
		if existingUser.Err() == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
			return
		} else if existingUser.Err() != nil && existingUser.Err() != mongo.ErrNoDocuments {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing user"})
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

		_, err = db.Collection("User").InsertOne(context.TODO(), newUser)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	
		ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
	}

func createAdmin(ctx *gin.Context) {

		db := GetClient()
		newAdmin := User{}

		if err := ctx.BindJSON(&newAdmin); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
	
		// Check if admin already exists
		filter := bson.D{{Key: "username", Value: newAdmin.Username}}
		existingAdmin := db.Collection("User").FindOne(context.TODO(), filter)
		if existingAdmin.Err() == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Admin already exists"})
			return
		} else if existingAdmin.Err() != nil && existingAdmin.Err() != mongo.ErrNoDocuments {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing admin"})
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
	
		// Insert the new admin user into the database
		_, err = db.Collection("User").InsertOne(context.TODO(), newAdmin)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin"})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Admin registered successfully"})
	}

var jwtSecret = []byte("your_jwt_secret_is_Yanet")

func loginUser(ctx *gin.Context) {
	db := GetClient()
	user := User{}

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the user from the database based on the provided username
	filter := bson.D{{Key: "username", Value: user.Username}}
	existingUser := db.Collection("User").FindOne(context.Background(), filter)
	if existingUser.Err() != nil {
		if existingUser.Err() == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		}
		return
	}

	// Verify the password using bcrypt
	var retrievedUser User
	if err := existingUser.Decode(&retrievedUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(retrievedUser.Password), []byte(user.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": retrievedUser.UserID,
		"username": retrievedUser.Username,
		"role": retrievedUser.Role,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // Token expires in 1 week
	})

	jwtToken, err := token.SignedString(jwtSecret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User logged in successfully", "token": jwtToken})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			c.Abort()
			return
		}

			token, err := jwt.Parse(authParts[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT"})
			c.Abort()
			return
		}

			// Extract claims from token
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				return
			}
	
			// Set user ID and username in the Gin context
			UserID, ok := claims["user_id"].(string)

			if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
			return
		}

			Username, ok := claims["username"].(string)

			
			if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Username not found in token"})
			return
		}

			Role, ok := claims["role"].(string)
	
	
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "role not found in token"})
				return
			}
			c.Set("AuthenticatedUser", &model.AuthenticatedUser{
				UserID: UserID,
				Username: Username,
				Role: Role,
			})
	
			c.Next()
		}

	
}

func createTask(ctx *gin.Context) {
    db := GetClient()
    newTask := Task{}

    // Get authenticated user from Gin context
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

    // Bind JSON data to newTask
    if err := ctx.BindJSON(&newTask); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    newTask.UserID = user.UserID
    taskID := primitive.NewObjectID()
    newTask.TaskID = taskID.Hex()

    // Insert the new task into the database
    _, err := db.Collection("Task").InsertOne(context.TODO(), newTask)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
        return
    }

    // Return success message
    ctx.JSON(http.StatusCreated, gin.H{"message": "Task created"})
}

func getTasks(ctx *gin.Context) {
    db := GetClient()
	fmt.Println("+++++++++++++++++++++++++++++++++++")
    _, exists := ctx.Get("AuthenticatedUser")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    // Retrieve tasks from the database
    cursor, err := db.Collection("Task").Find(context.TODO(), bson.D{})
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
        return
    }
    defer cursor.Close(context.TODO())

    var tasks []Task
    for cursor.Next(context.Background()) {
        var task Task
        err := cursor.Decode(&task)
        if err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode task"})
            return
        }
        tasks = append(tasks, task)
    }

    if err := cursor.Err(); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error"})
        return
    }

    if len(tasks) == 0 {
        ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "No tasks found"})
        return
    }

    ctx.IndentedJSON(http.StatusOK, gin.H{"Tasks": tasks})
}

func getTasksByID(ctx *gin.Context) {
    db := GetClient()
    id := ctx.Param("id")

    _, exists := ctx.Get("AuthenticatedUser")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    var task Task
    filter := bson.D{{Key: "_id", Value: id}}
    err := db.Collection("Task").FindOne(context.TODO(), filter).Decode(&task)
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

func updateTask(ctx *gin.Context) {
    db := GetClient()
    id := ctx.Param("id")
    updateTask := Task{}

    // Get authenticated user from Gin context
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

    // Retrieve the task from the database
    var existingTask Task
    filter := bson.D{{Key: "_id", Value: id}}
    err := db.Collection("Task").FindOne(context.Background(), filter).Decode(&existingTask)
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

    update := bson.D{}
    if updateTask.Title != "" {
        update = append(update, bson.E{Key: "$set", Value: bson.D{{Key: "title", Value: updateTask.Title}}})
    }
    if !updateTask.DueDate.IsZero() {
        update = append(update, bson.E{Key: "$set", Value: bson.D{{Key: "due_date", Value: updateTask.DueDate}}})
    }
    if updateTask.Status != "" {
        update = append(update, bson.E{Key: "$set", Value: bson.D{{Key: "status", Value: updateTask.Status}}})
    }
    if updateTask.Description != "" {
        update = append(update, bson.E{Key: "$set", Value: bson.D{{Key: "description", Value: updateTask.Description}}})
    }

    // Update the task in the database
    result, err := db.Collection("Task").UpdateOne(context.TODO(), filter, update)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
        return
    }

    if result.ModifiedCount == 0 {
        ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
    }

    ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Task Updated"})
}

func deleteTask(ctx *gin.Context) {
    db := GetClient()
    id := ctx.Param("id")

    // Get authenticated user from Gin context
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

    // Retrieve the task from the database
    var existingTask Task
    filter := bson.D{{Key: "_id", Value: id}}
    err := db.Collection("Task").FindOne(context.Background(), filter).Decode(&existingTask)
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
    result, err := db.Collection("Task").DeleteOne(context.TODO(), filter)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
        return
    }

    if result.DeletedCount == 0 {
        ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
    }

    ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Task removed"})
}


func main() {
	fmt.Println("Task Manger API")
	router := gin.Default()
	router.GET("/tasks", AuthMiddleware(), getTasks)
	router.GET("/tasks/:id", AuthMiddleware(), getTasksByID)
	router.PATCH("/tasks/:id/", AuthMiddleware(), updateTask)
	router.POST("/tasks", AuthMiddleware(), createTask)
	router.DELETE("/tasks/:id", AuthMiddleware(), deleteTask)


	router.POST("/register", createUser)
	router.POST("/register/admin", createAdmin)
	router.POST("/login", loginUser)
	// Listen and Server in :8080

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}

}
