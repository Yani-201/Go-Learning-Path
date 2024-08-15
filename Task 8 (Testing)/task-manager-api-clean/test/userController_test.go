package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"task-manager-api-clean/api/controller"
	"task-manager-api-clean/api/middleware"
	"task-manager-api-clean/domain"
	"task-manager-api-clean/domain/mocks"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserControllerTestSuite struct {
    suite.Suite
    useCase    *mocks.UserUseCase
    controller *controller.UserController
    router     *gin.Engine
    secret    string
}


func (suite *UserControllerTestSuite) SetupTest() {
    suite.useCase = new(mocks.UserUseCase)
    suite.controller = controller.NewUserController(suite.useCase)
    gin.SetMode(gin.TestMode)
    suite.router = gin.Default()
    auth := middleware.AuthMiddleware(suite.secret)


    suite.router.POST("/users", suite.controller.CreateUser)
    suite.router.POST("/login", suite.controller.LoginUser)
    suite.router.POST("/promote/:username", auth, suite.controller.PromoteUser)
}

func (suite *UserControllerTestSuite) TearDownTest() {
    suite.useCase.AssertExpectations(suite.T())
}


// Helper function to create a JWT for testing
func (suite *UserControllerTestSuite) createTestJWT(userID, username, role string) string {
	claims := &jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(suite.secret))
	suite.NoError(err)
	return tokenString
}


func (suite *UserControllerTestSuite) TestCreateUser_Success() {
    payload := &domain.UserCreate{Username: "test", Password: "test", Email: "test@example.com"}
    suite.useCase.On("RegisterUser", mock.Anything, payload).Return(&domain.UserInfo{UserId: "1", Username: "test", Email: "test@example.com"}, nil)

    body, _ := json.Marshal(payload)
    req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
    w := httptest.NewRecorder()
    suite.router.ServeHTTP(w, req)

    suite.Equal(http.StatusCreated, w.Code)
    suite.useCase.AssertCalled(suite.T(), "RegisterUser", mock.Anything, payload)
}

func (suite *UserControllerTestSuite) TestCreateUser_BadRequest() {
    body := []byte(`{"username":"test"}`)
    req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
    w := httptest.NewRecorder()
    suite.router.ServeHTTP(w, req)

	suite.useCase.AssertNotCalled(suite.T(), "RegisterUser", mock.Anything, mock.Anything)

    suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *UserControllerTestSuite) TestCreateUser_UseCaseError() {
    payload := &domain.UserCreate{Username: "test", Password: "test", Email: "test@example.com"}
    suite.useCase.On("RegisterUser", mock.Anything, payload).Return(nil, errors.New("user already exists"))

    body, _ := json.Marshal(payload)
    req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
    w := httptest.NewRecorder()
    suite.router.ServeHTTP(w, req)

    suite.Equal(http.StatusBadRequest, w.Code)
    suite.useCase.AssertCalled(suite.T(), "RegisterUser", mock.Anything, payload)
	

}

func (suite *UserControllerTestSuite) TestLoginUser_Success() {
    payload := &domain.UserLogin{Username: "test", Password: "test"}
    suite.useCase.On("Login", mock.Anything, payload).Return("jwtToken", nil)

    body, _ := json.Marshal(payload)
    req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
    w := httptest.NewRecorder()
    suite.router.ServeHTTP(w, req)

    suite.Equal(http.StatusOK, w.Code)
    suite.useCase.AssertCalled(suite.T(), "Login", mock.Anything, payload)
}

func (suite *UserControllerTestSuite) TestLoginUser_BadRequest() {
    body := []byte(`{"username":"test"}`)
    req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
    w := httptest.NewRecorder()
    suite.router.ServeHTTP(w, req)

    suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *UserControllerTestSuite) TestLoginUser_UseCaseError() {
    payload := &domain.UserLogin{Username: "test", Password: "test"}
    suite.useCase.On("Login", mock.Anything, payload).Return("", errors.New("invalid credentials"))

    body, _ := json.Marshal(payload)
    req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
    w := httptest.NewRecorder()
    suite.router.ServeHTTP(w, req)

    suite.Equal(http.StatusUnauthorized, w.Code)
    suite.useCase.AssertCalled(suite.T(), "Login", mock.Anything, payload)

}

func (suite *UserControllerTestSuite) TestPromoteUser_Success() {
    username := "testUser"
	tokenString := suite.createTestJWT("1", "adminUser", "admin")


    // Set up mock for Promote function
    expectedUserInfo := &domain.UserInfo{
        UserId:   "1",
        Username: username,
        Email:    "test@example.com",
    }
    suite.useCase.On("Promote", mock.Anything, username).Return(expectedUserInfo, nil)

    req, _ := http.NewRequest("POST", "/promote/"+username, nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()

    suite.router.ServeHTTP(w, req)
	suite.Equal(http.StatusOK, w.Code)

    expectedResponse := `{"message":"User promoted to admin"}`
    suite.JSONEq(expectedResponse, w.Body.String())
    suite.useCase.AssertCalled(suite.T(), "Promote", mock.Anything, username)
}

func (suite *UserControllerTestSuite) TestPromoteUser_Unauthorized() {
    username := "testUser"
	tokenString := suite.createTestJWT("1", "adminUser", "user")


    req, _ := http.NewRequest("POST", "/promote/"+username, nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()

    suite.router.ServeHTTP(w, req)
    suite.Equal(http.StatusUnauthorized, w.Code)

    expectedResponse := `{"error":"Unauthorized: Only the admin can promote a user"}`
    suite.JSONEq(expectedResponse, w.Body.String())
    suite.useCase.AssertNotCalled(suite.T(), "Promote", mock.Anything, username)
}

func (suite *UserControllerTestSuite) TestPromoteUser_UserNotFound() {
    username := "nonExistentUser"
    tokenString := suite.createTestJWT("1", "adminUser", "admin")

    suite.useCase.On("Promote", mock.Anything, username).Return(nil, errors.New("user not found"))

    req, _ := http.NewRequest("POST", "/promote/"+username, nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()

    suite.router.ServeHTTP(w, req)
    suite.Equal(http.StatusNotFound, w.Code)

    expectedResponse := `{"error":"user not found"}`
    suite.JSONEq(expectedResponse, w.Body.String())
    suite.useCase.AssertCalled(suite.T(), "Promote", mock.Anything, username)
}


func TestUserControllerTestSuite(t *testing.T) {
    suite.Run(t, new(UserControllerTestSuite))
}









