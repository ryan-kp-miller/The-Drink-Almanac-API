package service

import "the-drink-almanac-api/domain"

type FavoriteService interface {
	FindAllFavorites() ([]domain.Favorite, error)
	FindFavoritesByUser(userID int) ([]domain.Favorite, error)
	CreateNewFavorite(favorite domain.Favorite) error
}

type DefaultFavoriteService struct {
	repo domain.FavoriteRepository
}

func (s DefaultFavoriteService) FindAllFavorites() ([]domain.Favorite, error) {
	return s.repo.FindAll()
}

func (s DefaultFavoriteService) FindFavoritesByUser(userId int) ([]domain.Favorite, error) {
	return s.repo.FindFavoritesByUser(userId)
}

func (s DefaultFavoriteService) CreateNewFavorite(favorite domain.Favorite) error {
	// need to add logic to check if a favorite already exists with the same user and drink ids
	// if it doesn't, then determine an id for the favorite and store it in the repo
	return s.repo.CreateNewFavorite(favorite)
}

func NewDefaultFavoriteService(repository domain.FavoriteRepository) DefaultFavoriteService {
	return DefaultFavoriteService{repo: repository}
}
