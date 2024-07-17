//go:generate mockery --name=FavoriteService --output=./ --outpkg=service --filename=favorite_mock.go --inpackage
package service

import (
	"fmt"

	"the-drink-almanac-api/apperrors"
	"the-drink-almanac-api/model"
	"the-drink-almanac-api/repository"

	"github.com/google/uuid"
)

type FavoriteService interface {
	FindAllFavorites() ([]model.Favorite, error)

	// FindFavoritesByUser retrieves favorites based on the user's id
	FindFavoritesByUser(userId string) ([]model.Favorite, error)

	// CreateNewFavorite either creates a new favorite if one doesn't exist with the given drinkId and userId
	// or returns the existing favorite and the FavoriteAlreadyExistsError
	CreateNewFavorite(userId, drinkId string) (*model.Favorite, error)

	DeleteFavorite(id string) error
}

func NewDefaultFavoriteService(repo repository.FavoriteRepository) DefaultFavoriteService {
	return DefaultFavoriteService{repo: repo}
}

type DefaultFavoriteService struct {
	repo repository.FavoriteRepository
}

func (s DefaultFavoriteService) FindAllFavorites() ([]model.Favorite, error) {
	return s.repo.FindAll()
}

func (s DefaultFavoriteService) FindFavoritesByUser(userId string) ([]model.Favorite, error) {
	return s.repo.FindFavoritesByUser(userId)
}

func (s DefaultFavoriteService) CreateNewFavorite(drinkId, userId string) (*model.Favorite, error) {
	if drinkId == "" {
		return nil, fmt.Errorf("the drinkId must not be empty")
	}
	if userId == "" {
		return nil, fmt.Errorf("the userId must not be empty")
	}

	// check if a favorite already exists for this drink/user id pair
	userFavorites, err := s.repo.FindFavoritesByUser(userId)
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
		return existingFavorite, apperrors.NewFavoriteAlreadyExistsError("the user already favorited this drink")
	}

	favoriteUuid := uuid.New()
	newFavorite := model.Favorite{
		Id:      favoriteUuid.String(),
		UserId:  userId,
		DrinkId: drinkId,
	}
	err = s.repo.CreateNewFavorite(newFavorite)
	if err != nil {
		return nil, err
	}

	return &newFavorite, nil
}

func (s DefaultFavoriteService) DeleteFavorite(id string) error {
	return s.repo.DeleteFavorite(id)
}
