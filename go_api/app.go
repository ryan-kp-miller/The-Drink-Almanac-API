package main

import (
	"fmt"
	"net/http"
	"the-drink-almanac-api/handler"
	"the-drink-almanac-api/service"
	"the-drink-almanac-api/store"

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
	fr, _ := store.NewFavoriteStoreDDB()
	fs := service.NewDefaultFavoriteService(fr)
	fh := handler.FavoriteHandlers{Service: fs}
	favoriteRouteGroup := router.Group("/favorite")
	favoriteRouteGroup.GET("/", fh.FindAllFavorites)
	favoriteRouteGroup.GET("/:userId", fh.FindFavoritesByUser)
	favoriteRouteGroup.POST("/", fh.CreateNewFavorite)

	// set up user endpoints
	ur, _ := store.NewUserStoreDDB()
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
