package types

import (
	"fmt"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("secret-key")

type RegisterUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}

func NewUser(regUser RegisterUser) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(regUser.Password), 31)
	if err != nil {
		return User{}, err
	}
	return User{
		Username:     regUser.Username,
		PasswordHash: string(hashedPassword),
	}, nil
}

func ValidatePassword(hashPassword string, plainTextPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(plainTextPassword))
	return err == nil
}

func CreateToken(username string) (string, error) {
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
			jwt.MapClaims{ 
				"username": username, 
				"exp": time.Now().Add(time.Hour * 24).Unix(), 
			})
  tokenString, err := token.SignedString(secretKey)
  if err != nil {
		return "", err
  }
	return tokenString, nil
}

func VerifyToken(tokenString string) error {
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    return secretKey, nil
  })
  if err != nil {
    return err
  }
  if !token.Valid {
    return fmt.Errorf("invalid token")
  }
  return nil
}