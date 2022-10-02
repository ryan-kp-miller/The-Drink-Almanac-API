package database

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

//CreateLocalClient creates a dynamodb client using environment variables
func CreateLocalClient() (*dynamodb.Client, error) {
	awsEndpoint := DefaultEnv("AWS_ENDPOINT", "http://localstack:4566")
	awsRegion := DefaultEnv("AWS_REGION", "us-east-1")

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == dynamodb.ServiceID && region == "us-east-1" {
			fmt.Println("aws custom endpoint resolver called correctly")
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               awsEndpoint,
				SigningRegion:     awsRegion,
				HostnameImmutable: true,
			}, nil
		}
		// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))

	fmt.Println("Initializing config")
	// cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	fmt.Println("Initializing DDB service")
	return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.EndpointResolver = dynamodb.EndpointResolverFromURL(awsEndpoint)
	}), nil
}

// func CreateSession() (*session.Session, error) {
// 	awsEndpoint := DefaultEnv("test", "http://localhost:4566")
// 	awsRegion := DefaultEnv("AWS_REGION", "us-east-1")

// 	cfg, err := config.LoadDefaultConfig(context.TODO(),
// 		config.WithRegion(awsRegion),
// 		config.WithClientLogMode(aws.LogRequest|aws.LogRetries),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	cfg.Endpoint = aws.String(awsEndpoint)
// 	return session.New(cfg), nil
// 	// return dynamodb.NewFromConfig(cfg), nil
// }

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
