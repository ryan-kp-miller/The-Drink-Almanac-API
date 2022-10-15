package handler

import (
	"fmt"
	"net/http"
	"the-drink-almanac-api/dto"
	"the-drink-almanac-api/service"

	"github.com/gin-gonic/gin"
)

type FavoriteHandlers struct {
	Service service.FavoriteService
}

func (fh *FavoriteHandlers) FindAllFavorites(c *gin.Context) {
	favorites, err := fh.Service.FindAllFavorites()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	c.JSON(http.StatusOK, favorites)
}

func (fh *FavoriteHandlers) FindFavoritesByUser(c *gin.Context) {
	userId := c.Param("userId")
	favorites, err := fh.Service.FindFavoritesByUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if len(favorites) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("no favorites were found for user with id %s", userId)})
		return
	}
	c.JSON(http.StatusOK, favorites)
}

func (fh *FavoriteHandlers) CreateNewFavorite(c *gin.Context) {
	var newFavoritePostRequest dto.FavoritePostRequest
	if err := c.BindJSON(&newFavoritePostRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "please provide the drink_id and user_id in the body of your request"})
		return
	}
	if err := newFavoritePostRequest.ValidateRequest(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	newFavorite, err := fh.Service.CreateNewFavorite(newFavoritePostRequest.UserId, newFavoritePostRequest.DrinkId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("unable to add the new favorite due to %s", err.Error())})
		return
	}

	c.JSON(http.StatusCreated, newFavorite)
}

func (fh *FavoriteHandlers) DeleteFavorite(c *gin.Context) {
	favoriteId := c.Param("favoriteId")
	if favoriteId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "you must specify an id"})
		return
	}
	err := fh.Service.DeleteFavorite(favoriteId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "the favorite was deleted"})
}
