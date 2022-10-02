package domain

import (
	"fmt"
	"log"
	"the-drink-almanac-api/database"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type FavoriteRepositoryDDB struct {
	dynamodbClient *dynamodb.DynamoDB
}

func (frd *FavoriteRepositoryDDB) FindAll() ([]Favorite, error) {
	scanInput := dynamodb.ScanInput{
		TableName: aws.String("the-drink-almanac-favorites"),
	}
	scanOutput, err := frd.dynamodbClient.Scan(&scanInput)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(scanOutput.Items)
	favorites := []Favorite{}
	// err = attributevalue.UnmarshalListOfMaps(scanOutput.Items, &favorites)
	// if err != nil {
	// 	return nil, err
	// }

	return favorites, nil
}

func NewFavoriteRepositoryDDB() (*FavoriteRepositoryDDB, error) {
	ddbClient, err := database.CreateLocalClient()
	return &FavoriteRepositoryDDB{
		dynamodbClient: ddbClient,
	}, err
}
