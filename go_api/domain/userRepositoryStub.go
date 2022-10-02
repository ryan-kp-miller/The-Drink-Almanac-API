package domain

type UserRepositoryStub struct {
	users []User
}

func (s UserRepositoryStub) FindAll() ([]User, error) {
	return s.users, nil
}

func NewUserRepositoryStub() (UserRepositoryStub, error) {
	users := []User{
		{
			Id:       0,
			Username: "test0",
			Password: "test0",
		},
		{
			Id:       1,
			Username: "test1",
			Password: "test1",
		},
		{
			Id:       2,
			Username: "test2",
			Password: "test2",
		},
	}
	return UserRepositoryStub{users: users}, nil
}
