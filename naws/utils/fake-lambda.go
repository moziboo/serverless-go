package utils

import (
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type LambdaHandlerFunc func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func APIGatewayRequester(r *http.Request) (events.APIGatewayProxyRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return events.APIGatewayProxyRequest{}, err
	}
	defer r.Body.Close()

	headers := make(map[string]string)
	for key, values := range r.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	queryParams := make(map[string]string)
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			queryParams[key] = values[0]
		}
	}

	return events.APIGatewayProxyRequest{
		HTTPMethod:            r.Method,
		Headers:               headers,
		QueryStringParameters: queryParams,
		Path:                  r.URL.Path,
		Body:                  string(body),
		IsBase64Encoded:       false,
	}, nil
}

func APIGatewayResponder(lr events.APIGatewayProxyResponse, w http.ResponseWriter) {
	for key, value := range lr.Headers {
		w.Header().Set(key, value)
	}
	w.WriteHeader(lr.StatusCode)
	w.Write([]byte(lr.Body))
}

func LambdaProxy(handler LambdaHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiRequest, err := APIGatewayRequester(r)
		if err != nil {
			log.Printf("Error converting request: %v", err)
			http.Error(w, "Failed to convert request", http.StatusInternalServerError)
			return
		}

		lambdaResponse, err := handler(apiRequest)
		if err != nil {
			log.Printf("Error handling request: %v", err)
			http.Error(w, "Failed to handle request", http.StatusInternalServerError)
			return
		}

		APIGatewayResponder(lambdaResponse, w)
	}
}
