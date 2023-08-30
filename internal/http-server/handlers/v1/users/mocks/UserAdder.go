// Code generated by mockery v2.33.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/AlexZahvatkin/segments-users-service/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// UserAdder is an autogenerated mock type for the UserAdder type
type UserAdder struct {
	mock.Mock
}

// AddUser provides a mock function with given fields: _a0, _a1
func (_m *UserAdder) AddUser(_a0 context.Context, _a1 string) (models.User, error) {
	ret := _m.Called(_a0, _a1)

	var r0 models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (models.User, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) models.User); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserAdder creates a new instance of UserAdder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserAdder(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserAdder {
	mock := &UserAdder{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}