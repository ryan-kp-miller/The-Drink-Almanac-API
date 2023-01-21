package store

import (
	"context"
	"fmt"
	"strconv"
	"the-drink-almanac-api/client"
	"the-drink-almanac-api/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	USERS_TABLE_NAME string = "the-drink-almanac-users"
)

type UserStoreDDB struct {
	DynamodbClient client.DDBClient
}

func (usd *UserStoreDDB) FindAll() ([]model.User, error) {
	scanInput := dynamodb.ScanInput{
		TableName: aws.String(USERS_TABLE_NAME),
	}
	ctx := context.TODO()
	scanOutput, err := usd.DynamodbClient.Scan(ctx, &scanInput)
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

// FindUserByUsername checks the store's user table for any users with the given username;
//
// Returns:
//   - an error if there are multiple users with that username
//   - the user if there is only 1 user
//   - nil if there are 0 users
func (usd *UserStoreDDB) FindUserByUsername(username string) (*model.User, error) {
	filterExpression, err := expression.NewBuilder().WithFilter(
		expression.Equal(expression.Name("username"), expression.Value(username)),
	).Build()
	if err != nil {
		return nil, err
	}

	scanInput := dynamodb.ScanInput{
		TableName:                 aws.String(USERS_TABLE_NAME),
		FilterExpression:          filterExpression.Filter(),
		ExpressionAttributeNames:  filterExpression.Names(),
		ExpressionAttributeValues: filterExpression.Values(),
	}
	ctx := context.TODO()
	scanOutput, err := usd.DynamodbClient.Scan(ctx, &scanInput)
	if err != nil {
		return nil, err
	}
	users := []model.User{}
	err = attributevalue.UnmarshalListOfMaps(scanOutput.Items, &users)
	if err != nil {
		return nil, err
	}

	switch num_users := len(users); num_users {
	case 0:
		return nil, nil
	case 1:
		return &users[0], nil
	default:
		return nil, fmt.Errorf("there are %s users with the username %s", strconv.Itoa(len(users)), username)
	}
}

// CreateNewUser simply inserts the provided user into the store's user table
//
// Please ensure that you aren't inserting a duplicate record (i.e. user with that username already exists)
func (usd *UserStoreDDB) CreateNewUser(user model.User) error {
	_, err := usd.DynamodbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(USERS_TABLE_NAME),
		Item: map[string]types.AttributeValue{
			"id":       &types.AttributeValueMemberS{Value: user.Id},
			"username": &types.AttributeValueMemberS{Value: user.Username},
			"password": &types.AttributeValueMemberS{Value: user.Password},
		},
	})
	return err
}

// DeleteUser removes the record associated with the given id
// from the store's user table
func (usd *UserStoreDDB) DeleteUser(id string) error {
	_, err := usd.DynamodbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(USERS_TABLE_NAME),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	return err
}

func NewUserStoreDDB() (*UserStoreDDB, error) {
	ddbClient, err := client.CreateLocalDDBClient()
	return &UserStoreDDB{
		DynamodbClient: ddbClient,
	}, err
}
