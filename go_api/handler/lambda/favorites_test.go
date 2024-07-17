package lambda

import (
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"the-drink-almanac-api/apperrors"
	"the-drink-almanac-api/dto"
	"the-drink-almanac-api/model"
	"the-drink-almanac-api/service"
)

func TestNewFavoritesLambdaHandler(t *testing.T) {
	h := NewFavoritesLambdaHandler(nil, nil)
	assert.NotNil(t, h)
}

func TestFavoritesLambdaHandler_FindAllFavorites(t *testing.T) {
	favorites := []model.Favorite{
		{Id: "favorite1"},
		{Id: "favorite2"},
		{Id: "favorite3"},
	}
	dtoFavorites := []dto.FavoriteResponse{
		{Id: "favorite1"},
		{Id: "favorite2"},
		{Id: "favorite3"},
	}
	marshalledFavorites, err := jsoniter.MarshalToString(dtoFavorites)
	assert.NoError(t, err)

	testCases := map[string]struct {
		request        events.APIGatewayV2HTTPRequest
		mockCalls      func(ts *favoritesTestSuite)
		expectedResult events.APIGatewayV2HTTPResponse
		expectError    bool
	}{
		"Happy path": {
			request: events.APIGatewayV2HTTPRequest{
				Headers: map[string]string{
					"Token": "token",
				},
			},
			mockCalls: func(ts *favoritesTestSuite) {
				ts.mockAuthService.On("ValidateToken", "token").
					Return("userId", nil)

				ts.mockFavoriteService.On("FindAllFavorites").
					Return(favorites, nil)
			},
			expectedResult: events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusOK,
				Body:       marshalledFavorites,
			},
		},
		"Favorite service error": {
			request: events.APIGatewayV2HTTPRequest{
				Headers: map[string]string{
					"Token": "token",
				},
			},
			mockCalls: func(ts *favoritesTestSuite) {
				ts.mockAuthService.On("ValidateToken", "token").
					Return("userId", nil)

				ts.mockFavoriteService.On("FindAllFavorites").
					Return([]model.Favorite{}, errors.New("testing"))
			},
			expectedResult: events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       messageToResponseBody("testing"),
			},
		},
		"Auth service error": {
			request: events.APIGatewayV2HTTPRequest{
				Headers: map[string]string{
					"Token": "token",
				},
			},
			mockCalls: func(ts *favoritesTestSuite) {
				ts.mockAuthService.On("ValidateToken", "token").
					Return("", errors.New("testing"))
			},
			expectedResult: events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusForbidden,
				Body:       messageToResponseBody(InvalidTokenError.Error()),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ts := favoritesSetup(t)
			tc.mockCalls(ts)

			result, err := ts.handler.FindAllFavorites(tc.request)

			assert.Equal(t, tc.expectError, err != nil)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestFavoritesLambdaHandler_FindFavoritesByUser(t *testing.T) {
	favorites := []model.Favorite{
		{Id: "favorite1"},
		{Id: "favorite2"},
		{Id: "favorite3"},
	}
	dtoFavorites := []dto.FavoriteResponse{
		{Id: "favorite1"},
		{Id: "favorite2"},
		{Id: "favorite3"},
	}
	marshalledFavorites, err := jsoniter.MarshalToString(dtoFavorites)
	assert.NoError(t, err)

	testCases := map[string]struct {
		request        events.APIGatewayV2HTTPRequest
		mockCalls      func(ts *favoritesTestSuite)
		expectedResult events.APIGatewayV2HTTPResponse
		expectError    bool
	}{
		"Happy path": {
			request: events.APIGatewayV2HTTPRequest{
				Headers: map[string]string{
					"Token": "token",
				},
			},
			mockCalls: func(ts *favoritesTestSuite) {
				ts.mockAuthService.On("ValidateToken", "token").
					Return("userId", nil)

				ts.mockFavoriteService.On("FindFavoritesByUser", "userId").
					Return(favorites, nil)
			},
			expectedResult: events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusOK,
				Body:       marshalledFavorites,
			},
		},
		"No favorites found": {
			request: events.APIGatewayV2HTTPRequest{
				Headers: map[string]string{
					"Token": "token",
				},
			},
			mockCalls: func(ts *favoritesTestSuite) {
				ts.mockAuthService.On("ValidateToken", "token").
					Return("userId", nil)

				ts.mockFavoriteService.On("FindFavoritesByUser", "userId").
					Return([]model.Favorite{}, nil)
			},
			expectedResult: events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusNotFound,
				Body:       messageToResponseBody("no favorites were found for user with id userId"),
			},
		},
		"Favorite service error": {
			request: events.APIGatewayV2HTTPRequest{
				Headers: map[string]string{
					"Token": "token",
				},
			},
			mockCalls: func(ts *favoritesTestSuite) {
				ts.mockAuthService.On("ValidateToken", "token").
					Return("userId", nil)

				ts.mockFavoriteService.On("FindFavoritesByUser", "userId").
					Return([]model.Favorite{}, errors.New("testing"))
			},
			expectedResult: events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       messageToResponseBody("testing"),
			},
		},
		"Auth service error": {
			request: events.APIGatewayV2HTTPRequest{
				Headers: map[string]string{
					"Token": "token",
				},
			},
			mockCalls: func(ts *favoritesTestSuite) {
				ts.mockAuthService.On("ValidateToken", "token").
					Return("", errors.New("testing"))
			},
			expectedResult: events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusForbidden,
				Body:       messageToResponseBody(InvalidTokenError.Error()),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ts := favoritesSetup(t)
			tc.mockCalls(ts)

			result, err := ts.handler.FindFavoritesByUser(tc.request)

			assert.Equal(t, tc.expectError, err != nil)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestFavoritesLambdaHandler_CreateNewFavorite(t *testing.T) {
	favoriteRequest := dto.FavoritePostRequest{DrinkId: "drink1"}
	marshalledFavoriteRequest, err := jsoniter.MarshalToString(favoriteRequest)
	assert.NoError(t, err)

	favorite := model.Favorite{DrinkId: "drink1", UserId: "userId"}
	favoriteResponse := dto.FavoriteResponse{DrinkId: "drink1"}
	marshalledFavoriteResponse, err := jsoniter.MarshalToString(favoriteResponse)
	assert.NoError(t, err)

	testCases := map[string]struct {
		request        events.APIGatewayV2HTTPRequest
		mockCalls      func(ts *favoritesTestSuite)
		expectedResult events.APIGatewayV2HTTPResponse
		expectError    bool
	}{
		"Happy path": {
			request: events.APIGatewayV2HTTPRequest{
				Headers: map[string]string{
					"Token": "token",
				},
				Body: marshalledFavoriteRequest,
			},
			mockCalls: func(ts *favoritesTestSuite) {
				ts.mockAuthService.On("ValidateToken", "token").
					Return("userId", nil)

				ts.mockFavoriteService.On("CreateNewFavorite", "drink1", "userId").
					Return(&favorite, nil)
			},
			expectedResult: events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusOK,
				Body:       marshalledFavoriteResponse,
			},
		},
		"Favorite already exists": {
			request: events.APIGatewayV2HTTPRequest{
				Headers: map[string]string{
					"Token": "token",
				},
				Body: marshalledFavoriteRequest,
			},
			mockCalls: func(ts *favoritesTestSuite) {
				ts.mockAuthService.On("ValidateToken", "token").
					Return("userId", nil)

				ts.mockFavoriteService.On("CreateNewFavorite", "drink1", "userId").
					Return(&model.Favorite{}, apperrors.FavoriteAlreadyExistsError{})
			},
			expectedResult: events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusConflict,
				Body:       messageToResponseBody("the user 'userId' already favorited the drink with id 'drink1'"),
			},
		},
		"Favorite service error": {
			request: events.APIGatewayV2HTTPRequest{
				Headers: map[string]string{
					"Token": "token",
				},
				Body: marshalledFavoriteRequest,
			},
			mockCalls: func(ts *favoritesTestSuite) {
				ts.mockAuthService.On("ValidateToken", "token").
					Return("userId", nil)

				ts.mockFavoriteService.On("CreateNewFavorite", "drink1", "userId").
					Return(&model.Favorite{}, errors.New("testing"))
			},
			expectedResult: events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       messageToResponseBody("unable to add the new favorite due to testing"),
			},
		},
		"Auth service error": {
			request: events.APIGatewayV2HTTPRequest{
				Headers: map[string]string{
					"Token": "token",
				},
				Body: marshalledFavoriteRequest,
			},
			mockCalls: func(ts *favoritesTestSuite) {
				ts.mockAuthService.On("ValidateToken", "token").
					Return("", errors.New("testing"))
			},
			expectedResult: events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusForbidden,
				Body:       messageToResponseBody(InvalidTokenError.Error()),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ts := favoritesSetup(t)
			tc.mockCalls(ts)

			result, err := ts.handler.CreateNewFavorite(tc.request)

			assert.Equal(t, tc.expectError, err != nil)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestFavoritesLambdaHandler_DeleteFavorite(t *testing.T) {
	testCases := map[string]struct {
		request        events.APIGatewayV2HTTPRequest
		mockCalls      func(ts *favoritesTestSuite)
		expectedResult events.APIGatewayV2HTTPResponse
		expectError    bool
	}{
		"Happy path": {
			request: events.APIGatewayV2HTTPRequest{
				Headers: map[string]string{
					"Token": "token",
				},
				QueryStringParameters: map[string]string{
					"favoriteId": "favorite1",
				},
			},
			mockCalls: func(ts *favoritesTestSuite) {
				ts.mockAuthService.On("ValidateToken", "token").
					Return("userId", nil)

				ts.mockFavoriteService.On("DeleteFavorite", "favorite1").
					Return(nil)
			},
			expectedResult: events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusNoContent,
			},
		},
		"Favorite service error": {
			request: events.APIGatewayV2HTTPRequest{
				Headers: map[string]string{
					"Token": "token",
				},
				QueryStringParameters: map[string]string{
					"favoriteId": "favorite1",
				},
			},
			mockCalls: func(ts *favoritesTestSuite) {
				ts.mockAuthService.On("ValidateToken", "token").
					Return("userId", nil)

				ts.mockFavoriteService.On("DeleteFavorite", "favorite1").
					Return(errors.New("testing"))
			},
			expectedResult: events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       messageToResponseBody("testing"),
			},
		},
		"Auth service error": {
			request: events.APIGatewayV2HTTPRequest{
				Headers: map[string]string{
					"Token": "token",
				},
				QueryStringParameters: map[string]string{
					"favoriteId": "favorite1",
				},
			},
			mockCalls: func(ts *favoritesTestSuite) {
				ts.mockAuthService.On("ValidateToken", "token").
					Return("", errors.New("testing"))
			},
			expectedResult: events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusForbidden,
				Body:       messageToResponseBody(InvalidTokenError.Error()),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ts := favoritesSetup(t)
			tc.mockCalls(ts)

			result, err := ts.handler.DeleteFavorite(tc.request)

			assert.Equal(t, tc.expectError, err != nil)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

type favoritesTestSuite struct {
	mockFavoriteService *service.MockFavoriteService
	mockAuthService     *service.MockAuthService
	handler             *FavoritesLambdaHandler
}

func favoritesSetup(t *testing.T) *favoritesTestSuite {
	mockFavoriteService := service.NewMockFavoriteService(t)
	mockAuthService := service.NewMockAuthService(t)
	handler := &FavoritesLambdaHandler{
		favoriteService: mockFavoriteService,
		authService:     mockAuthService,
	}

	return &favoritesTestSuite{
		mockFavoriteService: mockFavoriteService,
		mockAuthService:     mockAuthService,
		handler:             handler,
	}
}
