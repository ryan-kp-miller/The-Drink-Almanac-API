package service

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"the-drink-almanac-api/mocks"
	"the-drink-almanac-api/model"
)

func TestDefaultFavoriteService_FindAllFavorites(t *testing.T) {
	mockFavorites := []model.Favorite{
		{
			Id:      "0",
			UserId:  "0",
			DrinkId: "0",
		},
		{
			Id:      "1",
			UserId:  "1",
			DrinkId: "1",
		},
		{
			Id:      "2",
			UserId:  "2",
			DrinkId: "2",
		},
	}
	tests := []struct {
		name              string
		returnedFavorites []model.Favorite
		returnedError     error
		expectError       bool
	}{
		{
			name:              "Successfully retrieved favorites",
			returnedFavorites: mockFavorites,
			returnedError:     nil,
			expectError:       false,
		},
		{
			name:              "Failed to retrieve favorites",
			returnedFavorites: mockFavorites,
			returnedError:     fmt.Errorf("failed to retrieve favorites"),
			expectError:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFavoriteStore := mocks.NewFavoriteStore(t)
			mockFavoriteStore.On("FindAll").Return(tt.returnedFavorites, tt.returnedError)
			favoriteService := NewDefaultFavoriteService(mockFavoriteStore)
			favorites, err := favoriteService.FindAllFavorites()
			assert.Equal(t, favorites, tt.returnedFavorites, "The favorites returned by FindAllFavorites() does not match the expected favorites; actual favorites = %v; expected favorites = %v", favorites, tt.returnedFavorites)
			if tt.expectError {
				assert.NotNil(t, err, "An error should have been returned from favoriteService.FindAllFavorites")
			} else {
				assert.Nil(t, err, "No error should have been returned from favoriteService.FindAllFavorites")
			}
			mockFavoriteStore.AssertExpectations(t)
		})
	}
}

func TestDefaultFavoriteService_CreateNewFavorite(t *testing.T) {
	mockFavorites := []model.Favorite{
		{
			Id:      "0",
			UserId:  "0",
			DrinkId: "0",
		},
		{
			Id:      "1",
			UserId:  "1",
			DrinkId: "1",
		},
		{
			Id:      "2",
			UserId:  "2",
			DrinkId: "2",
		},
	}
	tests := []struct {
		name                            string
		userId                          string
		drinkId                         string
		isStoreCreateNewFavoriteCalled  bool
		isStoreFindFavoriteByUserCalled bool
		returnedError                   error
		existingFavorites               []model.Favorite
		existingFavoritesError          error
		expectError                     bool
	}{
		{
			name:                            "Successfully create favorite",
			userId:                          "-1",
			drinkId:                         "0",
			isStoreCreateNewFavoriteCalled:  true,
			isStoreFindFavoriteByUserCalled: true,
			returnedError:                   nil,
			existingFavorites:               mockFavorites,
			existingFavoritesError:          nil,
			expectError:                     false,
		},
		{
			name:                            "Failed to create favorites",
			userId:                          "-1",
			drinkId:                         "0",
			isStoreCreateNewFavoriteCalled:  true,
			isStoreFindFavoriteByUserCalled: true,
			returnedError:                   fmt.Errorf("failed to create favorites"),
			existingFavorites:               mockFavorites,
			existingFavoritesError:          nil,
			expectError:                     true,
		},
		{
			name:                            "Favorite already exists",
			userId:                          "0",
			drinkId:                         "0",
			isStoreCreateNewFavoriteCalled:  false,
			isStoreFindFavoriteByUserCalled: true,
			returnedError:                   nil,
			existingFavorites:               mockFavorites,
			existingFavoritesError:          fmt.Errorf("favorite already exists"),
			expectError:                     true,
		},
		{
			name:                            "UserId is empty",
			userId:                          "",
			drinkId:                         "0",
			isStoreCreateNewFavoriteCalled:  false,
			isStoreFindFavoriteByUserCalled: false,
			returnedError:                   nil,
			existingFavorites:               nil,
			existingFavoritesError:          nil,
			expectError:                     true,
		},
		{
			name:                            "DrinkId is empty",
			userId:                          "0",
			drinkId:                         "",
			isStoreCreateNewFavoriteCalled:  false,
			isStoreFindFavoriteByUserCalled: false,
			returnedError:                   nil,
			existingFavorites:               nil,
			existingFavoritesError:          nil,
			expectError:                     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFavoriteStore := mocks.NewFavoriteStore(t)
			if tt.isStoreFindFavoriteByUserCalled {
				mockFavoriteStore.On("FindFavoritesByUser", tt.userId).Return(tt.existingFavorites, tt.existingFavoritesError)
			}
			if tt.isStoreCreateNewFavoriteCalled {
				mockFavoriteStore.On("CreateNewFavorite", mock.AnythingOfType("model.Favorite")).Return(tt.returnedError)
			}

			favoriteService := NewDefaultFavoriteService(mockFavoriteStore)
			favorite, err := favoriteService.CreateNewFavorite(tt.drinkId, tt.userId)
			if tt.expectError {
				assert.NotNil(t, err, "An error should have been returned from favoriteService.CreateNewFavorite")
			} else {
				assert.Nil(t, err, "No error should have been returned from favoriteService.CreateNewFavorite")
			}

			// check that the new favorite has a uuid for Id, UserId matches the provided userId,
			// and the DrinkId is the hashed version of the provided drinkId
			if favorite != nil {
				_, err = uuid.Parse(favorite.Id)
				assert.Nil(t, err, fmt.Sprintf("The favorite.Id field must be a uuid; actual value: %s", favorite.Id))
				assert.Equal(t, favorite.UserId, tt.userId)
				assert.Equal(t, favorite.DrinkId, tt.drinkId)
			}
			mockFavoriteStore.AssertExpectations(t)
		})
	}
}

func TestDefaultFavoriteService_FindFavoritesByUser(t *testing.T) {
	mockFavorites := []model.Favorite{
		{
			Id:      "0",
			UserId:  "0",
			DrinkId: "0",
		},
		{
			Id:      "1",
			UserId:  "0",
			DrinkId: "1",
		},
	}
	tests := []struct {
		name              string
		id                string
		returnedFavorites []model.Favorite
		returnedError     error
		expectError       bool
	}{
		{
			name:              "Successfully retrieved favorites",
			id:                "0",
			returnedFavorites: mockFavorites,
			returnedError:     nil,
			expectError:       false,
		},
		{
			name:              "Failed to retrieve favorites",
			id:                "0",
			returnedFavorites: nil,
			returnedError:     fmt.Errorf("failed to retrieve favorites"),
			expectError:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFavoriteStore := mocks.NewFavoriteStore(t)
			mockFavoriteStore.On("FindFavoritesByUser", tt.id).Return(tt.returnedFavorites, tt.returnedError)
			favoriteService := NewDefaultFavoriteService(mockFavoriteStore)
			favorites, err := favoriteService.FindFavoritesByUser(tt.id)
			assert.Equal(t, favorites, tt.returnedFavorites)
			if tt.expectError {
				assert.NotNil(t, err, "An error should have been returned from favoriteService.FindFavoritesByUser")
			} else {
				assert.Nil(t, err, "No error should have been returned from favoriteService.FindFavoritesByUser")
			}
			mockFavoriteStore.AssertExpectations(t)
		})
	}
}

func TestDefaultFavoriteService_DeleteFavorite(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		returnedError error
		expectError   bool
	}{
		{
			name:          "Successfully deleted the favorite",
			id:            "0",
			returnedError: nil,
			expectError:   false,
		},
		{
			name:          "Failed to delete the favorite",
			id:            "0",
			returnedError: fmt.Errorf("failed to delete the favorite"),
			expectError:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFavoriteStore := mocks.NewFavoriteStore(t)
			mockFavoriteStore.On("DeleteFavorite", tt.id).Return(tt.returnedError)
			favoriteService := NewDefaultFavoriteService(mockFavoriteStore)
			err := favoriteService.DeleteFavorite(tt.id)
			if tt.expectError {
				assert.NotNil(t, err, "An error should have been returned from favoriteService.DeleteFavorite")
			} else {
				assert.Nil(t, err, "No error should have been returned from favoriteService.DeleteFavorite")
			}
			mockFavoriteStore.AssertExpectations(t)
		})
	}
}
