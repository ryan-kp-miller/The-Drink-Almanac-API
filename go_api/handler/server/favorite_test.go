package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"the-drink-almanac-api/apperrors"
	"the-drink-almanac-api/dto"
	"the-drink-almanac-api/mocks"
	"the-drink-almanac-api/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestFindAllFavorites(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockFavorites := []model.Favorite{
		{
			Id:      "0",
			DrinkId: "0",
			UserId:  "0",
		},
		{
			Id:      "1",
			DrinkId: "1",
			UserId:  "0",
		},
		{
			Id:      "2",
			DrinkId: "1",
			UserId:  "1",
		},
	}
	data := []struct {
		testName           string
		returnedFavorites  []model.Favorite
		returnedError      error
		expectedStatusCode int
	}{
		{
			testName:           "Successfully retrieve favorites",
			returnedFavorites:  mockFavorites,
			returnedError:      nil,
			expectedStatusCode: http.StatusOK,
		},
		{
			testName:           "Failed to retrieve favorites",
			returnedFavorites:  nil,
			returnedError:      fmt.Errorf("failed to retrieve favorites"),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			mockFavoriteService := mocks.NewFavoriteService(t)
			mockFavoriteService.On("FindAllFavorites").Return(d.returnedFavorites, d.returnedError)
			favoriteHandler := FavoriteHandler{Service: mockFavoriteService}

			rr := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, "/favorite", nil)
			assert.NoError(t, err)

			router := gin.Default()
			router.GET("/favorite", favoriteHandler.FindAllFavorites)
			router.ServeHTTP(rr, request)

			assert.Equal(t, d.expectedStatusCode, rr.Code)
			mockFavoriteService.AssertExpectations(t)

			if d.returnedFavorites != nil {
				favoritesResponse := dto.NewFavoritesResponse(d.returnedFavorites)
				expectedResponseBody, err := json.Marshal(favoritesResponse)
				assert.NoError(t, err)
				assert.Equal(t, expectedResponseBody, rr.Body.Bytes())
			}
		})
	}
}

func TestFindFavoritesByUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockFavorites := []model.Favorite{
		{
			Id:      "0",
			DrinkId: "0",
			UserId:  "0",
		},
		{
			Id:      "1",
			DrinkId: "1",
			UserId:  "0",
		},
	}
	data := []struct {
		testName           string
		userId             string
		returnedFavorites  []model.Favorite
		returnedError      error
		expectedStatusCode int
	}{
		{
			testName:           "Successfully retrieve favorites",
			userId:             "0",
			returnedFavorites:  mockFavorites,
			returnedError:      nil,
			expectedStatusCode: http.StatusOK,
		},
		{
			testName:           "Failed to retrieve favorites",
			userId:             "0",
			returnedFavorites:  nil,
			returnedError:      fmt.Errorf("failed to retrieve favorites"),
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			testName:           "User doesn't have any favorites",
			userId:             "0",
			returnedFavorites:  make([]model.Favorite, 0),
			returnedError:      nil,
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			mockFavoriteService := mocks.NewFavoriteService(t)
			mockFavoriteService.On("FindFavoritesByUser", d.userId).Return(d.returnedFavorites, d.returnedError)
			favoriteHandler := FavoriteHandler{Service: mockFavoriteService}

			rr := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/favorite"), nil)
			assert.NoError(t, err)

			gin.SetMode(gin.TestMode)
			router := gin.Default()

			router.GET("/favorite", setUserIdInContext(d.userId), favoriteHandler.FindFavoritesByUser)
			router.ServeHTTP(rr, request)

			assert.Equal(t, d.expectedStatusCode, rr.Code)
			mockFavoriteService.AssertExpectations(t)

			if d.returnedFavorites != nil && len(d.returnedFavorites) != 0 {
				favoritesResponse := dto.NewFavoritesResponse(d.returnedFavorites)
				expectedResponseBody, err := json.Marshal(favoritesResponse)
				assert.NoError(t, err)
				assert.Equal(t, expectedResponseBody, rr.Body.Bytes())
			}
		})
	}
}

func TestCreateNewFavorite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockFavorite := &model.Favorite{
		Id:      "0",
		DrinkId: "0",
		UserId:  "0",
	}
	data := []struct {
		testName             string
		userId               string
		drinkId              string
		requestBody          []byte
		returnedFavorite     *model.Favorite
		returnedError        error
		expectedStatusCode   int
		shouldMethodBeCalled bool
	}{
		{
			testName:             "Successfully create favorite",
			userId:               "0",
			drinkId:              "0",
			requestBody:          []byte(`{"user_id": "0", "drink_id": "0"}`),
			returnedFavorite:     mockFavorite,
			returnedError:        nil,
			expectedStatusCode:   http.StatusCreated,
			shouldMethodBeCalled: true,
		},
		{
			testName:             "Failed to create favorite",
			userId:               "0",
			drinkId:              "0",
			requestBody:          []byte(`{"user_id": "0", "drink_id": "0"}`),
			returnedFavorite:     nil,
			returnedError:        fmt.Errorf("failed to create favorite"),
			expectedStatusCode:   http.StatusInternalServerError,
			shouldMethodBeCalled: true,
		},
		{
			testName:             "Favorite already exists",
			userId:               "0",
			drinkId:              "0",
			requestBody:          []byte(`{"user_id": "0", "drink_id": "0"}`),
			returnedFavorite:     mockFavorite,
			returnedError:        apperrors.NewFavoriteAlreadyExistsError("favorite already exists"),
			expectedStatusCode:   http.StatusConflict,
			shouldMethodBeCalled: true,
		},
		{
			testName:             "No request body",
			userId:               "",
			drinkId:              "",
			requestBody:          nil,
			returnedFavorite:     nil,
			returnedError:        fmt.Errorf("failed to create favorite"),
			expectedStatusCode:   http.StatusBadRequest,
			shouldMethodBeCalled: false,
		},
		{
			testName:             "Drink id not provided",
			userId:               "0",
			drinkId:              "",
			requestBody:          []byte(`{"user_id": "0"}`),
			returnedFavorite:     nil,
			returnedError:        fmt.Errorf("drink id not provided"),
			expectedStatusCode:   http.StatusBadRequest,
			shouldMethodBeCalled: false,
		},
		{
			testName:             "Non-string drink id provided",
			userId:               "0",
			drinkId:              "0",
			requestBody:          []byte(`{"user_id": "0", "drink_id": 0}`),
			returnedFavorite:     nil,
			returnedError:        fmt.Errorf("drink id is not a string"),
			expectedStatusCode:   http.StatusBadRequest,
			shouldMethodBeCalled: false,
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			mockFavoriteService := mocks.NewFavoriteService(t)
			if d.shouldMethodBeCalled {
				mockFavoriteService.On("CreateNewFavorite", d.drinkId, d.userId).Return(d.returnedFavorite, d.returnedError)
			}
			favoriteHandler := FavoriteHandler{Service: mockFavoriteService}

			rr := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodPost, "/favorite", bytes.NewBuffer(d.requestBody))
			assert.NoError(t, err)

			router := gin.Default()
			router.POST("/favorite", setUserIdInContext(d.userId), favoriteHandler.CreateNewFavorite)
			router.ServeHTTP(rr, request)

			assert.Equal(t, d.expectedStatusCode, rr.Code)
			mockFavoriteService.AssertExpectations(t)

			if d.returnedFavorite != nil && d.returnedError == nil {
				favoriteResponse := dto.NewFavoriteResponse(*d.returnedFavorite)
				expectedResponseBody, err := json.Marshal(favoriteResponse)
				assert.NoError(t, err)
				assert.Equal(t, expectedResponseBody, rr.Body.Bytes())
			}
		})
	}
}

func TestDeleteFavorite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	data := []struct {
		testName           string
		favoriteId         string
		returnedError      error
		expectedStatusCode int
	}{
		{
			testName:           "Successfully delete favorite",
			favoriteId:         "0",
			returnedError:      nil,
			expectedStatusCode: http.StatusNoContent,
		},
		{
			testName:           "Failed to delete favorite",
			favoriteId:         "0",
			returnedError:      fmt.Errorf("failed to delete favorite"),
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			testName:           "No favorite exists",
			favoriteId:         "0",
			returnedError:      nil,
			expectedStatusCode: http.StatusNoContent,
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			mockFavoriteService := mocks.NewFavoriteService(t)
			mockFavoriteService.On("DeleteFavorite", d.favoriteId).Return(d.returnedError)
			favoriteHandler := FavoriteHandler{Service: mockFavoriteService}

			rr := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/favorite/%s", d.favoriteId), nil)
			assert.NoError(t, err)

			router := gin.Default()
			router.DELETE("/favorite/:favoriteId", favoriteHandler.DeleteFavorite)
			router.ServeHTTP(rr, request)

			assert.Equal(t, d.expectedStatusCode, rr.Code)
			mockFavoriteService.AssertExpectations(t)
		})
	}
}
