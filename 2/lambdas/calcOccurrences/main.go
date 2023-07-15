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
	Word      string
	Char      string
}

func HandleRequest(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	event := Event{}
	err := json.Unmarshal([]byte(request.Body), &event)
	if err != nil {
		fmt.Println("ERROR:" + err.Error())
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	result := strings.Count(event.Word, event.Char)
	print(fmt.Sprintf("[INFO] Number of occurences of '%s' in '%s' is: '%d'\n", event.Char, event.Word, result))

	print("[INFO] Init AWS session...")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{},
	}))
	svc := dynamodb.New(sess)

	item := Item{
		Id:        uuid.New().String(),
		Timestamp: int(time.Now().UnixMilli()),
		Result:    result,
		Word:      event.Word,
		Char:      event.Char,
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

	out, err := svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}
	log.Printf("[INFO] Successfully added item id '%s' with result: '%d' to table '%s' with consumed capacity units: '%d'", item.Id, item.Result, tableName, out.ConsumedCapacity.CapacityUnits)

	return events.APIGatewayProxyResponse{Body: fmt.Sprintf("{\"result\": %d, \"id\": \"%s\"}", item.Result, item.Id), StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
