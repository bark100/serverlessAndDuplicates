package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"os"
)

type Event struct {
	Id string `json:"id"`
}

func HandleRequest(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	print("[INFO] Init...\n")

	// Parsing request...
	event := Event{}
	err := json.Unmarshal([]byte(request.Body), &event)
	if err != nil {
		fmt.Println("ERROR:" + err.Error())
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	// SDK session and credentials:
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	tableName := os.Getenv("DYNAMODB_TABLE")

	// Getting item from DB:
	print(fmt.Sprintf("[INFO] Getting ID '%s' from database...\n", event.Id))
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(event.Id),
			},
		},
	})
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}
	if result.Item == nil {
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("{\"message\": \"not found.\"}"), StatusCode: 500}, err
	}

	fmt.Println("[INFO] Item result:  ", *result.Item["Result"].N)
	return events.APIGatewayProxyResponse{Body: fmt.Sprintf("{\"result\": %s}", *result.Item["Result"].N), StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
