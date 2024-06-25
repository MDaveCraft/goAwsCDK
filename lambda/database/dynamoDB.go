package database

import (
	"fmt"
	"lambda-func/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDBClient struct {
	dynamoClient *dynamodb.DynamoDB
	tableName    string
}

func NewDynamoDBClient(tableName string) *DynamoDBClient {
	dbSession := session.Must(session.NewSession())
	db := dynamodb.New(dbSession)
	return &DynamoDBClient{
		dynamoClient: db,
		tableName:    tableName,
	}
}

func (client DynamoDBClient) InsertUser(user types.User) error {
	if user.Username == "" || user.PasswordHash == "" {
		return fmt.Errorf("username and password can't be empty")
	}
	_, err := client.dynamoClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(client.tableName),
		Item: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(user.Username),
			},
			"password": {
				S: aws.String(user.PasswordHash),
			},
		},
	})
	return err
}

func (client DynamoDBClient) UserExists(username string) (bool, error) {
	if username == "" {
		return false, fmt.Errorf("username can't be empty")
	}
	result, err := client.dynamoClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(client.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})
	if err != nil {
		return true, err
	}
	if result.Item == nil {
		return false, nil
	}
	return true, nil
}

func (client DynamoDBClient) GetUser(username string) (types.User, error) {
	var user types.User
	if username == "" {
		return user, fmt.Errorf("username can't be empty")
	}
	result, err := client.dynamoClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(client.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})
	if err != nil {
		return user, fmt.Errorf("error getting user: %w", err)
	}
	if result.Item == nil {
		return user, fmt.Errorf("user not found")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return user, err
	}
	return user, nil
}
