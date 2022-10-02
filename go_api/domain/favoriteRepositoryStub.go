package domain

import (
	"fmt"
	"strconv"
)

type FavoriteRepositoryStub struct {
	favorites []Favorite
}

func (s FavoriteRepositoryStub) FindAll() ([]Favorite, error) {
	return s.favorites, nil
}

func (s FavoriteRepositoryStub) FindFavoritesByUser(userId int) ([]Favorite, error) {
	filteredFavorites := make([]Favorite, 0)
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

func NewFavoriteRepositoryStub() (FavoriteRepositoryStub, error) {
	favorites := []Favorite{
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
	return FavoriteRepositoryStub{
		favorites: favorites,
	}, nil
}
