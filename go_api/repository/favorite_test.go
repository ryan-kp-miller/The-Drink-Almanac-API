package repository

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"the-drink-almanac-api/mocks"
	"the-drink-almanac-api/model"
)

func TestFavoriteStoreDDB_FindAll(t *testing.T) {
	numFavorites := 5
	favoriteItems := make([]map[string]types.AttributeValue, numFavorites)
	mockFavorites := make([]model.Favorite, numFavorites)
	for i := 0; i < numFavorites; i++ {
		iStr := strconv.Itoa(i)
		favoriteItems[i] = map[string]types.AttributeValue{
			"id":       &types.AttributeValueMemberS{Value: iStr},
			"user_id":  &types.AttributeValueMemberS{Value: iStr},
			"drink_id": &types.AttributeValueMemberS{Value: iStr},
		}
		mockFavorites[i] = model.Favorite{
			Id:      iStr,
			DrinkId: iStr,
			UserId:  iStr,
		}
	}
	scanOutput := &dynamodb.ScanOutput{Items: favoriteItems}
	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(""),
	}
	tests := []struct {
		name              string
		expectedFavorites []model.Favorite
		returnedError     error
		expectError       bool
	}{
		{
			name:              "Successfully retrieve favorites",
			expectedFavorites: mockFavorites,
			returnedError:     nil,
			expectError:       false,
		},
		{
			name:              "Failed to retrieve favorites",
			expectedFavorites: nil,
			returnedError:     fmt.Errorf("failed to retrieve favorites"),
			expectError:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDdbClient := mocks.NewDDBClient(t)
			mockDdbClient.On("Scan", context.TODO(), scanInput).Return(scanOutput, tt.returnedError)
			favoriteStore := FavoriteRepositoryDDB{DynamodbClient: mockDdbClient}
			got, err := favoriteStore.FindAll()
			if (err != nil) != tt.expectError {
				t.Errorf("FavoriteRepository.FindAll() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.expectedFavorites) {
				t.Errorf("FavoriteRepository.FindAll() = %v, want %v", got, tt.expectedFavorites)
			}
		})
	}
}

func TestFavoriteStoreDDB_FindFavoritesByUser(t *testing.T) {
	numFavorites := 5
	favoriteItems := make([]map[string]types.AttributeValue, numFavorites)
	mockFavorites := make([]model.Favorite, numFavorites)
	for i := 0; i < numFavorites; i++ {
		iStr := strconv.Itoa(i)
		favoriteItems[i] = map[string]types.AttributeValue{
			"id":       &types.AttributeValueMemberS{Value: iStr},
			"user_id":  &types.AttributeValueMemberS{Value: "0"},
			"drink_id": &types.AttributeValueMemberS{Value: iStr},
		}
		mockFavorites[i] = model.Favorite{
			Id:      iStr,
			UserId:  "0",
			DrinkId: iStr,
		}
	}
	tests := []struct {
		name              string
		userId            string
		expectedFavorites []model.Favorite
		scanOutput        *dynamodb.QueryOutput
		returnedError     error
		expectError       bool
	}{
		{
			name:              "Successfully retrieve favorites",
			userId:            "0",
			expectedFavorites: mockFavorites,
			scanOutput:        &dynamodb.QueryOutput{Items: favoriteItems},
			returnedError:     nil,
			expectError:       false,
		},
		{
			name:              "Failed to retrieve favorites",
			userId:            "0",
			expectedFavorites: nil,
			returnedError:     fmt.Errorf("failed to retrieve favorites"),
			expectError:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDdbClient := mocks.NewDDBClient(t)
			mockDdbClient.On("Query", context.TODO(), mock.AnythingOfType("*dynamodb.QueryInput")).Return(tt.scanOutput, tt.returnedError)
			favoriteStore := FavoriteRepositoryDDB{DynamodbClient: mockDdbClient}
			actualFavorites, err := favoriteStore.FindFavoritesByUser(tt.userId)
			assert.Equal(t, tt.expectError, err != nil, "FavoriteRepository.FindFavoritesByUser() error = %v", err)
			assert.Equal(t, tt.expectedFavorites, actualFavorites, "FavoriteRepository.FindFavoritesByUser() = %v, want %v", actualFavorites, tt.expectedFavorites)
		})
	}
}

func TestFavoriteStoreDDB_CreateNewFavorite(t *testing.T) {
	mockFavorite := model.Favorite{
		Id:      "0",
		DrinkId: "0",
		UserId:  "0",
	}
	putItemInput := &dynamodb.PutItemInput{
		TableName: aws.String(""),
		Item: map[string]types.AttributeValue{
			"id":       &types.AttributeValueMemberS{Value: mockFavorite.Id},
			"user_id":  &types.AttributeValueMemberS{Value: mockFavorite.DrinkId},
			"drink_id": &types.AttributeValueMemberS{Value: mockFavorite.UserId},
		},
	}
	putItemOutput := &dynamodb.PutItemOutput{}
	tests := []struct {
		name             string
		expectedFavorite model.Favorite
		returnedError    error
		expectError      bool
	}{
		{
			name:             "Successfully created a favorite",
			expectedFavorite: mockFavorite,
			returnedError:    nil,
			expectError:      false,
		},
		{
			name:             "Failed to create a favorite",
			expectedFavorite: mockFavorite,
			returnedError:    fmt.Errorf("failed to create the favorite"),
			expectError:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDdbClient := mocks.NewDDBClient(t)
			mockDdbClient.On("PutItem", context.TODO(), putItemInput).Return(putItemOutput, tt.returnedError)
			favoriteStore := FavoriteRepositoryDDB{DynamodbClient: mockDdbClient}
			err := favoriteStore.CreateNewFavorite(tt.expectedFavorite)
			if (err != nil) != tt.expectError {
				t.Errorf("FavoriteRepository.CreateNewFavorite() error = %v", err)
				return
			}
		})
	}
}

func TestFavoriteStoreDDB_DeleteFavorite(t *testing.T) {
	mockFavorite := model.Favorite{
		Id:      "0",
		DrinkId: "0",
		UserId:  "0",
	}
	deleteItemInput := &dynamodb.DeleteItemInput{
		TableName: aws.String(""),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: mockFavorite.Id},
		},
	}
	deleteItemOutput := &dynamodb.DeleteItemOutput{}
	tests := []struct {
		name             string
		expectedFavorite model.Favorite
		returnedError    error
		expectError      bool
	}{
		{
			name:             "Successfully deleted the favorite",
			expectedFavorite: mockFavorite,
			returnedError:    nil,
			expectError:      false,
		},
		{
			name:             "Failed to create a favorite",
			expectedFavorite: mockFavorite,
			returnedError:    fmt.Errorf("failed to delete the favorite"),
			expectError:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDdbClient := mocks.NewDDBClient(t)
			mockDdbClient.On("DeleteItem", context.TODO(), deleteItemInput).Return(deleteItemOutput, tt.returnedError)
			favoriteStore := FavoriteRepositoryDDB{DynamodbClient: mockDdbClient}
			err := favoriteStore.DeleteFavorite(tt.expectedFavorite.Id)
			if (err != nil) != tt.expectError {
				t.Errorf("FavoriteRepository.DeleteFavorite() error = %v", err)
				return
			}
		})
	}
}
