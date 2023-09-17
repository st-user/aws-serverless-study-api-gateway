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
	IsBase64Encoded bool   `json:"isBase64Encoded"`
	StatusCode      int    `json:"statusCode"`
	Body            string `json:"body"`
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {

	fmt.Printf("Request!!!: %+v\n", request)

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
		Body:            body,
	}, nil

}

func main() {
	lambda.Start(HandleRequest)
}
