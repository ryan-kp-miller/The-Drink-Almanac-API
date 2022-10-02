package domain

import (
	"context"
	"the-drink-almanac-api/database"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
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

func NewFavoriteRepositoryDDB() (*FavoriteRepositoryDDB, error) {
	ddbClient, err := database.CreateLocalClient()
	return &FavoriteRepositoryDDB{
		dynamodbClient: ddbClient,
	}, err
}
