package tests

import (
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
