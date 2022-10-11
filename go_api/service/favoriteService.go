package service

import (
	"fmt"
	"the-drink-almanac-api/model"

	"github.com/google/uuid"
)

type FavoriteService interface {
	FindAllFavorites() ([]model.Favorite, error)
	FindFavoritesByUser(userID string) ([]model.Favorite, error)
	CreateNewFavorite(userId, drinkId string) (*model.Favorite, error)
	DeleteFavorite(id string) error
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
	userFavorites, err := s.store.FindFavoritesByUser(userId)
	if err != nil {
		return nil, err
	}

	doesFavoriteExist := false
	var existingFavorite *model.Favorite
	for _, favorite := range userFavorites {
		if favorite.DrinkId == drinkId && favorite.UserId == userId {
			doesFavoriteExist = true
			existingFavorite = &favorite
			break
		}
	}
	if doesFavoriteExist {
		return existingFavorite, fmt.Errorf("the user already favorited this drink")
	}

	favoriteUuid := uuid.New()
	newFavorite := model.Favorite{
		Id:      favoriteUuid.String(),
		UserId:  userId,
		DrinkId: drinkId,
	}
	err = s.store.CreateNewFavorite(newFavorite)
	if err != nil {
		return nil, err
	}

	return &newFavorite, nil
}

func (s DefaultFavoriteService) DeleteFavorite(id string) error {
	return s.store.DeleteFavorite(id)
}

func NewDefaultFavoriteService(store model.FavoriteStore) DefaultFavoriteService {
	return DefaultFavoriteService{store: store}
}
