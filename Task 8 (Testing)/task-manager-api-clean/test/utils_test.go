package tests

import (
    "task-manager-api-clean/domain"
    "task-manager-api-clean/utils"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/suite"
)

type UtilsTestSuite struct {
    suite.Suite
}

func (suite *UtilsTestSuite) TestCheckUser_UserFound() {
    c, _ := gin.CreateTestContext(nil)
    authUser := &domain.AuthenticatedUser{
        UserID:   "123",
        Username: "testuser",
        Role:     "admin",
	}
    c.Set("AuthenticatedUser", authUser)

    user, err := utils.CheckUser(c)
    require.NoError(suite.T(), err, "CheckUser returned an error")
    assert.Equal(suite.T(), authUser, user, "The returned user does not match the expected authenticated user")
}

func (suite *UtilsTestSuite) TestCheckUser_UserNotFound() {
    c, _ := gin.CreateTestContext(nil)

    user, err := utils.CheckUser(c)
    assert.Error(suite.T(), err, "Expected an error when user is not found")
    assert.Nil(suite.T(), user, "Expected user to be nil when not found")
}

func (suite *UtilsTestSuite) TestCheckUser_UserNotInContext() {
    c, _ := gin.CreateTestContext(nil)
    c.Set("AuthenticatedUser", "invalidUserType") 

    user, err := utils.CheckUser(c)
    assert.Error(suite.T(), err, "Expected an error when user is not found in context")
    assert.Nil(suite.T(), user, "Expected user to be nil when not found in context")
}

func (suite *UtilsTestSuite) TestEncryptPassword_Success() {
    password := "mySecretPassword"
    hashedPassword, err := utils.EncryptPassword(password)
    require.NoError(suite.T(), err, "EncryptPassword should not return an error")
    assert.NotEqual(suite.T(), password, hashedPassword, "Hashed password should not be the same as the plain text password")
    assert.NotEmpty(suite.T(), hashedPassword, "Hashed password should not be empty")
}

func (suite *UtilsTestSuite) TestComparePasswords_Success() {
    password := "mySecretPassword"
    hashedPassword, _ := utils.EncryptPassword(password)

    err := utils.ComparePasswords(hashedPassword, password)
    assert.NoError(suite.T(), err, "Passwords should match")
}

func (suite *UtilsTestSuite) TestComparePasswords_Failure() {
    password := "mySecretPassword"
    hashedPassword, _ := utils.EncryptPassword(password)

    err := utils.ComparePasswords(hashedPassword, "wrongPassword")
    assert.Error(suite.T(), err, "Passwords should not match")
}

func (suite *UtilsTestSuite) TestTokenGenerate() {
    authUser := &domain.AuthenticatedUser{
        UserID:   "123",
        Username: "testuser",
        Role:     "admin",
    }

    token, err := utils.TokenGenerate(authUser, "secret")
    require.NoError(suite.T(), err, "Token generation should succeed")
    assert.NotEmpty(suite.T(), token, "Generated token should not be empty")
}

func TestUtilsTestSuite(t *testing.T) {
    suite.Run(t, new(UtilsTestSuite))
}