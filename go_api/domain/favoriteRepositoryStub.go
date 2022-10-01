package domain

type FavoriteRepositoryStub struct {
	favorites []Favorite
}

func (s FavoriteRepositoryStub) FindAll() ([]Favorite, error) {
	return s.favorites, nil
}

func NewFavoriteRepositoryStub() FavoriteRepositoryStub {
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
	}
}
