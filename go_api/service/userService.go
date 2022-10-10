package service

import "the-drink-almanac-api/model"

type UserService interface {
	FindAllUsers() ([]model.User, error)
}

type DefaultUserService struct {
	repo model.UserStore
}

func (s DefaultUserService) FindAllUsers() ([]model.User, error) {
	return s.repo.FindAll()
}

func NewDefaultUserService(repository model.UserStore) DefaultUserService {
	return DefaultUserService{repo: repository}
}
