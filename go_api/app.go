package main

import (
	"fmt"
	"net/http"

	"the-drink-almanac-api/handler/middleware"
	"the-drink-almanac-api/handler/server"
	"the-drink-almanac-api/model"
	"the-drink-almanac-api/repository"
	"the-drink-almanac-api/service"

	"github.com/gin-gonic/gin"
)

func Start(port string) {
	appConfig := model.NewAppConfig()
	router := gin.Default()

	// set up default endpoint
	router.GET("", hello_world_handler)

	// set up auth middleware
	authService := service.NewJwtAuthService(appConfig.JwtSecretKey)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// set up favorite endpoints
	favoriteStore, _ := repository.NewFavoriteRepository(appConfig.FavoritesTableName, appConfig.AwsEndpoint)
	favoriteService := service.NewDefaultFavoriteService(favoriteStore)
	favoriteHandler := server.FavoriteHandler{Service: favoriteService}
	favoriteRouteGroup := router.Group("/favorite")
	favoriteRouteGroup.GET("", authMiddleware.AuthUser, favoriteHandler.FindFavoritesByUser)
	favoriteRouteGroup.POST("", authMiddleware.AuthUser, favoriteHandler.CreateNewFavorite)
	favoriteRouteGroup.DELETE("/:favoriteId", authMiddleware.AuthUser, favoriteHandler.DeleteFavorite)

	// set up user endpoints
	userStore, _ := repository.NewUserRepository(appConfig.UsersTableName, appConfig.AwsEndpoint)
	userService := service.NewDefaultUserService(userStore)
	userHandler := server.NewUserHandler(userService, authService)
	userRouteGroup := router.Group("/user")
	userRouteGroup.GET("", authMiddleware.AuthUser, userHandler.FindUser)
	userRouteGroup.POST("", userHandler.CreateNewUser)
	userRouteGroup.DELETE("", authMiddleware.AuthUser, userHandler.DeleteUser)
	userRouteGroup.POST("/login", userHandler.Login)

	// running the app
	router.Run(fmt.Sprintf(":%s", port))
}

func hello_world_handler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World!",
	})
}
