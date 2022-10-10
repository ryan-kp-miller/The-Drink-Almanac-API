package store

import (
	"context"
	"the-drink-almanac-api/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type UserStoreDDB struct {
	dynamodbClient *dynamodb.Client
}

func (urd *UserStoreDDB) FindAll() ([]model.User, error) {
	scanInput := dynamodb.ScanInput{
		TableName: aws.String("the-drink-almanac-users"),
	}
	ctx := context.TODO()
	scanOutput, err := urd.dynamodbClient.Scan(ctx, &scanInput)
	if err != nil {
		return nil, err
	}
	users := []model.User{}
	err = attributevalue.UnmarshalListOfMaps(scanOutput.Items, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func NewUserStoreDDB() (*UserStoreDDB, error) {
	ddbClient, err := CreateLocalClient()
	return &UserStoreDDB{
		dynamodbClient: ddbClient,
	}, err
}
