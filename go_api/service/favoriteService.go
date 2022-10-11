package service

import (
	"fmt"
	"the-drink-almanac-api/model"
)

type FavoriteService interface {
	FindAllFavorites() ([]model.Favorite, error)
	FindFavoritesByUser(userID string) ([]model.Favorite, error)
	CreateNewFavorite(userId, drinkId string) (*model.Favorite, error)
}

type DefaultFavoriteService struct {
	store model.FavoriteStore
}

func (s DefaultFavoriteService) FindAllFavorites() ([]model.Favorite, error) {
	return s.store.FindAll()
}

func (s DefaultFavoriteService) FindFavoritesByUser(userId string) ([]model.Favorite, error) {
	return s.store.FindFavoritesByUser(userId)
}

// CreateNewFavorite checks if a favorite already exists in the FavoriteStore
// with the same user and drink ids
func (s DefaultFavoriteService) CreateNewFavorite(drinkId, userId string) (*model.Favorite, error) {
	// check if a favorite already exists for this drink/user id pair
	// userFavorites, err := s.store.FindFavoritesByUser(userId)

	// newFavorite := model.Favorite{}
	// need to add logic to check if a favorite already exists with the same user and drink ids
	// if it doesn't, then determine an id for the favorite and store it in the store

	// return s.store.CreateNewFavorite(favorite), nil
	return nil, fmt.Errorf("not implemented yet")
}

func NewDefaultFavoriteService(store model.FavoriteStore) DefaultFavoriteService {
	return DefaultFavoriteService{store: store}
}
