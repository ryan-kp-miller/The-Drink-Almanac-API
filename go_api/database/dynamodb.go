package database

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

//CreateLocalClient creates a dynamodb client using environment variables
func CreateLocalClient() (*dynamodb.Client, error) {
	awsEndpoint := DefaultEnv("AWS_ENDPOINT", "http://localstack:4566")
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.EndpointResolver = dynamodb.EndpointResolverFromURL(awsEndpoint)
	}), nil
}

// DefaultEnv takes the name of the environment variable and a default value;
// if the environment variable wasn't found, then the default value is returned;
//
// Note: if the environment variable exists but just contains an empty string,
// the empty string will be returned
func DefaultEnv(envVarName, defaultValue string) string {
	envValue, ok := os.LookupEnv(envVarName)
	if !ok {
		envValue = defaultValue
	}
	return envValue
}
