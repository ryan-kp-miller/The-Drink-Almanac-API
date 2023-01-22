package main

import (
	"fmt"
	"net/http"
	"the-drink-almanac-api/handler"
	"the-drink-almanac-api/service"
	"the-drink-almanac-api/store"

	"github.com/gin-gonic/gin"
)

func Start(port string) {
	router := gin.Default()

	// set up default endpoint
	router.GET("", hello_world_handler)

	// set up favorite endpoints
	favoriteStore, _ := store.NewFavoriteStoreDDB()
	favoriteService := service.NewDefaultFavoriteService(favoriteStore)
	favoriteHandlers := handler.FavoriteHandlers{Service: favoriteService}
	favoriteRouteGroup := router.Group("/favorite")
	favoriteRouteGroup.GET("/", favoriteHandlers.FindAllFavorites)
	favoriteRouteGroup.GET("/:userId", favoriteHandlers.FindFavoritesByUser)
	favoriteRouteGroup.POST("/", favoriteHandlers.CreateNewFavorite)
	favoriteRouteGroup.DELETE("/:favoriteId", favoriteHandlers.DeleteFavorite)

	// set up user endpoints
	userStore, _ := store.NewUserStoreDDB()
	authService := service.NewJwtAuthService()
	userService := service.NewDefaultUserService(userStore, authService)
	userHandlers := handler.UserHandlers{Service: userService}
	userRouteGroup := router.Group("/user")
	userRouteGroup.GET("", userHandlers.FindAllUsers)
	userRouteGroup.POST("", userHandlers.CreateNewUser)
	userRouteGroup.DELETE("", userHandlers.DeleteUser)
	userRouteGroup.POST("/login", userHandlers.Login)

	// running the app
	router.Run(fmt.Sprintf(":%s", port))
}

func hello_world_handler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World!",
	})
}
