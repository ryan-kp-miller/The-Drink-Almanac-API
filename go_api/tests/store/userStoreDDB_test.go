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

func TestUserStoreDDB_FindAll(t *testing.T) {
	numUsers := 5
	userItems := make([]map[string]types.AttributeValue, numUsers)
	mockUsers := make([]model.User, numUsers)
	for i := 0; i < numUsers; i++ {
		iStr := strconv.Itoa(i)
		userItems[i] = map[string]types.AttributeValue{
			"id":       &types.AttributeValueMemberS{Value: iStr},
			"username": &types.AttributeValueMemberS{Value: iStr},
			"password": &types.AttributeValueMemberS{Value: iStr},
		}
		mockUsers[i] = model.User{
			Id:       iStr,
			Username: iStr,
			Password: iStr,
		}
	}
	scanOutput := &dynamodb.ScanOutput{Items: userItems}
	scanInput := &dynamodb.ScanInput{TableName: aws.String("")}
	tests := []struct {
		name          string
		expectedUsers []model.User
		returnedError error
		expectError   bool
	}{
		{
			name:          "Successfully retrieve users",
			expectedUsers: mockUsers,
			returnedError: nil,
			expectError:   false,
		},
		{
			name:          "Failed to retrieve users",
			expectedUsers: nil,
			returnedError: fmt.Errorf("failed to retrieve users"),
			expectError:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDdbClient := mocks.NewDDBClient(t)
			mockDdbClient.On("Scan", context.TODO(), scanInput).Return(scanOutput, tt.returnedError)
			userStore := store.UserStoreDDB{DynamodbClient: mockDdbClient}
			got, err := userStore.FindAll()
			if (err != nil) != tt.expectError {
				t.Errorf("UserStoreDDB.FindAll() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.expectedUsers) {
				t.Errorf("UserStoreDDB.FindAll() = %v, want %v", got, tt.expectedUsers)
			}
		})
	}
}

func TestUserStoreDDB_FindUserByUsername(t *testing.T) {
	numUsers := 5
	userItems := make([]map[string]types.AttributeValue, numUsers)
	mockUsers := make([]model.User, numUsers)
	for i := 0; i < numUsers; i++ {
		userItems[i] = map[string]types.AttributeValue{
			"id":       &types.AttributeValueMemberS{Value: "0"},
			"username": &types.AttributeValueMemberS{Value: "0"},
			"password": &types.AttributeValueMemberS{Value: "0"},
		}
		mockUsers[i] = model.User{
			Id:       "0",
			Username: "0",
			Password: "0",
		}
	}
	tests := []struct {
		name          string
		username      string
		expectedUser  *model.User
		scanOutput    *dynamodb.ScanOutput
		returnedError error
		expectError   bool
	}{
		{
			name:          "Successfully retrieve users",
			username:      "0",
			expectedUser:  &mockUsers[0],
			scanOutput:    &dynamodb.ScanOutput{Items: userItems[:1]},
			returnedError: nil,
			expectError:   false,
		},
		{
			name:          "Failed to retrieve users",
			username:      "0",
			expectedUser:  nil,
			returnedError: fmt.Errorf("failed to retrieve users"),
			expectError:   true,
		},
		{
			name:          "No existing user",
			username:      "0",
			expectedUser:  nil,
			scanOutput:    &dynamodb.ScanOutput{Items: nil},
			returnedError: nil,
			expectError:   false,
		},
		{
			name:          "Too many users",
			username:      "0",
			expectedUser:  nil,
			scanOutput:    &dynamodb.ScanOutput{Items: userItems},
			returnedError: nil,
			expectError:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDdbClient := mocks.NewDDBClient(t)
			mockDdbClient.On("Scan", context.TODO(), mock.AnythingOfType("*dynamodb.ScanInput")).Return(tt.scanOutput, tt.returnedError)
			userStore := store.UserStoreDDB{DynamodbClient: mockDdbClient}
			got, err := userStore.FindUserByUsername(tt.username)
			if (err != nil) != tt.expectError {
				t.Errorf("UserStoreDDB.FindUserByUsername() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.expectedUser) {
				t.Errorf("UserStoreDDB.FindUserByUsername() = %v, want %v", got, tt.expectedUser)
			}
		})
	}
}

func TestUserStoreDDB_CreateNewUser(t *testing.T) {
	mockUser := model.User{
		Id:       "0",
		Username: "0",
		Password: "0",
	}
	putItemInput := &dynamodb.PutItemInput{
		TableName: aws.String(""),
		Item: map[string]types.AttributeValue{
			"id":       &types.AttributeValueMemberS{Value: mockUser.Id},
			"username": &types.AttributeValueMemberS{Value: mockUser.Username},
			"password": &types.AttributeValueMemberS{Value: mockUser.Password},
		},
	}
	putItemOutput := &dynamodb.PutItemOutput{}
	tests := []struct {
		name          string
		expectedUser  model.User
		returnedError error
		expectError   bool
	}{
		{
			name:          "Successfully created a user",
			expectedUser:  mockUser,
			returnedError: nil,
			expectError:   false,
		},
		{
			name:          "Failed to create a user",
			expectedUser:  mockUser,
			returnedError: fmt.Errorf("failed to create the user"),
			expectError:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDdbClient := mocks.NewDDBClient(t)
			mockDdbClient.On("PutItem", context.TODO(), putItemInput).Return(putItemOutput, tt.returnedError)
			userStore := store.UserStoreDDB{DynamodbClient: mockDdbClient}
			err := userStore.CreateNewUser(tt.expectedUser)
			if (err != nil) != tt.expectError {
				t.Errorf("UserStoreDDB.CreateNewUser() error = %v", err)
				return
			}
		})
	}
}

func TestUserStoreDDB_DeleteUser(t *testing.T) {
	mockUser := model.User{
		Id:       "0",
		Username: "0",
		Password: "0",
	}
	deleteItemInput := &dynamodb.DeleteItemInput{
		TableName: aws.String(""),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: mockUser.Id},
		},
	}
	deleteItemOutput := &dynamodb.DeleteItemOutput{}
	tests := []struct {
		name          string
		expectedUser  model.User
		returnedError error
		expectError   bool
	}{
		{
			name:          "Successfully deleted the user",
			expectedUser:  mockUser,
			returnedError: nil,
			expectError:   false,
		},
		{
			name:          "Failed to create a user",
			expectedUser:  mockUser,
			returnedError: fmt.Errorf("failed to delete the user"),
			expectError:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDdbClient := mocks.NewDDBClient(t)
			mockDdbClient.On("DeleteItem", context.TODO(), deleteItemInput).Return(deleteItemOutput, tt.returnedError)
			userStore := store.UserStoreDDB{DynamodbClient: mockDdbClient}
			err := userStore.DeleteUser(tt.expectedUser.Id)
			if (err != nil) != tt.expectError {
				t.Errorf("UserStoreDDB.DeleteUser() error = %v", err)
				return
			}
		})
	}
}
