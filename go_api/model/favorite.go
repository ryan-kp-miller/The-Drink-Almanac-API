package model

type Favorite struct {
	Id      string `dynamodbav:"id"`
	UserId  string `dynamodbav:"user_id"`
	DrinkId string `dynamodbav:"drink_id"`
}
