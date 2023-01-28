package store

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"the-drink-almanac-api/mocks"
	"the-drink-almanac-api/model"
	"the-drink-almanac-api/store"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/mock"
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
			favoriteStore := store.FavoriteStoreDDB{DynamodbClient: mockDdbClient}
			got, err := favoriteStore.FindAll()
			if (err != nil) != tt.expectError {
				t.Errorf("FavoriteStoreDDB.FindAll() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.expectedFavorites) {
				t.Errorf("FavoriteStoreDDB.FindAll() = %v, want %v", got, tt.expectedFavorites)
			}
		})
	}
}

func TestFavoriteStoreDDB_FindFavoritesByUser(t *testing.T) {
	numFavorites := 5
	favoriteItems := make([]map[string]types.AttributeValue, numFavorites)
	mockFavorites := make([]model.Favorite, numFavorites)
	for i := 0; i < numFavorites; i++ {
		favoriteItems[i] = map[string]types.AttributeValue{
			"id":       &types.AttributeValueMemberS{Value: "0"},
			"user_id":  &types.AttributeValueMemberS{Value: "0"},
			"drink_id": &types.AttributeValueMemberS{Value: "0"},
		}
		mockFavorites[i] = model.Favorite{
			Id:      "0",
			DrinkId: "0",
			UserId:  "0",
		}
	}
	tests := []struct {
		name             string
		userId           string
		expectedFavorite []model.Favorite
		scanOutput       *dynamodb.ScanOutput
		returnedError    error
		expectError      bool
	}{
		{
			name:             "Successfully retrieve favorites",
			userId:           "0",
			expectedFavorite: mockFavorites,
			scanOutput:       &dynamodb.ScanOutput{Items: favoriteItems},
			returnedError:    nil,
			expectError:      false,
		},
		{
			name:             "Failed to retrieve favorites",
			userId:           "0",
			expectedFavorite: nil,
			returnedError:    fmt.Errorf("failed to retrieve favorites"),
			expectError:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDdbClient := mocks.NewDDBClient(t)
			mockDdbClient.On("Scan", context.TODO(), mock.AnythingOfType("*dynamodb.ScanInput")).Return(tt.scanOutput, tt.returnedError)
			favoriteStore := store.FavoriteStoreDDB{DynamodbClient: mockDdbClient}
			got, err := favoriteStore.FindFavoritesByUser(tt.userId)
			if (err != nil) != tt.expectError {
				t.Errorf("FavoriteStoreDDB.FindFavoritesByUser() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.expectedFavorite) {
				t.Errorf("FavoriteStoreDDB.FindFavoritesByUser() = %v, want %v", got, tt.expectedFavorite)
			}
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
			favoriteStore := store.FavoriteStoreDDB{DynamodbClient: mockDdbClient}
			err := favoriteStore.CreateNewFavorite(tt.expectedFavorite)
			if (err != nil) != tt.expectError {
				t.Errorf("FavoriteStoreDDB.CreateNewFavorite() error = %v", err)
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
			favoriteStore := store.FavoriteStoreDDB{DynamodbClient: mockDdbClient}
			err := favoriteStore.DeleteFavorite(tt.expectedFavorite.Id)
			if (err != nil) != tt.expectError {
				t.Errorf("FavoriteStoreDDB.DeleteFavorite() error = %v", err)
				return
			}
		})
	}
}
