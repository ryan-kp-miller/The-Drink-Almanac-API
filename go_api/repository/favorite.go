//go:generate mockery --name=FavoriteRepository --output=./ --outpkg=repository --filename=favorite_mock.go --inpackage
package repository

import (
	"context"

	"the-drink-almanac-api/model"
	"the-drink-almanac-api/repository/client"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type FavoriteRepository interface {
	FindAll() ([]model.Favorite, error)
	FindFavoritesByUser(userId string) ([]model.Favorite, error)
	CreateNewFavorite(favorite model.Favorite) error
	DeleteFavorite(id string) error
}

func NewFavoriteRepository(tableName, awsEndpoint string) (*FavoriteRepositoryDDB, error) {
	ddbClient, err := client.CreateLocalDDBClient(awsEndpoint)
	return &FavoriteRepositoryDDB{
		DynamodbClient: ddbClient,
		TableName:      tableName,
	}, err
}

type FavoriteRepositoryDDB struct {
	DynamodbClient client.DDBClient
	TableName      string
}

func (r *FavoriteRepositoryDDB) FindAll() ([]model.Favorite, error) {
	scanInput := dynamodb.ScanInput{
		TableName: aws.String(r.TableName),
	}
	ctx := context.TODO()
	scanOutput, err := r.DynamodbClient.Scan(ctx, &scanInput)
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

func (r *FavoriteRepositoryDDB) FindFavoritesByUser(userId string) ([]model.Favorite, error) {
	filterExpression, err := expression.NewBuilder().WithKeyCondition(
		expression.Key("user_id").Equal(expression.Value(userId)),
	).Build()
	if err != nil {
		return nil, err
	}

	queryInput := dynamodb.QueryInput{
		TableName:                 aws.String(r.TableName),
		IndexName:                 aws.String("user-index"),
		ExpressionAttributeNames:  filterExpression.Names(),
		ExpressionAttributeValues: filterExpression.Values(),
		KeyConditionExpression:    filterExpression.KeyCondition(),
	}

	ctx := context.TODO()
	queryOutput, err := r.DynamodbClient.Query(ctx, &queryInput)
	if err != nil {
		return nil, err
	}
	favorites := []model.Favorite{}
	err = attributevalue.UnmarshalListOfMaps(queryOutput.Items, &favorites)
	if err != nil {
		return nil, err
	}

	return favorites, nil
}

func (r *FavoriteRepositoryDDB) CreateNewFavorite(favorite model.Favorite) error {
	_, err := r.DynamodbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.TableName),
		Item: map[string]types.AttributeValue{
			"id":       &types.AttributeValueMemberS{Value: favorite.Id},
			"drink_id": &types.AttributeValueMemberS{Value: favorite.DrinkId},
			"user_id":  &types.AttributeValueMemberS{Value: favorite.UserId},
		},
	})
	return err
}

func (r *FavoriteRepositoryDDB) DeleteFavorite(id string) error {
	_, err := r.DynamodbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(r.TableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	return err
}
