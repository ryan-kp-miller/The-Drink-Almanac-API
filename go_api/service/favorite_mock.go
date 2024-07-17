// Code generated by mockery v2.36.0. DO NOT EDIT.

package service

import (
	model "the-drink-almanac-api/model"

	mock "github.com/stretchr/testify/mock"
)

// MockFavoriteService is an autogenerated mock type for the FavoriteService type
type MockFavoriteService struct {
	mock.Mock
}

// CreateNewFavorite provides a mock function with given fields: userId, drinkId
func (_m *MockFavoriteService) CreateNewFavorite(userId string, drinkId string) (*model.Favorite, error) {
	ret := _m.Called(userId, drinkId)

	var r0 *model.Favorite
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (*model.Favorite, error)); ok {
		return rf(userId, drinkId)
	}
	if rf, ok := ret.Get(0).(func(string, string) *model.Favorite); ok {
		r0 = rf(userId, drinkId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Favorite)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(userId, drinkId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteFavorite provides a mock function with given fields: id
func (_m *MockFavoriteService) DeleteFavorite(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAllFavorites provides a mock function with given fields:
func (_m *MockFavoriteService) FindAllFavorites() ([]model.Favorite, error) {
	ret := _m.Called()

	var r0 []model.Favorite
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]model.Favorite, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []model.Favorite); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Favorite)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindFavoritesByUser provides a mock function with given fields: userId
func (_m *MockFavoriteService) FindFavoritesByUser(userId string) ([]model.Favorite, error) {
	ret := _m.Called(userId)

	var r0 []model.Favorite
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]model.Favorite, error)); ok {
		return rf(userId)
	}
	if rf, ok := ret.Get(0).(func(string) []model.Favorite); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Favorite)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockFavoriteService creates a new instance of MockFavoriteService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockFavoriteService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockFavoriteService {
	mock := &MockFavoriteService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
