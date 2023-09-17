package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Name string `json:"name"`
}

type Response struct {
	IsBase64Encoded bool              `json:"isBase64Encoded"`
	StatusCode      int               `json:"statusCode"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {

	fmt.Printf("Request!!!: %+v\n", request)

	if request.HTTPMethod == "OPTIONS" {
		return Response{
			IsBase64Encoded: false,
			StatusCode:      200,
			Headers: map[string]string{
				"Access-Control-Allow-Headers":  "*",
				"Access-Control-Allow-Methods":  "*",
				"Access-Control-Allow-Origin":   "*",
				"Access-Control-Expose-Headers": "*",
				"Access-Control-Max-Age":        "30",
			},
			Body: "",
		}, nil
	}

	myEvent := MyEvent{}
	if err := json.Unmarshal([]byte(request.Body), &myEvent); err != nil {
		return Response{
			IsBase64Encoded: false,
			StatusCode:      500,
		}, nil
	}

	body := fmt.Sprintf("{\"message\": \"Hello %s!\"}", myEvent.Name)

	return Response{
		IsBase64Encoded: false,
		StatusCode:      200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
		Body: body,
	}, nil

}

func main() {
	lambda.Start(HandleRequest)
}
