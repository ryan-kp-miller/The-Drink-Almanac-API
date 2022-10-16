package store

import "the-drink-almanac-api/model"

type UserStoreStub struct {
	users []model.User
}

func (s UserStoreStub) FindAll() ([]model.User, error) {
	return s.users, nil
}

func NewUserStoreStub() (UserStoreStub, error) {
	users := []model.User{
		{
			Id:       "0",
			Username: "test0",
			Password: "test0",
		},
		{
			Id:       "1",
			Username: "test1",
			Password: "test1",
		},
		{
			Id:       "2",
			Username: "test2",
			Password: "test2",
		},
	}
	return UserStoreStub{users: users}, nil
}
