package service

import (
	"fmt"
	"the-drink-almanac-api/model"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	FindAllUsers() ([]model.User, error)
	CreateNewUser(username, password string) (*model.User, error)
}

type DefaultUserService struct {
	store model.UserStore
}

func (s DefaultUserService) FindAllUsers() ([]model.User, error) {
	return s.store.FindAll()
}

func (s DefaultUserService) CreateNewUser(username, password string) (*model.User, error) {
	// ensure that we aren't creating a duplicate user by check for an existing record with the same username
	user, err := s.store.FindUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, fmt.Errorf("a user already exists with the username %s", username)
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

func NewDefaultUserService(store model.UserStore) DefaultUserService {
	return DefaultUserService{store: store}
}
