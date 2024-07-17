package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"the-drink-almanac-api/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func checkIfUserIdContextVarIsSet(t *testing.T, expectedUserId string) func(c *gin.Context) {
	return func(c *gin.Context) {
		actualUserId := c.GetString("userId")
		assert.Equal(t, expectedUserId, actualUserId)
	}
}

func TestAuthUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	data := []struct {
		testName           string
		userId             string
		token              string
		authError          error
		isTokenExpected    bool
		expectedStatusCode int
	}{
		{
			testName:           "Successfully retrieved user id",
			userId:             "0",
			token:              "testToken",
			authError:          nil,
			isTokenExpected:    true,
			expectedStatusCode: http.StatusOK,
		},
		{
			testName:           "No token in request headers",
			userId:             "0",
			token:              "",
			authError:          nil,
			isTokenExpected:    false,
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			testName:           "Invalid token",
			userId:             "0",
			token:              "testToken",
			authError:          fmt.Errorf("invalid token"),
			isTokenExpected:    false,
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			mockAuthService := service.NewMockAuthService(t)
			if d.token != "" {
				mockAuthService.On("ValidateToken", d.token).Return(d.userId, d.authError)
			}
			authMiddleware := NewAuthMiddleware(mockAuthService)

			rr := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, "/user", nil)
			assert.NoError(t, err)
			if d.token != "" {
				request.Header["Token"] = []string{d.token}
			}

			router := gin.Default()
			router.GET("/user", authMiddleware.AuthUser, checkIfUserIdContextVarIsSet(t, d.userId))
			router.ServeHTTP(rr, request)

			assert.Equal(t, d.expectedStatusCode, rr.Code)
			mockAuthService.AssertExpectations(t)
		})
	}
}
