package app

import (
	"lambda-func/api"
	"lambda-func/middleware"
	db "lambda-func/database"
	"net/http"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	tableName = "user"
)

type App struct {
	ApiHandler *api.ApiHandler
}

func NewApp() *App{
	dynamoDB := db.NewDynamoDBClient(tableName)
	return &App{
		ApiHandler: api.NewApiHandler(*dynamoDB),
	}
}

func (app* App) Run() {
	lambda.Start(
		func (request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			switch request.Path {
				case "/register":
					return app.ApiHandler.InsertUserHandler(request)
				case "/login":
					return app.ApiHandler.LoginUser(request)
				case "/protected":
					return middleware.ValidateJWTMiddleware(api.ProtectedHandler)(request)
				default:
					return events.APIGatewayProxyResponse{
						Body: "Not Found",
						StatusCode: http.StatusNotFound,
					}, nil
			}
	})
}
