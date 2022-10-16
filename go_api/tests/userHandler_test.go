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
			userHandler := handler.UserHandlers{Service: mockUserService}

			rr := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, "/user", nil)
			assert.NoError(t, err)

			router := gin.Default()
			router.GET("/user", userHandler.FindAllUsers)
			router.ServeHTTP(rr, request)

			assert.Equal(t, d.expectedStatusCode, rr.Code)
			mockUserService.AssertExpectations(t)

			if d.returnedUsers != nil {
				expectedResponseBody, err := json.Marshal(d.returnedUsers)
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
			userHandler := handler.UserHandlers{Service: mockUserService}

			rr := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(d.requestBody))
			assert.NoError(t, err)

			router := gin.Default()
			router.POST("/user", userHandler.CreateNewUser)
			router.ServeHTTP(rr, request)

			assert.Equal(t, d.expectedStatusCode, rr.Code)
			mockUserService.AssertExpectations(t)

			if d.returnedUser != nil {
				expectedResponseBody, err := json.Marshal(d.returnedUser)
				assert.NoError(t, err)
				assert.Equal(t, expectedResponseBody, rr.Body.Bytes())
			}
		})
	}
}

// func TestDeleteUser(t *testing.T) {
// 	gin.SetMode(gin.TestMode)
// 	data := []struct {
// 		testName           string
// 		username             string
// 		returnedError      error
// 		expectedStatusCode int
// 	}{
// 		{
// 			testName:           "Successfully delete user",
// 			username:             "0",
// 			returnedError:      nil,
// 			expectedStatusCode: http.StatusNoContent,
// 		},
// 		{
// 			testName:           "Failed to create user",
// 			username:             "0",
// 			returnedError:      fmt.Errorf("failed to create user"),
// 			expectedStatusCode: http.StatusInternalServerError,
// 		},
// 		{
// 			testName:           "No user exists",
// 			username:             "0",
// 			returnedError:      nil,
// 			expectedStatusCode: http.StatusNoContent,
// 		},
// 	}

// 	for _, d := range data {
// 		t.Run(d.testName, func(t *testing.T) {
// 			mockUserService := mocks.NewUserService(t)
// 			mockUserService.On("DeleteUser", d.username).Return(d.returnedError)
// 			userHandler := handler.UserHandlers{Service: mockUserService}

// 			rr := httptest.NewRecorder()
// 			request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/user/%s", d.username), nil)
// 			assert.NoError(t, err)

// 			router := gin.Default()
// 			router.DELETE("/user/:username", userHandler.DeleteUser)
// 			router.ServeHTTP(rr, request)

// 			assert.Equal(t, d.expectedStatusCode, rr.Code)
// 			mockUserService.AssertExpectations(t)
// 		})
// 	}
// }
