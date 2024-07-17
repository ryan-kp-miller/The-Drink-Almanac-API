package service

import (
	"testing"
	"the-drink-almanac-api/service"

	"github.com/stretchr/testify/assert"
)

func TestJwtAuthService(t *testing.T) {
	tests := []struct {
		name         string
		userId       string
		tokenExpiry  int
		expectError  bool
		useFakeToken bool
	}{
		{
			name:         "Valid token, successfully retrieve user id",
			userId:       "testId",
			tokenExpiry:  10,
			expectError:  false,
			useFakeToken: false,
		},
		{
			name:         "Expired token returns error",
			userId:       "testId",
			tokenExpiry:  -10,
			expectError:  true,
			useFakeToken: false,
		},
		{
			name:         "Bad token is marked invalid",
			userId:       "testId",
			tokenExpiry:  10,
			expectError:  true,
			useFakeToken: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService := service.NewJwtAuthService("testToken")

			var tokenString string
			if !tt.useFakeToken {
				tokenString, _ = authService.CreateNewToken(tt.userId, tt.tokenExpiry)
			} else {
				tokenString = "fakeToken"
			}

			actualUserId, err := authService.ValidateToken(tokenString)
			if tt.expectError {
				assert.NotNil(t, err, "An error should have been returned from authService.ValidateToken")
			} else {
				assert.Nil(t, err, "No error should have been returned from authService.ValidateToken")
				assert.Equal(t, tt.userId, actualUserId, "The extracted userId does not match the one put into the token")
			}
		})
	}

}
