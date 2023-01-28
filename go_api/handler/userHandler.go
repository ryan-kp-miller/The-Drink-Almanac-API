package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"the-drink-almanac-api/appErrors"
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
	usersResponse := dto.NewUsersResponse(users)
	c.JSON(http.StatusOK, usersResponse)
}

func (uh *UserHandlers) FindUser(c *gin.Context) {
	tokens := c.Request.Header["Token"]
	if len(tokens) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "the 'Token' header was not included in the request"})
		return
	}

	tokenString := tokens[0]
	user, err := uh.Service.FindUser(tokenString)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.As(err, &appErrors.InvalidAuthTokenError{}) {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"message": err.Error()})
		return
	}
	userResponse := dto.NewUserResponse(*user)
	c.JSON(http.StatusOK, userResponse)
}

func (uh *UserHandlers) CreateNewUser(c *gin.Context) {
	var userRequest dto.UserPostRequest
	err := c.BindJSON(&userRequest)
	if err != nil {
		err_msg := err.Error()

		if strings.Contains(err_msg, "cannot unmarshal number into Go struct field UserPostRequest.username of type string") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "the username must be a string"})
			return
		}
		if strings.Contains(err_msg, "cannot unmarshal number into Go struct field UserPostRequest.password of type string") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "the password must be a string"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "please provide the username and password in the body of your request"})
		return
	}
	if err = userRequest.ValidateRequest(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	user, err := uh.Service.CreateNewUser(userRequest.Username, userRequest.Password)
	if err != nil {
		if errors.As(err, &appErrors.UserAlreadyExistsError{}) {
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
	c.Status(http.StatusCreated)
}

func (uh *UserHandlers) DeleteUser(c *gin.Context) {
	tokens := c.Request.Header["Token"]
	if len(tokens) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "the 'Token' header was not included in the request"})
		return
	}

	tokenString := tokens[0]
	err := uh.Service.DeleteUser(tokenString)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.As(err, &appErrors.InvalidAuthTokenError{}) {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "the user was deleted"})
}

func (uh *UserHandlers) Login(c *gin.Context) {
	var userRequest dto.UserPostRequest
	err := c.BindJSON(&userRequest)
	if err != nil {
		err_msg := err.Error()

		if strings.Contains(err_msg, "cannot unmarshal number into Go struct field UserPostRequest.username of type string") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "the username must be a string"})
			return
		}
		if strings.Contains(err_msg, "cannot unmarshal number into Go struct field UserPostRequest.password of type string") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "the password must be a string"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "please provide the username and password in the body of your request"})
		return
	}
	if err = userRequest.ValidateRequest(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	tokenString, err := uh.Service.Login(userRequest.Username, userRequest.Password)
	if err != nil {
		statusCode := http.StatusInternalServerError
		// I tried using a switch statement to clean this up,
		// but it wasn't picking up my custom errors types
		if errors.As(err, &appErrors.UserNotFoundError{}) {
			statusCode = http.StatusNotFound
		}
		if errors.As(err, &appErrors.IncorrectPasswordError{}) {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{"message": err.Error()})
		return
	}

	c.Header("Token", tokenString)
	c.Status(http.StatusAccepted)
}
