package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"the-drink-almanac-api/appErrors"
	"the-drink-almanac-api/dto"
	"the-drink-almanac-api/handler"
	"the-drink-almanac-api/mocks"
	"the-drink-almanac-api/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestFindAllUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsers := []model.User{
		{
			Id:       "0",
			Username: "0",
			Password: "0",
		},
		{
			Id:       "1",
			Username: "1",
			Password: "0",
		},
		{
			Id:       "2",
			Username: "1",
			Password: "1",
		},
	}
	data := []struct {
		testName           string
		returnedUsers      []model.User
		returnedError      error
		expectedStatusCode int
	}{
		{
			testName:           "Successfully retrieve users",
			returnedUsers:      mockUsers,
			returnedError:      nil,
			expectedStatusCode: http.StatusOK,
		},
		{
			testName:           "Failed to retrieve users",
			returnedUsers:      nil,
			returnedError:      fmt.Errorf("failed to retrieve users"),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			mockUserService := mocks.NewUserService(t)
			mockUserService.On("FindAllUsers").Return(d.returnedUsers, d.returnedError)
			mockAuthService := mocks.NewAuthService(t)
			userHandler := handler.NewUserHandler(mockUserService, mockAuthService)

			rr := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, "/user", nil)
			assert.NoError(t, err)

			router := gin.Default()
			router.GET("/user", userHandler.FindAllUsers)
			router.ServeHTTP(rr, request)

			assert.Equal(t, d.expectedStatusCode, rr.Code)
			mockUserService.AssertExpectations(t)

			if d.returnedUsers != nil {
				usersResponse := dto.NewUsersResponse(d.returnedUsers)
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
			returnedError:        appErrors.NewUserAlreadyExistsError("user exists"),
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
			userHandler := handler.NewUserHandler(mockUserService, mockAuthService)

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
			testName:           "Not authorized",
			userId:             "0",
			returnedError:      nil,
			expectedStatusCode: http.StatusUnauthorized,
			authError:          fmt.Errorf("tsk tsk"),
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			mockUserService := mocks.NewUserService(t)
			mockAuthService := mocks.NewAuthService(t)
			mockAuthService.On("ValidateToken", "testToken").Return(d.userId, d.authError)
			if d.authError == nil {
				mockUserService.On("DeleteUser", d.userId).Return(d.returnedError)
			}
			userHandler := handler.NewUserHandler(mockUserService, mockAuthService)

			rr := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodDelete, "/user", nil)
			request.Header["Token"] = []string{"testToken"}
			assert.NoError(t, err)

			router := gin.Default()
			router.DELETE("/user", userHandler.DeleteUser)
			router.ServeHTTP(rr, request)

			assert.Equal(t, d.expectedStatusCode, rr.Code)
			mockUserService.AssertExpectations(t)
		})
	}
}
