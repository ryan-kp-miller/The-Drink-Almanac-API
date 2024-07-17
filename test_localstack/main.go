package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	awsEndpoint := "http://localhost:4566"
	awsRegion := "us-east-1"
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(awsEndpoint),
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials("test", "test", ""),
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create DynamoDB client
	svc := dynamodb.New(sess)

	scanInput := dynamodb.ScanInput{
		TableName: aws.String("the-drink-almanac-favorites"),
	}
	scanOutput, err := svc.Scan(&scanInput)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(scanOutput.Items)
}
