package tests

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "task-manager-api-clean/api/middleware"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

type MiddlewareTestSuite struct {
    suite.Suite
    router *gin.Engine
}

func (suite *MiddlewareTestSuite) SetupTest() {
    suite.router = gin.New()
    suite.router.Use(middleware.AuthMiddleware("secret"))
    suite.router.GET("/test", func(c *gin.Context) {
        c.Status(http.StatusOK)
    })
}

func (suite *MiddlewareTestSuite) TearDownTest() {
    suite.router = nil
}

func (suite *MiddlewareTestSuite) TestAuthMiddleware_MissingAuthorizationHeader() {
    req, err := http.NewRequest("GET", "/test", nil)
    suite.NoError(err)

    w := httptest.NewRecorder()
    suite.router.ServeHTTP(w, req)

    assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
    assert.Contains(suite.T(), w.Body.String(), "Authorization header is required")
}

func (suite *MiddlewareTestSuite) TestAuthMiddleware_InvalidAuthorizationHeaderFormat() {
    req, err := http.NewRequest("GET", "/test", nil)
    suite.NoError(err)

    req.Header.Set("Authorization", "InvalidFormat")
    w := httptest.NewRecorder()
    suite.router.ServeHTTP(w, req)

    assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
    assert.Contains(suite.T(), w.Body.String(), "Invalid authorization header")
}

func (suite *MiddlewareTestSuite) TestAuthMiddleware_InvalidJWTToken() {
    req, err := http.NewRequest("GET", "/test", nil)
    suite.NoError(err)

    req.Header.Set("Authorization", "Bearer invalid.token")
    w := httptest.NewRecorder()
    suite.router.ServeHTTP(w, req)

    assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
    assert.Contains(suite.T(), w.Body.String(), "Invalid JWT")
}

func (suite *MiddlewareTestSuite) TestAuthMiddleware_ValidToken() {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id":  "123",
        "username": "testuser",
        "role":     "admin",
    })
    tokenString, err := token.SignedString([]byte("secret"))
    suite.NoError(err)

    req, err := http.NewRequest("GET", "/test", nil)
    suite.NoError(err)

    req.Header.Set("Authorization", "Bearer "+tokenString)
    w := httptest.NewRecorder()
    suite.router.ServeHTTP(w, req)

    assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func TestMiddlewareTestSuite(t *testing.T) {
    suite.Run(t, new(MiddlewareTestSuite))
}