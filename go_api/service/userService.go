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
}

type DefaultUserService struct {
	store model.UserStore
}

// FindAllUsers returns all users in the user store
//
// This method will be restricted to admin users only
// when auth is set up
func (s DefaultUserService) FindAllUsers() ([]model.User, error) {
	return s.store.FindAll()
}

// CreateNewUser either creates a new user if one doesn't exist with the given username and password
// or returns the existing user and the UserAlreadyExistsError
func (s DefaultUserService) CreateNewUser(username, password string) (*model.User, error) {
	// ensure that we aren't creating a duplicate user by check for an existing record with the same username
	user, err := s.store.FindUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, appErrors.NewUserAlreadyExistsError(fmt.Sprintf("a user already exists with the username %s", username))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return nil, err
	}

	user = &model.User{
		Id:       uuid.NewString(),
		Username: username,
		Password: string(hashedPassword),
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

func NewDefaultUserService(store model.UserStore) DefaultUserService {
	return DefaultUserService{store: store}
}
