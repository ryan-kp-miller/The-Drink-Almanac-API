package handler

import (
	"net/http"
	"the-drink-almanac-api/service"

	"github.com/gin-gonic/gin"
)

type UserHandlers struct {
	Service service.UserService
}

func (uh *UserHandlers) FindAllUsers(c *gin.Context) {
	users, err := uh.Service.FindAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, users)
}
