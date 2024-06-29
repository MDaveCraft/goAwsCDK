package utils

import (
	"time"
	"lambda-func/types"
	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type HandlerFunc func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func ValidatePassword(hashedPassword, plainTextPasword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainTextPasword))
	return err == nil
}

func CreateToken(user types.User) string {
	now := time.Now()
	validUntil := now.Add(time.Hour * 1).Unix()

	claims := jwt.MapClaims{
		"user":    user.Username,
		"expires": validUntil,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims, nil)
	secret := "secret"

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return ""
	}

	return tokenString
}