package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

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

	resp := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "We prolly going to do GO things here",
		Headers:    headers,
	}

	return resp, nil
}

func main() {
	lambda.Start(HandleRequest)
}
