package handler

import (
	"fmt"
	"net/http"
	"the-drink-almanac-api/dto"
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
		return
	}

	c.JSON(http.StatusOK, users)
}

func (uh *UserHandlers) CreateNewUser(c *gin.Context) {
	var userRequest dto.UserPostRequest
	err := c.BindJSON(&userRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "the request body must contain username and password fields",
		})
		return
	}
	if userRequest.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "the username field must not be blank",
		})
		return
	}
	if userRequest.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "the password field must not be blank",
		})
		return
	}
	user, err := uh.Service.CreateNewUser(userRequest.Username, userRequest.Password)
	if err != nil {
		if user != nil {
			c.JSON(http.StatusConflict, gin.H{
				"message": fmt.Sprintf("a user already exists with the username %s", user.Username),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, user)
}
