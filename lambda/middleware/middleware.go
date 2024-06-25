package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
)

type handlerFunc func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func ValidateJWTMiddleware(next handlerFunc) handlerFunc {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		tokenString := extractTokenFromHeaders(request.Headers)
		if tokenString == "" {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Body: "Missing Auth Token",
			}, nil
		}
		claims, err := parseToken(tokenString)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Body: "Unauthorized",
			}, nil
		}
		exp := int64(claims["expire"].(float64))
		if exp < time.Now().Unix() {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusUnauthorized,
				Body: "Token Expired",
			}, nil
		}
		return next(request)
	}
}

func extractTokenFromHeaders(header map[string]string) string{
	authHeader, ok := header["Authorization"]
	if !ok {
		return ""
	}

	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return ""
	}
	return splitToken[1]
}

func parseToken(tokenString string) (jwt.MapClaims,error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		return nil, fmt.Errorf("unauthorized - %v", http.StatusUnauthorized)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid - unauthorized [%v]", http.StatusUnauthorized)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unauthorized - %v", http.StatusUnauthorized)
	}

	return claims, nil
}