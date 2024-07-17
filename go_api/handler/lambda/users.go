package lambda

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	jsoniter "github.com/json-iterator/go"
	"the-drink-almanac-api/apperrors"
	"the-drink-almanac-api/dto"
	"the-drink-almanac-api/service"
)

type UsersLambdaHandler struct {
	userService service.UserService
	authService service.AuthService
}

func NewUsersLambdaHandler(userService service.UserService, authService service.AuthService) UsersLambdaHandler {
	return UsersLambdaHandler{
		userService: userService,
		authService: authService,
	}
}

func (h *UsersLambdaHandler) FindUser(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	userId, err := authorizeUser(request.Headers, h.authService)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusForbidden,
			Body:       messageToResponseBody(err.Error()),
		}, nil
	}

	user, err := h.userService.FindUser(userId)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       messageToResponseBody(err.Error()),
		}, nil
	}

	if user == nil {
		response := events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusNotFound,
			Body:       messageToResponseBody(fmt.Sprintf("no user was found for user id %s", userId)),
		}
		return response, nil
	}

	userResponse := dto.NewUserResponse(*user)
	body, err := jsoniter.MarshalToString(userResponse)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       messageToResponseBody(err.Error()),
		}, nil
	}

	response := events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       body,
	}
	return response, nil
}

func (h *UsersLambdaHandler) CreateNewUser(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var userRequest dto.UserPostRequest
	if err := jsoniter.Unmarshal([]byte(request.Body), userRequest); err != nil {
		response := events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       messageToResponseBody(err.Error()),
		}
		return response, nil
	}

	if err := userRequest.ValidateRequest(); err != nil {
		response := events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       messageToResponseBody(err.Error()),
		}
		return response, nil
	}

	user, err := h.userService.CreateNewUser(userRequest.Username, userRequest.Password)
	if err != nil {
		if errors.As(err, &apperrors.UserAlreadyExistsError{}) {
			response := events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusConflict,
				Body:       messageToResponseBody(fmt.Sprintf("a user already exists with the username %s", user.Username)),
			}
			return response, nil
		}
		response := events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       messageToResponseBody(err.Error()),
		}
		return response, nil
	}
	response := events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusCreated,
	}
	return response, nil
}

func (h *UsersLambdaHandler) DeleteUser(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	userId, err := authorizeUser(request.Headers, h.authService)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusForbidden,
			Body:       messageToResponseBody(err.Error()),
		}, nil
	}

	err = h.userService.DeleteUser(userId)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.As(err, &apperrors.InvalidAuthTokenError{}) {
			statusCode = http.StatusForbidden
		}
		return events.APIGatewayV2HTTPResponse{
			StatusCode: statusCode,
			Body:       messageToResponseBody(err.Error()),
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusNoContent,
		Body:       messageToResponseBody("the user was deleted"),
	}, nil
}

func (h *UsersLambdaHandler) Login(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var userRequest dto.UserPostRequest
	if err := jsoniter.Unmarshal([]byte(request.Body), userRequest); err != nil {
		response := events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       messageToResponseBody(err.Error()),
		}
		return response, nil
	}

	if err := userRequest.ValidateRequest(); err != nil {
		response := events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       messageToResponseBody(err.Error()),
		}
		return response, nil
	}

	user, err := h.userService.Login(userRequest.Username, userRequest.Password)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.As(err, &apperrors.UserNotFoundError{}) {
			statusCode = http.StatusNotFound
		}
		if errors.As(err, &apperrors.IncorrectPasswordError{}) {
			statusCode = http.StatusBadRequest
		}

		return events.APIGatewayV2HTTPResponse{
			StatusCode: statusCode,
			Body:       messageToResponseBody(err.Error()),
		}, nil
	}

	tokenString, err := h.authService.CreateNewToken(user.Id, 60*24)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       messageToResponseBody(err.Error()),
		}, nil
	}

	response := events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Token": tokenString,
		},
	}
	return response, nil
}

func (h *UsersLambdaHandler) RouteRequest(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	requestMarshalled, _ := jsoniter.MarshalToString(request)
	fmt.Printf("request: %v", requestMarshalled)
	switch request.RouteKey {
	case "GET /user":
		return h.FindUser(request)
	case "DELETE /user":
		return h.DeleteUser(request)
	case "POST /user/login":
		return h.Login(request)
	case "POST /user/register":
		return h.CreateNewUser(request)
	default:
		fmt.Printf("invalid path in request: %v", request)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       messageToResponseBody(fmt.Sprintf("invalid request path: '%s'", request.RawPath)),
		}, nil
	}
}
