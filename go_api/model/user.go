package model

type User struct {
	Id       string `json:"id" dynamodbav:"id"`
	Username string `json:"username" dynamodbav:"username"`
	Password string `json:"password" dynamodbav:"password"`
}

type UserStore interface {
	FindAll() ([]User, error)
	FindUserByUsername(username string) (*User, error)
	CreateNewUser(User) error
}
