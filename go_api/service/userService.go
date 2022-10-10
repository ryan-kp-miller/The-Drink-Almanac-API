package service

import "the-drink-almanac-api/model"

type UserService interface {
	FindAllUsers() ([]model.User, error)
}

type DefaultUserService struct {
	store model.UserStore
}

func (s DefaultUserService) FindAllUsers() ([]model.User, error) {
	return s.store.FindAll()
}

func NewDefaultUserService(store model.UserStore) DefaultUserService {
	return DefaultUserService{store: store}
}
