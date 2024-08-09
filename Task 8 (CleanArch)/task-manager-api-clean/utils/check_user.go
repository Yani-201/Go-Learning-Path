package utils

import (
	"task-manager-api/model"
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckUser(c *gin.Context) (*model.AuthenticatedUser, error) {
	value, exist := c.Get("AuthenticatedUser")
	if !exist {
		return nil, errors.New("user not found")
	}
	currUser, ok := value.(*model.AuthenticatedUser)
	if !ok {
		return nil, errors.New("user not found in context")
	} 
	return currUser, nil
}