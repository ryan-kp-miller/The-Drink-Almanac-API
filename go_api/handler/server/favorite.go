package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"the-drink-almanac-api/apperrors"
	"the-drink-almanac-api/dto"
	"the-drink-almanac-api/service"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	Service service.FavoriteService
}

func (fh *FavoriteHandler) FindAllFavorites(c *gin.Context) {
	favorites, err := fh.Service.FindAllFavorites()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	favoritesResponse := dto.NewFavoritesResponse(favorites)
	c.JSON(http.StatusOK, favoritesResponse)
}

func (fh *FavoriteHandler) FindFavoritesByUser(c *gin.Context) {
	userId := c.GetString("userId")
	favorites, err := fh.Service.FindFavoritesByUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if len(favorites) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("no favorites were found for user with id %s", userId)})
		return
	}
	favoritesResponse := dto.NewFavoritesResponse(favorites)
	c.JSON(http.StatusOK, favoritesResponse)
}

func (fh *FavoriteHandler) CreateNewFavorite(c *gin.Context) {
	var newFavoritePostRequest dto.FavoritePostRequest
	if err := c.BindJSON(&newFavoritePostRequest); err != nil {
		err_msg := err.Error()
		if strings.Contains(err_msg, "cannot unmarshal number into Go struct field FavoritePostRequest.drink_id of type string") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "the drink_id must be a string"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "please provide the drink_id in the body of your request"})
		return
	}
	if err := newFavoritePostRequest.ValidateRequest(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	userId := c.GetString("userId")
	newFavorite, err := fh.Service.CreateNewFavorite(newFavoritePostRequest.DrinkId, userId)
	if err != nil {
		if errors.As(err, &apperrors.FavoriteAlreadyExistsError{}) {
			c.JSON(http.StatusConflict, gin.H{
				"message": fmt.Sprintf("the user '%s' already favorited the drink with id '%s'", newFavorite.UserId, newFavorite.DrinkId),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("unable to add the new favorite due to %s", err.Error())})
		return
	}
	favoriteResponse := dto.NewFavoriteResponse(*newFavorite)
	c.JSON(http.StatusCreated, favoriteResponse)
}

func (fh *FavoriteHandler) DeleteFavorite(c *gin.Context) {
	favoriteId := c.Param("favoriteId")
	err := fh.Service.DeleteFavorite(favoriteId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "the favorite was deleted"})
}
