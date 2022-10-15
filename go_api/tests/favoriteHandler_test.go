package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"the-drink-almanac-api/handler"
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
			favoriteHandler := handler.FavoriteHandlers{mockFavoriteService}

			rr := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, "/favorite", nil)
			assert.NoError(t, err)

			router := gin.Default()
			router.GET("/favorite", favoriteHandler.FindAllFavorites)
			router.ServeHTTP(rr, request)

			assert.Equal(t, d.expectedStatusCode, rr.Code)
			mockFavoriteService.AssertExpectations(t)

			if d.returnedFavorites != nil {
				expectedResponseBody, err := json.Marshal(d.returnedFavorites)
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
			favoriteHandler := handler.FavoriteHandlers{mockFavoriteService}

			rr := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/favorite/%s", d.userId), nil)
			assert.NoError(t, err)

			router := gin.Default()
			router.GET("/favorite/:userId", favoriteHandler.FindFavoritesByUser)
			router.ServeHTTP(rr, request)

			assert.Equal(t, d.expectedStatusCode, rr.Code)
			mockFavoriteService.AssertExpectations(t)

			if d.returnedFavorites != nil && len(d.returnedFavorites) != 0 {
				expectedResponseBody, err := json.Marshal(d.returnedFavorites)
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
		testName           string
		userId             string
		drinkId            string
		requestBody        []byte
		returnedFavorite   *model.Favorite
		returnedError      error
		expectedStatusCode int
	}{
		{
			testName:           "Successfully create favorite",
			userId:             "0",
			drinkId:            "0",
			requestBody:        []byte(`{"user_id": "0", "drink_id": "0"}`),
			returnedFavorite:   mockFavorite,
			returnedError:      nil,
			expectedStatusCode: http.StatusCreated,
		},
		{
			testName:           "Failed to create favorite",
			userId:             "0",
			drinkId:            "0",
			requestBody:        []byte(`{"user_id": "0", "drink_id": "0"}`),
			returnedFavorite:   nil,
			returnedError:      fmt.Errorf("failed to create favorite"),
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			testName:           "No request body",
			userId:             "",
			drinkId:            "",
			requestBody:        nil,
			returnedFavorite:   nil,
			returnedError:      fmt.Errorf("failed to create favorite"),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testName:           "User id not provided",
			userId:             "",
			drinkId:            "0",
			requestBody:        []byte(`{"drink_id": "0"}`),
			returnedFavorite:   nil,
			returnedError:      fmt.Errorf("user id not provided"),
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			testName:           "Drink id not provided",
			userId:             "0",
			drinkId:            "",
			requestBody:        []byte(`{"user_id": "0"}`),
			returnedFavorite:   nil,
			returnedError:      fmt.Errorf("drink id not provided"),
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			mockFavoriteService := mocks.NewFavoriteService(t)
			if d.userId != "" && d.drinkId != "" {
				mockFavoriteService.On("CreateNewFavorite", d.userId, d.drinkId).Return(d.returnedFavorite, d.returnedError)
			}
			favoriteHandler := handler.FavoriteHandlers{mockFavoriteService}

			rr := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodPost, "/favorite", bytes.NewBuffer(d.requestBody))
			assert.NoError(t, err)

			router := gin.Default()
			router.POST("/favorite", favoriteHandler.CreateNewFavorite)
			router.ServeHTTP(rr, request)

			assert.Equal(t, d.expectedStatusCode, rr.Code)
			mockFavoriteService.AssertExpectations(t)

			if d.returnedFavorite != nil {
				expectedResponseBody, err := json.Marshal(d.returnedFavorite)
				assert.NoError(t, err)
				assert.Equal(t, expectedResponseBody, rr.Body.Bytes())
			}
		})
	}
}
