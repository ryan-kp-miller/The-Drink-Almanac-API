package domain

import (
	"context"
	"fmt"
	"the-drink-almanac-api/database"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type FavoriteRepositoryDDB struct {
	dynamodbClient *dynamodb.Client
}

func (frd *FavoriteRepositoryDDB) FindAll() ([]Favorite, error) {
	scanInput := dynamodb.ScanInput{
		TableName: aws.String("the-drink-almanac-favorites"),
	}
	ctx := context.TODO()
	scanOutput, err := frd.dynamodbClient.Scan(ctx, &scanInput)
	if err != nil {
		return nil, err
	}
	favorites := []Favorite{}
	err = attributevalue.UnmarshalListOfMaps(scanOutput.Items, &favorites)
	if err != nil {
		return nil, err
	}

	return favorites, nil
}

func (frd *FavoriteRepositoryDDB) FindFavoritesByUser(userId int) ([]Favorite, error) {
	filterExpression, err := expression.NewBuilder().WithFilter(
		expression.Equal(expression.Name("user_id"), expression.Value(userId)),
	).Build()
	if err != nil {
		return nil, err
	}

	scanInput := dynamodb.ScanInput{
		TableName:                 aws.String("the-drink-almanac-favorites"),
		FilterExpression:          filterExpression.Filter(),
		ExpressionAttributeNames:  filterExpression.Names(),
		ExpressionAttributeValues: filterExpression.Values(),
	}
	ctx := context.TODO()
	scanOutput, err := frd.dynamodbClient.Scan(ctx, &scanInput)
	if err != nil {
		return nil, err
	}
	favorites := []Favorite{}
	err = attributevalue.UnmarshalListOfMaps(scanOutput.Items, &favorites)
	if err != nil {
		return nil, err
	}

	return favorites, nil
}

func (frd *FavoriteRepositoryDDB) CreateNewFavorite(favorite Favorite) error {
	return fmt.Errorf("not implemented yet")
}

func NewFavoriteRepositoryDDB() (*FavoriteRepositoryDDB, error) {
	ddbClient, err := database.CreateLocalClient()
	return &FavoriteRepositoryDDB{
		dynamodbClient: ddbClient,
	}, err
}
