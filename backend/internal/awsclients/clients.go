package awsclients

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type Clients struct {
	DynamoDB *dynamodb.Client
	SSM      *ssm.Client
}

func New() (*Clients, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	return &Clients{
		DynamoDB: dynamodb.NewFromConfig(cfg),
		SSM:      ssm.NewFromConfig(cfg),
	}, nil
}
