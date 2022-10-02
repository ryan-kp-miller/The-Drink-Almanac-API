package main

import (
	"fmt"
	"net/http"
	"the-drink-almanac-api/domain"
	"the-drink-almanac-api/handler"
	"the-drink-almanac-api/service"

	"github.com/gin-gonic/gin"
)

type HelloWorld struct {
	Hello string `json:"hello"`
}

func Start(port string) {
	router := gin.Default()

	// set up default endpoint
	router.GET("", hello_world_handler)

	// set up favorite endpoints
	fr, _ := domain.NewFavoriteRepositoryDDB()
	fs := service.NewDefaultFavoriteService(fr)
	fh := handler.FavoriteHandlers{Service: fs}
	router.GET("/favorites", fh.FindAllFavorites)

	// set up user endpoints
	ur, _ := domain.NewUserRepositoryDDB()
	us := service.NewDefaultUserService(ur)
	uh := handler.UserHandlers{Service: us}
	router.GET("/users", uh.FindAllUsers)

	// running the app
	router.Run(fmt.Sprintf(":%s", port))
}

func hello_world_handler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World!",
	})
}
