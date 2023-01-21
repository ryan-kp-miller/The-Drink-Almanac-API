package service

import (
	"fmt"
	"the-drink-almanac-api/appErrors"
	"the-drink-almanac-api/model"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	FindAllUsers() ([]model.User, error)
	CreateNewUser(username, password string) (*model.User, error)
	DeleteUser(id string) error
	Login(username, password string) (string, error)
}

type DefaultUserService struct {
	store       model.UserStore
	authService AuthService
}

// FindAllUsers returns all users in the user store
//
// This method will be restricted to admin users only
// when auth is set up
func (s DefaultUserService) FindAllUsers() ([]model.User, error) {
	return s.store.FindAll()
}

func (s DefaultUserService) getHashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CreateNewUser either creates a new user if one doesn't exist with the given username and password
// or returns the existing user and the UserAlreadyExistsError
func (s DefaultUserService) CreateNewUser(username, password string) (*model.User, error) {
	if username == "" {
		return nil, fmt.Errorf("the username must not be empty")
	}
	if password == "" {
		return nil, fmt.Errorf("the password must not be empty")
	}

	// ensure that we aren't creating a duplicate user by check for an existing record with the same username
	user, err := s.store.FindUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, appErrors.NewUserAlreadyExistsError(username)
	}

	hashedPassword, err := s.getHashedPassword(password)
	if err != nil {
		return nil, err
	}

	user = &model.User{
		Id:       uuid.NewString(),
		Username: username,
		Password: hashedPassword,
	}
	err = s.store.CreateNewUser(*user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUser removes the user's record from the user store
// based on the provided id
func (s DefaultUserService) DeleteUser(id string) error {
	return s.store.DeleteUser(id)
}

// Login checks if a user exists with the provided username and password;
// if a user exists:
//   - and the password matches, returns a JWT string
//   - and the password doesn't match, returns an error
// if a user doesn't exist with the given username, returns the UserNotFoundError
func (s DefaultUserService) Login(username, password string) (string, error) {
	if username == "" {
		return "", fmt.Errorf("the username must not be empty")
	}
	if password == "" {
		return "", fmt.Errorf("the password must not be empty")
	}

	// ensure that we aren't creating a duplicate user by check for an existing record with the same username
	user, err := s.store.FindUserByUsername(username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", appErrors.NewUserNotFoundError(username)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", appErrors.NewIncorrectPasswordError(user.Username)
	}

	// set the ttl to 24 hours
	return s.authService.CreateNewToken(user.Id, 60*24)
}

func NewDefaultUserService(store model.UserStore) DefaultUserService {
	return DefaultUserService{
		store:       store,
		authService: NewJwtAuthService(),
	}
}
