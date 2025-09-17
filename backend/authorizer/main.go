package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var cloudFrontSecret = os.Getenv("CLOUDFRONT_SECRET")

func handler(ctx context.Context, event events.APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	key := event.Headers["x-cloudfront-secret"]
	log.Println("env key:", cloudFrontSecret)
	log.Println("header key:", key)

	if key != cloudFrontSecret {
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{
			IsAuthorized: false,
		}, nil
	}

	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: true,
		Context: map[string]interface{}{
			"user": "cloudfront",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
