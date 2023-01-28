package dto

import "the-drink-almanac-api/model"

type UserResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

func NewUserResponse(user model.User) UserResponse {
	return UserResponse{
		Id:       user.Id,
		Username: user.Username,
	}
}

func NewUsersResponse(users []model.User) []UserResponse {
	usersResponse := make([]UserResponse, len(users))
	for i, user := range users {
		usersResponse[i] = NewUserResponse(user)
	}
	return usersResponse
}
