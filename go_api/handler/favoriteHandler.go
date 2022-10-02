package handler

import (
	"net/http"
	"strconv"
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
	userIdStr := c.Param("userId")
	if userIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "you must specify a user id"})
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "the id must be an integer"})
		return
	}

	favorites, err := fh.Service.FindFavoritesByUser(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, favorites)
}
