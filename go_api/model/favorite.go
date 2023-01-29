package model

type Favorite struct {
	Id      string `dynamodbav:"id"`
	UserId  string `dynamodbav:"user_id"`
	DrinkId string `dynamodbav:"drink_id"`
}

type FavoriteStore interface {
	FindAll() ([]Favorite, error)
	FindFavoritesByUser(userId string) ([]Favorite, error)
	CreateNewFavorite(favorite Favorite) error
	DeleteFavorite(id string) error
}
