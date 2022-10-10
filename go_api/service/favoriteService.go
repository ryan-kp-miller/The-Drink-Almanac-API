package service

import "the-drink-almanac-api/model"

type FavoriteService interface {
	FindAllFavorites() ([]model.Favorite, error)
	FindFavoritesByUser(userID int) ([]model.Favorite, error)
	CreateNewFavorite(favorite model.Favorite) error
}

type DefaultFavoriteService struct {
	store model.FavoriteStore
}

func (s DefaultFavoriteService) FindAllFavorites() ([]model.Favorite, error) {
	return s.store.FindAll()
}

func (s DefaultFavoriteService) FindFavoritesByUser(userId int) ([]model.Favorite, error) {
	return s.store.FindFavoritesByUser(userId)
}

// CreateNewFavorite checks if a favorite already exists in the FavoriteStore
func (s DefaultFavoriteService) CreateNewFavorite(favorite model.Favorite) error {

	// need to add logic to check if a favorite already exists with the same user and drink ids
	// if it doesn't, then determine an id for the favorite and store it in the store

	return s.store.CreateNewFavorite(favorite)
}

func NewDefaultFavoriteService(store model.FavoriteStore) DefaultFavoriteService {
	return DefaultFavoriteService{store: store}
}
