package lambda

import (
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"the-drink-almanac-api/appErrors"
	"the-drink-almanac-api/dto"
	"the-drink-almanac-api/service"
)

type FavoritesLambdaHandler struct {
	favoriteService service.FavoriteService
	authService     service.AuthService
}

func NewFavoritesLambdaHandler(favoriteService service.FavoriteService, authService service.AuthService) FavoritesLambdaHandler {
	return FavoritesLambdaHandler{
		favoriteService: favoriteService,
		authService:     authService,
	}
}

func (h *FavoritesLambdaHandler) FindAllFavorites(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	favorites, err := h.favoriteService.FindAllFavorites()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	favoritesResponse := dto.NewFavoritesResponse(favorites)
	body, err := jsoniter.MarshalToString(favoritesResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       body,
	}
	return response, nil
}

func (h *FavoritesLambdaHandler) FindFavoritesByUser(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	userId, err := authorizeUser(request.Headers, h.authService)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	favorites, err := h.favoriteService.FindFavoritesByUser(userId)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if len(favorites) == 0 {
		response := events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       messageToResponseBody(fmt.Sprintf("no favorites were found for user with id %s", userId)),
		}
		return response, nil
	}

	favoritesResponse := dto.NewFavoritesResponse(favorites)
	body, err := jsoniter.MarshalToString(favoritesResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       body,
	}
	return response, nil
}

func (h *FavoritesLambdaHandler) CreateNewFavorite(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	userId, err := authorizeUser(request.Headers, h.authService)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	var newFavoritePostRequest dto.FavoritePostRequest
	if err := jsoniter.Unmarshal([]byte(request.Body), newFavoritePostRequest); err != nil {
		response := events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       messageToResponseBody(err.Error()),
		}
		return response, nil
	}

	if err := newFavoritePostRequest.ValidateRequest(); err != nil {
		response := events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       messageToResponseBody(err.Error()),
		}
		return response, nil
	}

	newFavorite, err := h.favoriteService.CreateNewFavorite(newFavoritePostRequest.DrinkId, userId)
	if err != nil {
		if errors.As(err, &appErrors.FavoriteAlreadyExistsError{}) {
			response := events.APIGatewayProxyResponse{
				StatusCode: http.StatusConflict,
				Body:       messageToResponseBody(fmt.Sprintf("the user '%s' already favorited the drink with id '%s'", newFavorite.UserId, newFavorite.DrinkId)),
			}
			return response, nil
		}
		response := events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       messageToResponseBody(fmt.Sprintf("unable to add the new favorite due to %s", err.Error())),
		}
		return response, nil
	}

	favoriteResponse := dto.NewFavoriteResponse(*newFavorite)
	body, err := jsoniter.MarshalToString(favoriteResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       body,
	}
	return response, nil
}

func (h *FavoritesLambdaHandler) DeleteFavorite(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	_, err := authorizeUser(request.Headers, h.authService)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	favoriteId := request.QueryStringParameters["favoriteId"]
	err = h.favoriteService.DeleteFavorite(favoriteId)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: http.StatusNoContent,
	}
	return response, nil
}

func (h *FavoritesLambdaHandler) RouteRequest(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	requestMarshalled, _ := jsoniter.MarshalToString(request)
	fmt.Printf("request: %v", requestMarshalled)
	switch request.RouteKey {
	case "GET /favorites":
		return h.FindAllFavorites(request)
	case "ANY /favorite/{drinkId}":
		switch request.RequestContext.HTTP.Method {
		case "GET":
			return h.FindFavoritesByUser(request)
		case "POST":
			return h.CreateNewFavorite(request)
		case "DELETE":
			return h.DeleteFavorite(request)
		default:
			fmt.Println("invalid method in request:", request)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       messageToResponseBody(fmt.Sprintf("invalid request method: '%s'", request.RouteKey)),
			}, nil
		}
	default:
		fmt.Printf("invalid path in request: %v", request)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       messageToResponseBody(fmt.Sprintf("invalid request path: '%s'", request.RawPath)),
		}, nil
	}
}

func messageToResponseBody(message string) string {
	m := map[string]string{
		"message": message,
	}
	body, _ := jsoniter.MarshalToString(m)
	return body
}
