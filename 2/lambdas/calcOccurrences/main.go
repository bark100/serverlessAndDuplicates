package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"log"
	"os"
	"strings"
	"time"
)

type Event struct {
	Word string `json:"word"`
	Char string `json:"character"`
}

type Item struct {
	Id        string
	Timestamp int
	Result    int
}

func HandleRequest(ctx context.Context, event Event) (int, error) {
	result := strings.Count(event.Word, event.Char)
	print(fmt.Sprintf("[INFO] Number of occurences of '%s' in '%s' is: '%d'\n", event.Char, event.Word, result))

	print("[INFO] Init...")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{},
	}))
	svc := dynamodb.New(sess)

	item := Item{
		Id:        uuid.New().String(),
		Timestamp: int(time.Now().UnixMilli()),
		Result:    result,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("Got error marshalling item: %s", err)
	}

	tableName := os.Getenv("DYNAMODB_TABLE")
	input := &dynamodb.PutItemInput{
		TableName:              aws.String(tableName),
		Item:                   av,
		ReturnConsumedCapacity: aws.String("TOTAL"),
	}

	fmt.Println("[DEBUG] Saving new item: " + fmt.Sprint(input))
	out, err := svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}
	log.Printf("[INFO] Successfully added item id '%s' with result: '%d' to table '%s' with consumed capacity units: '%s'", item.Id, item.Result, tableName, out.ConsumedCapacity.CapacityUnits)

	return result, nil
}

func main() {
	lambda.Start(HandleRequest)
}
