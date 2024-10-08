package utils

import (
	"time"
    "task-manager-api-clean/domain"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (string, error){
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPass), err
}


func ComparePasswords(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func TokenGenerate(auth *domain.AuthenticatedUser, secret string) (string, error) {
    token := jwt.New(jwt.SigningMethodHS256)

    // Set claims
    claims := token.Claims.(jwt.MapClaims)
    claims["username"] = auth.Username
    claims["role"] = auth.Role
    claims["user_id"] = auth.UserID
    claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

    tokenString, err := token.SignedString([]byte(secret))
    if err != nil {
        return "", err
    }

    return tokenString, nil
	
}

