package database

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

//CreateLocalClient creates a dynamodb client using environment variables
func CreateLocalClient() (*dynamodb.DynamoDB, error) {
	awsEndpoint := DefaultEnv("AWS_ENDPOINT", "http://localstack:4566")
	awsRegion := DefaultEnv("AWS_REGION", "us-east-1")

	fmt.Println(awsEndpoint)
	fmt.Println("Creating new aws session")
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(awsEndpoint),
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials("test", "test", ""),
	})
	if err != nil {
		return nil, err
	}
	// Create DynamoDB client
	svc := dynamodb.New(sess)
	return svc, nil
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
