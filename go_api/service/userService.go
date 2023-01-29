package service

import (
	"fmt"
	"the-drink-almanac-api/appErrors"
	"the-drink-almanac-api/model"
	"the-drink-almanac-api/store"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	// FindAllUsers returns all users in the user store
	FindAllUsers() ([]model.User, error)

	// FindUser retrieves the user's data based on their id
	FindUser(userId string) (*model.User, error)

	// CreateNewUser either creates a new user if one doesn't exist with the given username and password
	// or returns the existing user and the UserAlreadyExistsError
	CreateNewUser(username, password string) (*model.User, error)

	// DeleteUser removes the user's record from the user store
	DeleteUser(userId string) error

	// Login checks if a user exists with the provided username and password;
	// if a user exists:
	//   - and the password matches, returns the user's data
	//   - and the password doesn't match, returns an error
	// if a user doesn't exist with the given username, returns the UserNotFoundError
	Login(username, password string) (*model.User, error)
}

type DefaultUserService struct {
	store store.UserStore
}

func (s DefaultUserService) FindAllUsers() ([]model.User, error) {
	return s.store.FindAll()
}

func (s DefaultUserService) FindUser(userId string) (*model.User, error) {
	return s.store.FindUserById(userId)
}

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

	hashedPassword, err := getHashedPassword(password)
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

func (s DefaultUserService) DeleteUser(userId string) error {
	return s.store.DeleteUser(userId)
}

func (s DefaultUserService) Login(username, password string) (*model.User, error) {
	if username == "" {
		return nil, fmt.Errorf("the username must not be empty")
	}
	if password == "" {
		return nil, fmt.Errorf("the password must not be empty")
	}

	user, err := s.store.FindUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, appErrors.NewUserNotFoundError(username)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, appErrors.NewIncorrectPasswordError(user.Username)
	}

	return user, nil
}

func NewDefaultUserService(store store.UserStore) DefaultUserService {
	return DefaultUserService{
		store: store,
	}
}

// getHashedPassword takes the raw password and hashes it for protection
func getHashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
