package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"task-manager-api-clean/domain"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secret string) gin.HandlerFunc {

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

				return []byte(secret), nil
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
			c.Set("AuthenticatedUser", &domain.AuthenticatedUser{
				UserID: UserID,
				Username: Username,
				Role: Role,
			})
	
			c.Next()
		}

	
}
