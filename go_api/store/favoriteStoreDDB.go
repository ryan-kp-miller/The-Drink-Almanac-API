package store

import (
	"context"
	"the-drink-almanac-api/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	FAVORITES_TABLE_NAME string = "the-drink-almanac-favorites"
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
		TableName:                 aws.String(FAVORITES_TABLE_NAME),
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
	_, err := frd.dynamodbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(FAVORITES_TABLE_NAME),
		Item: map[string]types.AttributeValue{
			"id":       &types.AttributeValueMemberS{Value: favorite.Id},
			"drink_id": &types.AttributeValueMemberS{Value: favorite.DrinkId},
			"user_id":  &types.AttributeValueMemberS{Value: favorite.UserId},
		},
	})
	return err
}

func NewFavoriteStoreDDB() (*FavoriteStoreDDB, error) {
	ddbClient, err := CreateLocalClient()
	return &FavoriteStoreDDB{
		dynamodbClient: ddbClient,
	}, err
}
