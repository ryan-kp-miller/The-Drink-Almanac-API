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

type UserHandler struct {
	userService service.UserService
	authService service.AuthService
}

func (uh *UserHandler) FindAllUsers(c *gin.Context) {
	users, err := uh.userService.FindAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	usersResponse := dto.NewUsersResponse(users)
	c.JSON(http.StatusOK, usersResponse)
}

func (uh *UserHandler) FindUser(c *gin.Context) {
	userId := c.GetString("userId")
	if userId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "token was not successfully retrieved"})
	}

	user, err := uh.userService.FindUser(userId)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.As(err, &appErrors.InvalidAuthTokenError{}) {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"message": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("no user was found with the userId '%s'", userId)})
		return
	}

	userResponse := dto.NewUserResponse(*user)
	c.JSON(http.StatusOK, userResponse)
}

func (uh *UserHandler) CreateNewUser(c *gin.Context) {
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
	user, err := uh.userService.CreateNewUser(userRequest.Username, userRequest.Password)
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

func (uh *UserHandler) DeleteUser(c *gin.Context) {
	userId := c.GetString("userId")
	if userId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "token was not successfully retrieved"})
	}

	err := uh.userService.DeleteUser(userId)
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

func (uh *UserHandler) Login(c *gin.Context) {
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

	user, err := uh.userService.Login(userRequest.Username, userRequest.Password)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.As(err, &appErrors.UserNotFoundError{}) {
			statusCode = http.StatusNotFound
		}
		if errors.As(err, &appErrors.IncorrectPasswordError{}) {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{"message": err.Error()})
		return
	}

	tokenString, err := uh.authService.CreateNewToken(user.Id, 60*24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Header("Token", tokenString)
	c.Status(http.StatusAccepted)
}

func NewUserHandler(userService service.UserService, authService service.AuthService) UserHandler {
	return UserHandler{
		userService: userService,
		authService: authService,
	}
}
