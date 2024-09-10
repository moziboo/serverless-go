package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// PK : USER#KELLY
// SK : PROFILE
type UserProfile struct {
	PK    string `dynamodbav:"PK" json:"pk"`
	SK    string `dynamodbav:"SK" json:"sk"`
	Email string `dynamodbav:"Email" json:"email"`
	PfP   string `dynamodbav:"PfP" json:"pfp"`
	Phone string `dynamodbav:"Phone" json:"phone"`
	Name  string `dynamodbav:"Name" json:"name"`
}

var svc *dynamodb.Client

func acceptableOrigin(origin string) bool {
	acceptableOrigins := []string{"http://localhost:3000", "http://localhost:3001", "https://mmmelton.com"}

	for _, v := range acceptableOrigins {
		if v == origin {
			return true
		}
	}
	return false
}

func HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := make(map[string]string)

	headers["Content-Type"] = "application/json"
	headers["Access-Control-Allow-Credentials"] = "true"
	headers["Access-Control-Allow-Headers"] = "Content-Type, Authorization, Api-Key"
	headers["Access-Control-Allow-Methods"] = "GET, OPTIONS"

	if acceptableOrigin(req.Headers["origin"]) {
		headers["Access-Control-Allow-Origin"] = req.Headers["origin"]
	}

	var resp events.APIGatewayProxyResponse

	user, ok := req.PathParameters["user"]
	if !ok {
		resp = events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "no user",
			Headers:    headers,
		}

		return resp, nil
	}

	user = "USER#" + strings.ToUpper(user)

	key := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: user},
		"SK": &types.AttributeValueMemberS{Value: "PROFILE"},
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String("SINGLE-TABLE"),
		Key:       key,
	}

	result, err := svc.GetItem(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to get item from DynamoDB, %v", err)
	}

	var userprofile UserProfile
	err = attributevalue.UnmarshalMap(result.Item, &userprofile)
	if err != nil {
		log.Fatalf("failed to unmarshal item into struct, %v", err)
	}

	userProfileJSON, err := json.Marshal(userprofile)
	if err != nil {
		log.Fatalf("failed to marshal struct into JSON, %v", err)
	}

	resp = events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(userProfileJSON),
		Headers:    headers,
	}

	return resp, nil
}

func main() {
	lambda.Start(HandleRequest)
}

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc = dynamodb.NewFromConfig(cfg)
}
