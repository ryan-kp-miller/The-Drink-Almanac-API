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

func start(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	fmt.Println("starting users lambda")
	appConfig := model.NewAppConfig()
	authService := service.NewJwtAuthService(appConfig.JwtSecretKey)
	userStore, _ := store.NewUserStoreDDB(appConfig.UsersTableName, appConfig.AwsEndpoint)
	userService := service.NewDefaultUserService(userStore)
	userHandler := lambdaHandler.NewUsersLambdaHandler(userService, authService)

	return userHandler.RouteRequest(request)
}

func main() {
	lambda.Start(start)
}
