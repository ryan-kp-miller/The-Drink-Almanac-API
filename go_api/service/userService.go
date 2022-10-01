package service

import "the-drink-almanac-api/domain"

type UserService interface {
	FindAllUsers() ([]domain.User, error)
}

type DefaultUserService struct {
	repo domain.UserRepository
}

func (s DefaultUserService) FindAllUsers() ([]domain.User, error) {
	return s.repo.FindAll()
}

func NewDefaultUserService(repository domain.UserRepository) DefaultUserService {
	return DefaultUserService{repo: repository}
}
