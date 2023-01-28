package main

import (
	"fmt"
	"net/http"
	"the-drink-almanac-api/handler"
	"the-drink-almanac-api/model"
	"the-drink-almanac-api/service"
	"the-drink-almanac-api/store"

	"github.com/gin-gonic/gin"
)

func Start(port string) {
	appConfig := model.NewAppConfig()
	router := gin.Default()

	// set up default endpoint
	router.GET("", hello_world_handler)

	// set up favorite endpoints
	favoriteStore, _ := store.NewFavoriteStoreDDB(appConfig.FavoritesTableName, appConfig.AwsEndpoint)
	favoriteService := service.NewDefaultFavoriteService(favoriteStore)
	favoriteHandlers := handler.FavoriteHandlers{Service: favoriteService}
	favoriteRouteGroup := router.Group("/favorite")
	favoriteRouteGroup.GET("", favoriteHandlers.FindFavoritesByUser)
	favoriteRouteGroup.POST("", favoriteHandlers.CreateNewFavorite)
	favoriteRouteGroup.DELETE("/:favoriteId", favoriteHandlers.DeleteFavorite)

	// set up user endpoints
	userStore, _ := store.NewUserStoreDDB(appConfig.UsersTableName, appConfig.AwsEndpoint)
	userService := service.NewDefaultUserService(userStore)
	authService := service.NewJwtAuthService(appConfig.JwtSecretKey)
	userHandler := handler.NewUserHandler(userService, authService)
	userRouteGroup := router.Group("/user")
	userRouteGroup.GET("", userHandler.FindUser)
	userRouteGroup.POST("", userHandler.CreateNewUser)
	userRouteGroup.DELETE("", userHandler.DeleteUser)
	userRouteGroup.POST("/login", userHandler.Login)

	// running the app
	router.Run(fmt.Sprintf(":%s", port))
}

func hello_world_handler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World!",
	})
}
