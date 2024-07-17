// Code generated by mockery v2.36.0. DO NOT EDIT.

package repository

import (
	model "the-drink-almanac-api/model"

	mock "github.com/stretchr/testify/mock"
)

// MockUserRepository is an autogenerated mock type for the UserRepository type
type MockUserRepository struct {
	mock.Mock
}

// CreateNewUser provides a mock function with given fields: _a0
func (_m *MockUserRepository) CreateNewUser(_a0 model.User) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(model.User) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteUser provides a mock function with given fields: id
func (_m *MockUserRepository) DeleteUser(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAll provides a mock function with given fields:
func (_m *MockUserRepository) FindAll() ([]model.User, error) {
	ret := _m.Called()

	var r0 []model.User
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]model.User, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []model.User); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.User)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindUserById provides a mock function with given fields: userId
func (_m *MockUserRepository) FindUserById(userId string) (*model.User, error) {
	ret := _m.Called(userId)

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.User, error)); ok {
		return rf(userId)
	}
	if rf, ok := ret.Get(0).(func(string) *model.User); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindUserByUsername provides a mock function with given fields: username
func (_m *MockUserRepository) FindUserByUsername(username string) (*model.User, error) {
	ret := _m.Called(username)

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.User, error)); ok {
		return rf(username)
	}
	if rf, ok := ret.Get(0).(func(string) *model.User); ok {
		r0 = rf(username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockUserRepository creates a new instance of MockUserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUserRepository {
	mock := &MockUserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
