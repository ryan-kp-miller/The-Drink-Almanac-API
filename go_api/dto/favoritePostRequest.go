package dto

type FavoritePostRequest struct {
	UserId  string `json:"user_id"`
	DrinkId string `json:"drink_id"`
}
