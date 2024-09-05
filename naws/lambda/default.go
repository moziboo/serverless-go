package lambda

import (
	"context"
	"naws/models"
	"naws/utils"

	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var svc *dynamodb.Client

func User(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := utils.HandleCORS(req)

	key := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: "USER#KELLY"},
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

	var userprofile models.UserProfile
	err = attributevalue.UnmarshalMap(result.Item, &userprofile)
	if err != nil {
		log.Fatalf("failed to unmarshal item into struct, %v", err)
	}

	userProfileJSON, err := json.Marshal(userprofile)
	if err != nil {
		log.Fatalf("failed to marshal struct into JSON, %v", err)
	}

	resp := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(userProfileJSON),
		Headers:    headers,
	}

	return resp, nil
}

func Posts(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := utils.HandleCORS(req)

	input := &dynamodb.QueryInput{
		TableName:              aws.String("SINGLE-TABLE"),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :skPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":       &types.AttributeValueMemberS{Value: "USER#KELLY"},
			":skPrefix": &types.AttributeValueMemberS{Value: "POST#"},
		},
	}

	result, err := svc.Query(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to query items from DynamoDB, %v", err)
	}

	var posts []models.Post
	err = attributevalue.UnmarshalListOfMaps(result.Items, &posts)
	if err != nil {
		log.Fatalf("failed to unmarshal items into struct, %v", err)
	}

	postsJSON, err := json.Marshal(posts)
	if err != nil {
		log.Fatalf("failed to marshal struct into JSON, %v", err)
	}

	resp := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(postsJSON),
		Headers:    headers,
	}

	return resp, nil
}

func init() {
	// Load the Shared AWS Configuration (e.g., ~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc = dynamodb.NewFromConfig(cfg)
}
