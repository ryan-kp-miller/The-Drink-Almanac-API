package handler

import (
	"net/http"
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
