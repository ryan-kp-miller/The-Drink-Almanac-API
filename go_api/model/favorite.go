package model

type Favorite struct {
	Id      string `json:"id" dynamodbav:"id"`
	UserId  string `json:"userId" dynamodbav:"user_id"`
	DrinkId string `json:"drinkId" dynamodbav:"drink_id"`
}

type FavoriteStore interface {
	FindAll() ([]Favorite, error)
	FindFavoritesByUser(userId string) ([]Favorite, error)
	CreateNewFavorite(favorite Favorite) error
}
