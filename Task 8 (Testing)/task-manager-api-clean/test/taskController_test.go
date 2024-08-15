package tests

import (
	"bytes"
	
	"encoding/json"
	
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"task-manager-api-clean/api/controller"
	"task-manager-api-clean/api/middleware"
	"task-manager-api-clean/domain"
	"task-manager-api-clean/domain/mocks"
	// "task-manager-api-clean/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskControllerTestSuite struct {
	suite.Suite
	router         *gin.Engine
	useCase        *mocks.TaskUseCase
	taskController *controller.TaskController
	secret         string
}

func (suite *TaskControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()
	suite.useCase = new(mocks.TaskUseCase)
	suite.taskController = controller.NewTaskController(suite.useCase)
	suite.secret = "secret" 

	// Apply the middleware
	auth := middleware.AuthMiddleware(suite.secret)
	suite.router.Use(auth)

	suite.router.POST("/tasks", suite.taskController.CreateTask)
	suite.router.GET("/tasks", suite.taskController.GetTasks)
	suite.router.GET("/tasks/:id", suite.taskController.GetTaskByID)
	suite.router.PUT("/tasks/:id", suite.taskController.UpdateTask)
	suite.router.DELETE("/tasks/:id", suite.taskController.DeleteTask)
}

func (suite *TaskControllerTestSuite) TearDownTest() {
	suite.useCase.AssertExpectations(suite.T())
}

// Helper function to create a JWT for testing
func (suite *TaskControllerTestSuite) createTestJWT(userID, username, role string) string {
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

func (suite *TaskControllerTestSuite) TestCreateTask_Success() {
	taskInput := &domain.TaskInput{Title: "Test Task", Description: "Test Description", Status: "Pending"}
	task := &domain.Task{Id: "1", Title: taskInput.Title, Description: taskInput.Description, Status: taskInput.Status}

	suite.useCase.On("Create", mock.Anything, taskInput).Return(task, nil)

	body, _ := json.Marshal(taskInput)
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	token := suite.createTestJWT("123", "testuser", "admin")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	suite.Equal(http.StatusCreated, resp.Code)
	suite.useCase.AssertCalled(suite.T(), "Create", mock.Anything, taskInput)
}

func (suite *TaskControllerTestSuite) TestCreateTask_Failure_Unauthorized() {
	taskInput := &domain.TaskInput{Title: "Test Task", Description: "Test Description", Status: "Pending", DueDate: time.Now()}

	body, _ := json.Marshal(taskInput)
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
	token := suite.createTestJWT("123", "testuser", "user")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	suite.Equal(http.StatusUnauthorized, resp.Code)
	suite.Contains(resp.Body.String(), "Unauthorized: Only the admin can create a task")
	suite.useCase.AssertNotCalled(suite.T(), "Create", mock.Anything, mock.Anything)
}

func (suite *TaskControllerTestSuite) TestGetTasks_Success() {
	tasks := []*domain.Task{
		{Id: "1", Title: "Task 1", Description: "Description 1", Status: "Pending", DueDate: time.Now()},
		{Id: "2", Title: "Task 2", Description: "Description 2", Status: "Completed", DueDate: time.Now()},
	}

	suite.useCase.On("GetAll", mock.Anything).Return(&tasks, nil)

	req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
	token := suite.createTestJWT("123", "testuser", "admin")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	suite.Equal(http.StatusOK, resp.Code)
	suite.Contains(resp.Body.String(), "Task 1")
	suite.Contains(resp.Body.String(), "Task 2")
	suite.useCase.AssertCalled(suite.T(), "GetAll", mock.Anything)
}

func (suite *TaskControllerTestSuite) TestGetTasks_Failure_NotFound() {
	var tasks []*domain.Task
	suite.useCase.On("GetAll", mock.Anything).Return(&tasks, nil)

	req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
	token := suite.createTestJWT("123", "testuser", "admin")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	suite.Equal(http.StatusNotFound, resp.Code)
	suite.Contains(resp.Body.String(), "No tasks found")
	suite.useCase.AssertCalled(suite.T(), "GetAll", mock.Anything)
}

func (suite *TaskControllerTestSuite) TestGetTaskByID_Success() {
	task := &domain.Task{Id: "1", Title: "Task 1", Description: "Description 1", Status: "Pending"}

	suite.useCase.On("GetById", mock.Anything, "1").Return(task, nil)

	req, _ := http.NewRequest(http.MethodGet, "/tasks/1", nil)
	token := suite.createTestJWT("123", "testuser", "admin")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	suite.Equal(http.StatusOK, resp.Code)
	suite.Contains(resp.Body.String(), "Task 1")
	suite.useCase.AssertCalled(suite.T(), "GetById", mock.Anything, "1")
}


func (suite *TaskControllerTestSuite) TestGetTaskByID_Failure_NotFound() {
	suite.useCase.On("GetById", mock.Anything, "nonExistentId").Return(nil, mongo.ErrNoDocuments)

	req, _ := http.NewRequest(http.MethodGet, "/tasks/nonExistentId", nil)
	token := suite.createTestJWT("123", "testuser", "admin")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	suite.Equal(http.StatusNotFound, resp.Code)
	suite.Contains(resp.Body.String(), "Task not found")
	suite.useCase.AssertCalled(suite.T(), "GetById", mock.Anything, "nonExistentId")
}

func (suite *TaskControllerTestSuite) TestUpdateTask_Success() {
	taskInput := &domain.TaskInput{Title: "Updated Task", Description: "Updated Description", Status: "Completed"}
	task := &domain.Task{Id: "1", Title: taskInput.Title, Description: taskInput.Description, Status: taskInput.Status}

	suite.useCase.On("GetById", mock.Anything, "1").Return(task, nil)
	suite.useCase.On("Update", mock.Anything, "1", taskInput).Return(task, nil)

	body, _ := json.Marshal(taskInput)
	req, _ := http.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBuffer(body))
		
	token := suite.createTestJWT("123", "testuser", "admin")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	suite.Equal(http.StatusOK, resp.Code)
	suite.Contains(resp.Body.String(), "Updated Task")
	suite.useCase.AssertCalled(suite.T(), "GetById", mock.Anything, "1")
	suite.useCase.AssertCalled(suite.T(), "Update", mock.Anything, "1", taskInput)

}

func (suite *TaskControllerTestSuite) TestUpdateTask_Failure_Unauthorized() {
	taskInput := &domain.TaskInput{Title: "Updated Task", Description: "Updated Description", Status: "Completed", DueDate: time.Now()}

	body, _ := json.Marshal(taskInput)
	req, _ := http.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBuffer(body))
	token := suite.createTestJWT("123", "testuser", "user")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	suite.Equal(http.StatusUnauthorized, resp.Code)
	suite.Contains(resp.Body.String(), "Unauthorized: Only the admin can update this task")
	suite.useCase.AssertNotCalled(suite.T(), "Update", mock.Anything, mock.Anything, mock.Anything)
}

func (suite *TaskControllerTestSuite) TestDeleteTask_Success() {
	task := &domain.Task{Id: "1", Title: "Task 1", Description: "Description 1", Status: "Pending"}

	suite.useCase.On("GetById", mock.Anything, "1").Return(task, nil)
	suite.useCase.On("Delete", mock.Anything, "1").Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/tasks/1", nil)
	token := suite.createTestJWT("123", "testuser", "admin")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	suite.Equal(http.StatusOK, resp.Code)
	suite.Contains(resp.Body.String(), "Task removed")
	suite.useCase.AssertCalled(suite.T(), "GetById", mock.Anything, "1")
	suite.useCase.AssertCalled(suite.T(), "Delete", mock.Anything, "1")

}

func (suite *TaskControllerTestSuite) TestDeleteTask_Failure_Unauthorized() {
	req, _ := http.NewRequest(http.MethodDelete, "/tasks/1", nil)

	// Set the JWT token for a non-admin user
	token := suite.createTestJWT("123", "testuser", "user")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	suite.router.ServeHTTP(resp, req)

	suite.Equal(http.StatusUnauthorized, resp.Code)
	suite.Contains(resp.Body.String(), "Unauthorized: Only the admin can delete this task")
	suite.useCase.AssertNotCalled(suite.T(), "Delete", mock.Anything, mock.Anything)
	
}
func TestTaskControllerTestSuite(t *testing.T) {
	suite.Run(t, new(TaskControllerTestSuite))
}


































