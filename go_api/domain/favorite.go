package domain

type Favorite struct {
	Id      int `json:"id" db:"id"`
	DrinkId int `json:"drinkId" db:"drink_id"`
	UserId  int `json:"userId" db:"user_id"`
}

type FavoriteRepository interface {
	FindAll() ([]Favorite, error)
	FindFavoritesByUser(userId int) ([]Favorite, error)
}
