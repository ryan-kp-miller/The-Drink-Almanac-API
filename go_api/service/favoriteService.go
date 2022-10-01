package service

import "the-drink-almanac-api/domain"

type FavoriteService interface {
	FindAllFavorites() ([]domain.Favorite, error)
}

type DefaultFavoriteService struct {
	repo domain.FavoriteRepository
}

func (s DefaultFavoriteService) FindAllFavorites() ([]domain.Favorite, error) {
	return s.repo.FindAll()
}

func NewDefaultFavoriteService(repository domain.FavoriteRepository) DefaultFavoriteService {
	return DefaultFavoriteService{repo: repository}
}
