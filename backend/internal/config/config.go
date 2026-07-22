package config

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type Config struct {
	AuthServiceURL string
	DynamoDBTable  string
}

func Load(ctx context.Context, ssmClient *ssm.Client) (*Config, error) {
	authURL, err := getParameter(ctx, ssmClient, "/leoding/auth/service_url")
	if err != nil {
		return nil, fmt.Errorf("load auth service url: %w", err)
	}

	tableName, err := getParameter(ctx, ssmClient, "/leoding/dynamodb/coffee")
	if err != nil {
		return nil, fmt.Errorf("load dynamodb table: %w", err)
	}

	return &Config{
		AuthServiceURL: authURL,
		DynamoDBTable:  tableName,
	}, nil
}

func getParameter(ctx context.Context, client *ssm.Client, name string) (string, error) {
	resp, err := client.GetParameter(ctx, &ssm.GetParameterInput{
		Name: aws.String(name),
	})
	if err != nil {
		return "", err
	}

	return aws.ToString(resp.Parameter.Value), nil
}
