package domain

type Favorite struct {
	Id      int `json:"id" db:"id"`
	DrinkId int `json:"drinkId" db:"drink_id"`
	UserId  int `json:"userId" db:"user_id"`
}

type FavoriteRepository interface {
	FindAll() ([]Favorite, error)
	FindFavoritesByUser(userId int) ([]Favorite, error)
	// need to add a method to determine if a record already exists with a given user id and drink id
	// need to add a method to find the max favorite id in order to determine what to use for new favorites
	CreateNewFavorite(favorite Favorite) error
}
