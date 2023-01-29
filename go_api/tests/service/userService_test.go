package service

import (
	"fmt"
	"testing"
	"the-drink-almanac-api/appErrors"
	"the-drink-almanac-api/mocks"
	"the-drink-almanac-api/model"
	"the-drink-almanac-api/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestDefaultUserService_FindAllUsers(t *testing.T) {
	mockUsers := []model.User{
		{
			Id:       "0",
			Username: "0",
			Password: "0",
		},
		{
			Id:       "1",
			Username: "1",
			Password: "1",
		},
		{
			Id:       "2",
			Username: "2",
			Password: "2",
		},
	}
	tests := []struct {
		name          string
		returnedUsers []model.User
		returnedError error
		expectError   bool
	}{
		{
			name:          "Successfully retrieved users",
			returnedUsers: mockUsers,
			returnedError: nil,
			expectError:   false,
		},
		{
			name:          "Failed to retrieve users",
			returnedUsers: mockUsers,
			returnedError: fmt.Errorf("failed to retrieve users"),
			expectError:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore := mocks.NewUserStore(t)
			mockUserStore.On("FindAll").Return(tt.returnedUsers, tt.returnedError)
			userService := service.NewDefaultUserService(mockUserStore)
			users, err := userService.FindAllUsers()
			assert.Equal(t, users, tt.returnedUsers, "The users returned by FindAllUsers() does not match the expected users; actual users = %v; expected users = %v", users, tt.returnedUsers)
			if tt.expectError {
				assert.NotNil(t, err, "An error should have been returned from userService.FindAllUsers")
			} else {
				assert.Nil(t, err, "No error should have been returned from userService.FindAllUsers")
			}
			mockUserStore.AssertExpectations(t)
		})
	}
}

func TestDefaultUserService_CreateNewUser(t *testing.T) {
	mockUser := model.User{
		Id:       "0",
		Username: "0",
		Password: "0",
	}
	tests := []struct {
		name                            string
		username                        string
		password                        string
		isStoreCreateNewUserCalled      bool
		isStoreFindUserByUsernameCalled bool
		returnedError                   error
		existingUser                    *model.User
		existingUserError               error
		expectError                     bool
	}{
		{
			name:                            "Successfully create user",
			username:                        "0",
			password:                        "0",
			isStoreCreateNewUserCalled:      true,
			isStoreFindUserByUsernameCalled: true,
			returnedError:                   nil,
			existingUser:                    nil,
			existingUserError:               nil,
			expectError:                     false,
		},
		{
			name:                            "Failed to create users",
			username:                        "0",
			password:                        "0",
			isStoreCreateNewUserCalled:      true,
			isStoreFindUserByUsernameCalled: true,
			returnedError:                   fmt.Errorf("failed to retrieve users"),
			existingUser:                    nil,
			existingUserError:               nil,
			expectError:                     true,
		},
		{
			name:                            "User already exists",
			username:                        "0",
			password:                        "0",
			isStoreCreateNewUserCalled:      false,
			isStoreFindUserByUsernameCalled: true,
			returnedError:                   nil,
			existingUser:                    &mockUser,
			existingUserError:               fmt.Errorf("user already exists"),
			expectError:                     true,
		},
		{
			name:                            "Username is empty",
			username:                        "",
			password:                        "0",
			isStoreCreateNewUserCalled:      false,
			isStoreFindUserByUsernameCalled: false,
			returnedError:                   nil,
			existingUser:                    nil,
			existingUserError:               nil,
			expectError:                     true,
		},
		{
			name:                            "Password is empty",
			username:                        "0",
			password:                        "",
			isStoreCreateNewUserCalled:      false,
			isStoreFindUserByUsernameCalled: false,
			returnedError:                   nil,
			existingUser:                    nil,
			existingUserError:               nil,
			expectError:                     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore := mocks.NewUserStore(t)
			if tt.isStoreFindUserByUsernameCalled {
				mockUserStore.On("FindUserByUsername", tt.username).Return(tt.existingUser, tt.existingUserError)
			}
			if tt.isStoreCreateNewUserCalled {
				mockUserStore.On("CreateNewUser", mock.AnythingOfType("model.User")).Return(tt.returnedError)
			}

			userService := service.NewDefaultUserService(mockUserStore)
			user, err := userService.CreateNewUser(tt.username, tt.password)
			if tt.expectError {
				assert.NotNil(t, err, "An error should have been returned from userService.CreateNewUser")
			} else {
				assert.Nil(t, err, "No error should have been returned from userService.CreateNewUser")
			}

			// check that the new user has a uuid for Id, Username matches the provided username,
			// and the Password is the hashed version of the provided password
			if user != nil {
				_, err = uuid.Parse(user.Id)
				assert.Nil(t, err, fmt.Sprintf("The user.Id field must be a uuid; actual value: %s", user.Id))
				assert.Equal(t, user.Username, tt.username)
				err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(tt.password))
				assert.Nil(t, err, "The user.Password field was not hashed correctly; actual value: %s", user.Password)
			}
			mockUserStore.AssertExpectations(t)
		})
	}
}

func TestDefaultUserService_DeleteUser(t *testing.T) {
	tests := []struct {
		name               string
		userId             string
		storeReturnedError error
		authReturnedError  error
		expectError        bool
	}{
		{
			name:               "Successfully delete user",
			userId:             "0",
			storeReturnedError: nil,
			authReturnedError:  nil,
			expectError:        false,
		},
		{
			name:               "Failed to delete user",
			userId:             "0",
			storeReturnedError: fmt.Errorf("failed to delete users"),
			authReturnedError:  nil,
			expectError:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore := mocks.NewUserStore(t)
			if tt.authReturnedError == nil {
				mockUserStore.On("DeleteUser", tt.userId).Return(tt.storeReturnedError)
			}

			userService := service.NewDefaultUserService(mockUserStore)
			err := userService.DeleteUser(tt.userId)

			if tt.expectError {
				assert.NotNil(t, err, "An error should have been returned from userService.DeleteUser")
			} else {
				assert.Nil(t, err, "No error should have been returned from userService.DeleteUser")
			}
			mockUserStore.AssertExpectations(t)
		})
	}
}

func TestDefaultUserService_Login(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("0"), 8)
	mockUser := &model.User{
		Id:       "0",
		Username: "0",
		Password: string(hashedPassword),
	}
	tests := []struct {
		name                            string
		username                        string
		password                        string
		isStoreFindUserByUsernameCalled bool
		FindUserByUsernameError         error
		existingUser                    *model.User
		expectedError                   error
	}{
		{
			name:                            "Successfully logged in",
			username:                        "0",
			password:                        "0",
			isStoreFindUserByUsernameCalled: true,
			FindUserByUsernameError:         nil,
			existingUser:                    mockUser,
			expectedError:                   nil,
		},
		{
			name:                            "Failed to log in",
			username:                        "0",
			password:                        "0",
			isStoreFindUserByUsernameCalled: true,
			FindUserByUsernameError:         fmt.Errorf("failed to log in"),
			existingUser:                    nil,
			expectedError:                   fmt.Errorf("failed to log in"),
		},
		{
			name:                            "User not found",
			username:                        "0",
			password:                        "0",
			isStoreFindUserByUsernameCalled: true,
			FindUserByUsernameError:         nil,
			existingUser:                    nil,
			expectedError:                   appErrors.NewUserNotFoundError("0"),
		},
		{
			name:                            "Username is empty",
			username:                        "",
			password:                        "0",
			isStoreFindUserByUsernameCalled: false,
			FindUserByUsernameError:         nil,
			existingUser:                    nil,
			expectedError:                   fmt.Errorf("the username must not be empty"),
		},
		{
			name:                            "Password is empty",
			username:                        "0",
			password:                        "",
			isStoreFindUserByUsernameCalled: false,
			FindUserByUsernameError:         nil,
			existingUser:                    nil,
			expectedError:                   fmt.Errorf("the password must not be empty"),
		},
		{
			name:                            "Password doesn't match",
			username:                        "0",
			password:                        "badPassword",
			isStoreFindUserByUsernameCalled: true,
			FindUserByUsernameError:         nil,
			existingUser:                    mockUser,
			expectedError:                   appErrors.NewIncorrectPasswordError("0"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore := mocks.NewUserStore(t)
			if tt.isStoreFindUserByUsernameCalled {
				mockUserStore.On("FindUserByUsername", tt.username).Return(tt.existingUser, tt.FindUserByUsernameError)
			}

			userService := service.NewDefaultUserService(mockUserStore)
			user, err := userService.Login(tt.username, tt.password)
			assert.Equal(t, tt.expectedError, err)
			if tt.expectedError == nil {
				assert.Equal(t, user, tt.existingUser)
			}
			mockUserStore.AssertExpectations(t)
		})
	}
}
