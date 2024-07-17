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
	"github.com/stretchr/testify/mock"
)

func TestFindUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	data := []struct {
		testName           string
		userId             string
		returnedUser       *model.User
		returnedError      error
		expectedStatusCode int
	}{
		{
			testName: "Successfully retrieve user",
			userId:   "0",
			returnedUser: &model.User{
				Id:       "0",
				Username: "0",
				Password: "0",
			},
			returnedError:      nil,
			expectedStatusCode: http.StatusOK,
		},
		{
			testName:           "Failed to retrieve user",
			userId:             "0",
			returnedUser:       nil,
			returnedError:      nil,
			expectedStatusCode: http.StatusNotFound,
		},
		{
			testName:           "No user id in context",
			userId:             "",
			returnedUser:       nil,
			returnedError:      nil,
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			mockUserService := mocks.NewUserService(t)
			if d.userId != "" {
				mockUserService.On("FindUser", d.userId).Return(d.returnedUser, d.returnedError)
			}
			mockAuthService := mocks.NewAuthService(t)
			userHandler := NewUserHandler(mockUserService, mockAuthService)

			rr := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, "/user", nil)
			assert.NoError(t, err)

			router := gin.Default()
			router.GET("/user", setUserIdInContext(d.userId), userHandler.FindUser)
			router.ServeHTTP(rr, request)

			assert.Equal(t, d.expectedStatusCode, rr.Code)
			mockUserService.AssertExpectations(t)

			if d.returnedUser != nil {
				usersResponse := dto.NewUserResponse(*d.returnedUser)
				expectedResponseBody, err := json.Marshal(usersResponse)
				assert.NoError(t, err)
				assert.Equal(t, expectedResponseBody, rr.Body.Bytes())
			}
		})
	}
}

func TestCreateNewUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUser := &model.User{
		Id:       "0",
		Username: "0",
		Password: "0",
	}
	data := []struct {
		testName             string
		username             string
		password             string
		requestBody          []byte
		returnedUser         *model.User
		returnedError        error
		expectedStatusCode   int
		shouldMethodBeCalled bool
	}{
		{
			testName:             "Successfully create user",
			username:             "0",
			password:             "0",
			requestBody:          []byte(`{"username": "0", "password": "0"}`),
			returnedUser:         mockUser,
			returnedError:        nil,
			expectedStatusCode:   http.StatusCreated,
			shouldMethodBeCalled: true,
		},
		{
			testName:             "Failed to create user",
			username:             "0",
			password:             "0",
			requestBody:          []byte(`{"username": "0", "password": "0"}`),
			returnedUser:         nil,
			returnedError:        fmt.Errorf("failed to create user"),
			expectedStatusCode:   http.StatusInternalServerError,
			shouldMethodBeCalled: true,
		},
		{
			testName:             "User already exists",
			username:             "0",
			password:             "0",
			requestBody:          []byte(`{"username": "0", "password": "0"}`),
			returnedUser:         mockUser,
			returnedError:        apperrors.NewUserAlreadyExistsError("user exists"),
			expectedStatusCode:   http.StatusConflict,
			shouldMethodBeCalled: true,
		},
		{
			testName:             "No request body",
			username:             "",
			password:             "",
			requestBody:          nil,
			returnedUser:         nil,
			returnedError:        fmt.Errorf("failed to create user"),
			expectedStatusCode:   http.StatusBadRequest,
			shouldMethodBeCalled: false,
		},
		{
			testName:             "username not provided",
			username:             "",
			password:             "0",
			requestBody:          []byte(`{"password": "0"}`),
			returnedUser:         nil,
			returnedError:        fmt.Errorf("username not provided"),
			expectedStatusCode:   http.StatusBadRequest,
			shouldMethodBeCalled: false,
		},
		{
			testName:             "password not provided",
			username:             "0",
			password:             "",
			requestBody:          []byte(`{"username": "0"}`),
			returnedUser:         nil,
			returnedError:        fmt.Errorf("password not provided"),
			expectedStatusCode:   http.StatusBadRequest,
			shouldMethodBeCalled: false,
		},
		{
			testName:             "Non-string username provided",
			username:             "0",
			password:             "0",
			requestBody:          []byte(`{"username": 0, "password": "0"}`),
			returnedUser:         nil,
			returnedError:        fmt.Errorf("username is not a string"),
			expectedStatusCode:   http.StatusBadRequest,
			shouldMethodBeCalled: false,
		},
		{
			testName:             "Non-string password provided",
			username:             "0",
			password:             "0",
			requestBody:          []byte(`{"username": "0", "password": 0}`),
			returnedUser:         nil,
			returnedError:        fmt.Errorf("password is not a string"),
			expectedStatusCode:   http.StatusBadRequest,
			shouldMethodBeCalled: false,
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			mockUserService := mocks.NewUserService(t)
			if d.shouldMethodBeCalled {
				mockUserService.On("CreateNewUser", d.username, d.password).Return(d.returnedUser, d.returnedError)
			}
			mockAuthService := mocks.NewAuthService(t)
			userHandler := NewUserHandler(mockUserService, mockAuthService)

			rr := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(d.requestBody))
			assert.NoError(t, err)

			router := gin.Default()
			router.POST("/user", userHandler.CreateNewUser)
			router.ServeHTTP(rr, request)

			assert.Equal(t, d.expectedStatusCode, rr.Code)
			mockUserService.AssertExpectations(t)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	data := []struct {
		testName           string
		userId             string
		returnedError      error
		expectedStatusCode int
		authError          error
	}{
		{
			testName:           "Successfully delete user",
			userId:             "0",
			returnedError:      nil,
			expectedStatusCode: http.StatusNoContent,
			authError:          nil,
		},
		{
			testName:           "Failed to create user",
			userId:             "0",
			returnedError:      fmt.Errorf("failed to delete user"),
			expectedStatusCode: http.StatusInternalServerError,
			authError:          nil,
		},
		{
			testName:           "No user exists",
			userId:             "0",
			returnedError:      nil,
			expectedStatusCode: http.StatusNoContent,
			authError:          nil,
		},
		{
			testName:           "User id not retrieved",
			userId:             "",
			returnedError:      nil,
			expectedStatusCode: http.StatusUnauthorized,
			authError:          nil,
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			mockUserService := mocks.NewUserService(t)
			mockAuthService := mocks.NewAuthService(t)
			if d.userId != "" {
				mockUserService.On("DeleteUser", d.userId).Return(d.returnedError)
			}
			userHandler := NewUserHandler(mockUserService, mockAuthService)

			rr := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodDelete, "/user", nil)
			assert.NoError(t, err)

			router := gin.Default()
			router.DELETE("/user", setUserIdInContext(d.userId), userHandler.DeleteUser)
			router.ServeHTTP(rr, request)

			assert.Equal(t, d.expectedStatusCode, rr.Code)
			mockUserService.AssertExpectations(t)
		})
	}
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUser := &model.User{
		Id:       "0",
		Username: "0",
		Password: "0",
	}
	data := []struct {
		testName             string
		username             string
		password             string
		requestBody          []byte
		returnedUser         *model.User
		returnedError        error
		expectedStatusCode   int
		shouldMethodBeCalled bool
		shouldReturnToken    bool
	}{
		{
			testName:             "Successfully logged in",
			username:             "0",
			password:             "0",
			requestBody:          []byte(`{"username": "0", "password": "0"}`),
			returnedUser:         mockUser,
			returnedError:        nil,
			expectedStatusCode:   http.StatusOK,
			shouldMethodBeCalled: true,
			shouldReturnToken:    true,
		},
		{
			testName:             "Failed to log in",
			username:             "0",
			password:             "0",
			requestBody:          []byte(`{"username": "0", "password": "0"}`),
			returnedUser:         nil,
			returnedError:        fmt.Errorf("failed to log in"),
			expectedStatusCode:   http.StatusInternalServerError,
			shouldMethodBeCalled: true,
			shouldReturnToken:    false,
		},
		{
			testName:             "No request body",
			username:             "",
			password:             "",
			requestBody:          nil,
			returnedUser:         nil,
			returnedError:        fmt.Errorf("failed to log in"),
			expectedStatusCode:   http.StatusBadRequest,
			shouldMethodBeCalled: false,
			shouldReturnToken:    false,
		},
		{
			testName:             "username not provided",
			username:             "",
			password:             "0",
			requestBody:          []byte(`{"password": "0"}`),
			returnedUser:         nil,
			returnedError:        fmt.Errorf("username not provided"),
			expectedStatusCode:   http.StatusBadRequest,
			shouldMethodBeCalled: false,
			shouldReturnToken:    false,
		},
		{
			testName:             "password not provided",
			username:             "0",
			password:             "",
			requestBody:          []byte(`{"username": "0"}`),
			returnedUser:         nil,
			returnedError:        fmt.Errorf("password not provided"),
			expectedStatusCode:   http.StatusBadRequest,
			shouldMethodBeCalled: false,
			shouldReturnToken:    false,
		},
		{
			testName:             "Non-string username provided",
			username:             "0",
			password:             "0",
			requestBody:          []byte(`{"username": 0, "password": "0"}`),
			returnedUser:         nil,
			returnedError:        fmt.Errorf("username is not a string"),
			expectedStatusCode:   http.StatusBadRequest,
			shouldMethodBeCalled: false,
			shouldReturnToken:    false,
		},
		{
			testName:             "Non-string password provided",
			username:             "0",
			password:             "0",
			requestBody:          []byte(`{"username": "0", "password": 0}`),
			returnedUser:         nil,
			returnedError:        fmt.Errorf("password is not a string"),
			expectedStatusCode:   http.StatusBadRequest,
			shouldMethodBeCalled: false,
			shouldReturnToken:    false,
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			mockUserService := mocks.NewUserService(t)
			if d.shouldMethodBeCalled {
				mockUserService.On("Login", d.username, d.password).Return(d.returnedUser, d.returnedError)
			}
			mockAuthService := mocks.NewAuthService(t)
			if d.shouldReturnToken {
				mockAuthService.On("CreateNewToken", d.returnedUser.Id, mock.Anything).Return("testToken", nil)
			}
			userHandler := NewUserHandler(mockUserService, mockAuthService)

			rr := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodPost, "/user/login", bytes.NewBuffer(d.requestBody))
			assert.NoError(t, err)

			router := gin.Default()
			router.POST("/user/login", userHandler.Login)
			router.ServeHTTP(rr, request)

			assert.Equal(t, d.expectedStatusCode, rr.Code)
			_, ok := rr.HeaderMap["Token"]
			assert.Equal(t, d.shouldReturnToken, ok)
			mockUserService.AssertExpectations(t)
		})
	}
}
