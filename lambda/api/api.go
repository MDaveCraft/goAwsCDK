package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/database/interfaces"
	"lambda-func/types"
	"net/http"
	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	dynamoClient interfaces.IDbCLient
}

func NewApiHandler(dynamoClient interfaces.IDbCLient) *ApiHandler {
	return &ApiHandler{
		dynamoClient: dynamoClient,
	}
}

func (api* ApiHandler) InsertUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	
	var registerUser types.RegisterUser

	err:= json.Unmarshal([]byte(request.Body), &registerUser)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Invalid request body",
			StatusCode:  http.StatusBadRequest,
		}, err
	}

	if registerUser.Username == "" || registerUser.Password == "" {
		return events.APIGatewayProxyResponse{
			Body: "Invalid request - fields empty",
			StatusCode:  http.StatusBadRequest,
		}, err
	}

	userExists, err := api.dynamoClient.UserExists(registerUser.Username)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Internal Server Error",
			StatusCode:  http.StatusInternalServerError,
		}, err
	}

	if userExists {
		return events.APIGatewayProxyResponse{
			Body: "User already exists",
			StatusCode:  http.StatusConflict,
		}, err
	}

	user, err := types.NewUser(registerUser)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Internal Server Error",
			StatusCode:  http.StatusInternalServerError,
		}, fmt.Errorf("error creating user: %w", err)
	}

	err = api.dynamoClient.InsertUser(user)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Internal Server Error",
			StatusCode:  http.StatusInternalServerError,
		}, fmt.Errorf("error inserting user: %w", err)
	}

	return events.APIGatewayProxyResponse{
		Body: "User inserted successfully",
		StatusCode:  http.StatusCreated,
	}, nil
}

func (api* ApiHandler) LoginUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var login types.RegisterUser

	err:= json.Unmarshal([]byte(request.Body), &login)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Invalid request body",
			StatusCode:  http.StatusBadRequest,
		}, err
	}

	if login.Username == "" || login.Password == "" {
		return events.APIGatewayProxyResponse{
			Body: "Invalid request - fields empty",
			StatusCode:  http.StatusBadRequest,
		}, err
	}

	userExists, err := api.dynamoClient.UserExists(login.Username)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Internal Server Error",
			StatusCode:  http.StatusInternalServerError,
		}, err
	}

	if !userExists {
		return events.APIGatewayProxyResponse{
			Body: "User does not exist",
			StatusCode:  http.StatusNotFound,
		}, err
	}

	user, err := api.dynamoClient.GetUser(login.Username)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Internal Server Error",
			StatusCode:  http.StatusInternalServerError,
		}, fmt.Errorf("error getting user: %w", err)
	}

	if !types.ValidatePassword(user.PasswordHash, login.Password) {
		return events.APIGatewayProxyResponse{
			Body: "Invalid password",
			StatusCode:  http.StatusUnauthorized,
		}, err
	}

	accessToken,err := types.CreateToken(user.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Internal Server Error",
			StatusCode:  http.StatusInternalServerError,
		}, fmt.Errorf("error creating token: %w", err)
	}
	
	successMsg := fmt.Sprintf(`{"accessToken": "%s"}`, accessToken)

	return events.APIGatewayProxyResponse{
		Body: successMsg,
		StatusCode:  http.StatusOK,
	}, nil
}

func ProtectedHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body: "This is a protected route",
		StatusCode: http.StatusOK,
	}, nil
}