package store

import (
	"fmt"
	"strconv"
	"the-drink-almanac-api/model"
)

type FavoriteStoreStub struct {
	favorites []model.Favorite
}

func (s *FavoriteStoreStub) FindAll() ([]model.Favorite, error) {
	return s.favorites, nil
}

func (s *FavoriteStoreStub) FindFavoritesByUser(userId int) ([]model.Favorite, error) {
	filteredFavorites := make([]model.Favorite, 0)
	for _, favorite := range s.favorites {
		if favorite.UserId == userId {
			filteredFavorites = append(filteredFavorites, favorite)
		}
	}
	if len(filteredFavorites) == 0 {
		return nil, fmt.Errorf("no favorites found for user with id %s", strconv.Itoa(userId))
	}
	return filteredFavorites, nil
}

func (s *FavoriteStoreStub) CreateNewFavorite(favorite model.Favorite) error {
	s.favorites = append(s.favorites, favorite)
	return nil
}

func NewFavoriteStoreStub() (*FavoriteStoreStub, error) {
	favorites := []model.Favorite{
		{
			Id:      0,
			DrinkId: 0,
			UserId:  0,
		},
		{
			Id:      1,
			DrinkId: 1,
			UserId:  0,
		},
		{
			Id:      2,
			DrinkId: 1,
			UserId:  1,
		},
	}
	return &FavoriteStoreStub{
		favorites: favorites,
	}, nil
}
