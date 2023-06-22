package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	lambdaHandler "the-drink-almanac-api/handler/lambda"
	"the-drink-almanac-api/model"
	"the-drink-almanac-api/service"
	"the-drink-almanac-api/store"
)

func start(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("starting favorites lambda")
	appConfig := model.NewAppConfig()
	authService := service.NewJwtAuthService(appConfig.JwtSecretKey)
	favoriteStore, _ := store.NewFavoriteStoreDDB(appConfig.FavoritesTableName, appConfig.AwsEndpoint)
	favoriteService := service.NewDefaultFavoriteService(favoriteStore)
	favoriteHandler := lambdaHandler.NewFavoritesLambdaHandler(favoriteService, authService)

	return favoriteHandler.RouteRequest(request)
}

func main() {
	lambda.Start(start)
}
