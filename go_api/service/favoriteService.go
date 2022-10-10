package service

import "the-drink-almanac-api/model"

type FavoriteService interface {
	FindAllFavorites() ([]model.Favorite, error)
	FindFavoritesByUser(userID int) ([]model.Favorite, error)
	CreateNewFavorite(favorite model.Favorite) error
}

type DefaultFavoriteService struct {
	repo model.FavoriteStore
}

func (s DefaultFavoriteService) FindAllFavorites() ([]model.Favorite, error) {
	return s.repo.FindAll()
}

func (s DefaultFavoriteService) FindFavoritesByUser(userId int) ([]model.Favorite, error) {
	return s.repo.FindFavoritesByUser(userId)
}

func (s DefaultFavoriteService) CreateNewFavorite(favorite model.Favorite) error {
	// need to add logic to check if a favorite already exists with the same user and drink ids
	// if it doesn't, then determine an id for the favorite and store it in the repo
	return s.repo.CreateNewFavorite(favorite)
}

func NewDefaultFavoriteService(repository model.FavoriteStore) DefaultFavoriteService {
	return DefaultFavoriteService{repo: repository}
}
