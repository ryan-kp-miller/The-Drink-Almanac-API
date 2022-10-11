package store

import (
	"context"
	"fmt"
	"the-drink-almanac-api/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type FavoriteStoreDDB struct {
	dynamodbClient *dynamodb.Client
}

func (frd *FavoriteStoreDDB) FindAll() ([]model.Favorite, error) {
	scanInput := dynamodb.ScanInput{
		TableName: aws.String("the-drink-almanac-favorites"),
	}
	ctx := context.TODO()
	scanOutput, err := frd.dynamodbClient.Scan(ctx, &scanInput)
	if err != nil {
		return nil, err
	}
	favorites := []model.Favorite{}
	err = attributevalue.UnmarshalListOfMaps(scanOutput.Items, &favorites)
	if err != nil {
		return nil, err
	}

	return favorites, nil
}

func (frd *FavoriteStoreDDB) FindFavoritesByUser(userId string) ([]model.Favorite, error) {
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
	favorites := []model.Favorite{}
	err = attributevalue.UnmarshalListOfMaps(scanOutput.Items, &favorites)
	if err != nil {
		return nil, err
	}

	return favorites, nil
}

func (frd *FavoriteStoreDDB) CreateNewFavorite(favorite model.Favorite) error {
	return fmt.Errorf("not implemented yet")
}

func NewFavoriteStoreDDB() (*FavoriteStoreDDB, error) {
	ddbClient, err := CreateLocalClient()
	return &FavoriteStoreDDB{
		dynamodbClient: ddbClient,
	}, err
}
