package controller

import (
	"BlogApp/config"
	"BlogApp/domain/model"
	"BlogApp/domain/usecase"

	"github.com/gin-gonic/gin"
)

type BlogController struct {
	environment *config.Environment
	BlogUseCase usecase.BlogUseCase
}

func NewBlogController(environment *config.Environment, blogUseCase *usecase.BlogUseCase) *BlogController {
	return &BlogController{
		environment: environment,
		BlogUseCase: *blogUseCase,
	}

}

func (b *BlogController) CreateBlog(c *gin.Context) {
	PostHandler(c, b.BlogUseCase.CreateBlog, &model.BlogCreate{}, nil)
}

func (b *BlogController) GetByBlogID(c *gin.Context) {
	GetHandler(c, b.BlogUseCase.GetBlogByID, nil, &model.IdParam{ID: c.Param("blog_id")})
}

func (b *BlogController) GetAllBlogs(c *gin.Context) {
	GetHandler(c, b.BlogUseCase.GetBlogs, nil, &model.SearchParam{Search: c.Query("search")})
}

func (b *BlogController) UpdateBlog(c *gin.Context) {
	PutHandler(c, b.BlogUseCase.UpdateBlogByID, &model.BlogUpdate{}, &model.IdParam{ID: c.Param("blog_id")})
}

func (b BlogController) DeleteBlog(c *gin.Context) {
	DeleteHandler(c, b.BlogUseCase.DeleteBlogByID, nil, &model.IdParam{ID: c.Param("blog_id")})
}
package controller

import (
	"BlogApp/domain/model"
	"BlogApp/middleware"
	"BlogApp/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)





func PostHandler[T any, U any, V any](
	c *gin.Context, 
	handler func(currUser *model.AuthenticatedUser, dto T, param U) (V, string, error),
	dto T,
	param U) {
    err := c.BindJSON(dto)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }
	currUser, _ := utils.CheckUser(c)
    result, message, err := handler(currUser, dto, param)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    middleware.SuccessResponseHandler(c, http.StatusCreated, message, result)
}



// PUT
func PutHandler[T any, U any, V any](
    c *gin.Context, 
    handler func(currUser *model.AuthenticatedUser, dto T, param U) (V, string, error),
    dto T,
    param U) {
    err := c.BindJSON(dto)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }
    currUser, _ := utils.CheckUser(c)
    result, message, err := handler(currUser, dto, param)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    middleware.SuccessResponseHandler(c, http.StatusOK, message, result)
}

func GetHandler[T any, U any, V any](
    c *gin.Context, 
    handler func(currUser *model.AuthenticatedUser, dto T, param U) (V, string, error),
    dto T,
    param U) {
    currUser, _ := utils.CheckUser(c)
    result, message, err := handler(currUser, dto, param)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    middleware.SuccessResponseHandler(c, http.StatusOK, message, result)
}

func DeleteHandler[T any, U any, V any](
    c *gin.Context, 
    handler func(currUser *model.AuthenticatedUser, dto T, param U) (V, string, error),
    dto T,
    param U) {
    currUser, _ := utils.CheckUser(c)
    result, message, err := handler(currUser, dto, param)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    middleware.SuccessResponseHandler(c, http.StatusOK, message, result)
}


package controller

import (
	"BlogApp/config"
	"BlogApp/domain/model"
	"BlogApp/domain/usecase"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	environment *config.Environment
	AuthUseCase usecase.AuthUseCase
}


func NewAuthController(environment *config.Environment, authUseCase *usecase.AuthUseCase) *AuthController {
	return &AuthController{
		environment: environment,
		AuthUseCase: *authUseCase,
	}
}

func (a *AuthController) Register(c *gin.Context) {
	PostHandler(c, a.AuthUseCase.Register, &model.UserCreate{}, nil)
}

func (a *AuthController) Login(c *gin.Context) {
	PostHandler(c, a.AuthUseCase.Login, &model.UserLogin{}, nil)
}

func (a *AuthController) AdminRegister(c *gin.Context) {
	PostHandler(c, a.AuthUseCase.AdminRegister, &model.UserCreate{}, nil)
}


package controller

import (
	"BlogApp/config"
	"BlogApp/domain/model"
	"BlogApp/domain/usecase"

	"github.com/gin-gonic/gin"
)

type ProfileController struct{
	environment *config.Environment
	ProfileUseCase usecase.ProfileUseCase
}


func NewProfileController(environment *config.Environment, profileUseCase *usecase.ProfileUseCase) *ProfileController {
	return &ProfileController{
		environment: environment,
		ProfileUseCase: *profileUseCase,
	}

}

func (p *ProfileController) GetProfile(c *gin.Context) {
	GetHandler(c, p.ProfileUseCase.GetProfile, nil, nil)
}

func (p *ProfileController) UpdateProfile(c *gin.Context) {
	PutHandler(c, p.ProfileUseCase.UpdateProfile, &model.UserUpdateProfile{}, nil)
}

func (p *ProfileController) DeleteProfile(c *gin.Context) {
	DeleteHandler(c, p.ProfileUseCase.DeleteProfile, nil, nil)
}

func (p *ProfileController) UpdatePassword(c *gin.Context) {
	PutHandler(c, p.ProfileUseCase.UpdatePassword, &model.UserUpdatePassword{}, nil)
}

func (p *ProfileController) UpdateEmail(c *gin.Context) {
	PutHandler(c, p.ProfileUseCase.UpdateEmail, &model.UserUpdateEmail{}, nil)
}

func (p *ProfileController) UpdateUsername(c *gin.Context) {
	PutHandler(c, p.ProfileUseCase.UpdateUsername, &model.UserUpdateUsername{}, nil)
}


// 