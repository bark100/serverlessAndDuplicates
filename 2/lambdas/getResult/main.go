package main

import (
	"context"
	"fmt"
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

//type Item struct {
//	Id     string
//	Result string
//}

func HandleRequest(ctx context.Context, event Event) (string, error) {
	print("[INFO] Init...")

	// SDK session and credentials:
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	tableName := os.Getenv("DYNAMODB_TABLE")

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
		return "Could not find '" + event.Id + "'", nil
	}

	//err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	//if err != nil {
	//	panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	//}

	fmt.Println("[INFO] Item:  ", result.Item)

	return fmt.Sprintf("%s", result.Item), nil
}

func main() {
	lambda.Start(HandleRequest)
}
