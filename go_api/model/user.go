package model

type User struct {
	Id       string `dynamodbav:"id"`
	Username string `dynamodbav:"username"`
	Password string `dynamodbav:"password"`
}
