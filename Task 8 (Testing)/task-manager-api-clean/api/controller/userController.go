package controller

import (
	"net/http"

	"task-manager-api-clean/domain"
	"task-manager-api-clean/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUseCase domain.UserUseCase
}

func NewUserController(userUseCase domain.UserUseCase) *UserController {
	return &UserController{
		userUseCase: userUseCase,
	}
}


func (uc *UserController) CreateUser(ctx *gin.Context) {
	var newUser domain.UserCreate

	// Bind JSON to new user
	if err := ctx.BindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	_, err := uc.userUseCase.RegisterUser(ctx, &newUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (uc *UserController) LoginUser(ctx *gin.Context) {
	var user domain.UserLogin

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwtToken, err := uc.userUseCase.Login(ctx, &user)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"message": "User logged in successfully", "token": jwtToken})
}

func (uc *UserController) PromoteUser(ctx *gin.Context) {
	// Get authenticated user from gin context
	user, err := utils.CheckUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if user.Role != "admin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Only the admin can promote a user"})
		return
	}

	username := ctx.Param("username")

	_, er := uc.userUseCase.Promote(ctx, username)
	if er != nil {
		if er.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": er.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": er.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User promoted to admin"})
}

