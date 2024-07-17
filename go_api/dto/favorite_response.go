package dto

import "the-drink-almanac-api/model"

type FavoriteResponse struct {
	Id      string `json:"id"`
	DrinkId string `json:"drinkId"`
}

func NewFavoriteResponse(favorite model.Favorite) FavoriteResponse {
	return FavoriteResponse{
		Id:      favorite.Id,
		DrinkId: favorite.DrinkId,
	}
}

func NewFavoritesResponse(favorites []model.Favorite) []FavoriteResponse {
	favoritesResponse := make([]FavoriteResponse, len(favorites))
	for i, user := range favorites {
		favoritesResponse[i] = NewFavoriteResponse(user)
	}
	return favoritesResponse
}
