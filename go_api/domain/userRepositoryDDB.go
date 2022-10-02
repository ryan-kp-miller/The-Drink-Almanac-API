package domain

import (
	"context"
	"the-drink-almanac-api/database"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type UserRepositoryDDB struct {
	dynamodbClient *dynamodb.Client
}

func (urd *UserRepositoryDDB) FindAll() ([]User, error) {
	scanInput := dynamodb.ScanInput{
		TableName: aws.String("the-drink-almanac-users"),
	}
	ctx := context.TODO()
	scanOutput, err := urd.dynamodbClient.Scan(ctx, &scanInput)
	if err != nil {
		return nil, err
	}
	users := []User{}
	err = attributevalue.UnmarshalListOfMaps(scanOutput.Items, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func NewUserRepositoryDDB() (*UserRepositoryDDB, error) {
	ddbClient, err := database.CreateLocalClient()
	return &UserRepositoryDDB{
		dynamodbClient: ddbClient,
	}, err
}
