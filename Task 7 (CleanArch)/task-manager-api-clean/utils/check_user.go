package utils

import (
	"task-manager-api-clean/domain"
	"errors"
	"github.com/gin-gonic/gin"
)

func CheckUser(c *gin.Context) (*domain.AuthenticatedUser, error) {
	value, exist := c.Get("AuthenticatedUser")
	if !exist {
		return nil, errors.New("user not found")
	}
	currUser, ok := value.(*domain.AuthenticatedUser)
	if !ok {
		return nil, errors.New("user not found in context")
	} 
	return currUser, nil
}