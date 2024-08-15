// Code generated by mockery v2.44.1. DO NOT EDIT.

package mocks

import (
	context "context"
	domain "task-manager-api-clean/domain"

	mock "github.com/stretchr/testify/mock"
)

// UserUseCase is an autogenerated mock type for the UserUseCase type
type UserUseCase struct {
	mock.Mock
}

// Login provides a mock function with given fields: c, payload
func (_m *UserUseCase) Login(c context.Context, payload *domain.UserLogin) (string, error) {
	ret := _m.Called(c, payload)

	if len(ret) == 0 {
		panic("no return value specified for Login")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.UserLogin) (string, error)); ok {
		return rf(c, payload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *domain.UserLogin) string); ok {
		r0 = rf(c, payload)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *domain.UserLogin) error); ok {
		r1 = rf(c, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Promote provides a mock function with given fields: c, username
func (_m *UserUseCase) Promote(c context.Context, username string) (*domain.UserInfo, error) {
	ret := _m.Called(c, username)

	if len(ret) == 0 {
		panic("no return value specified for Promote")
	}

	var r0 *domain.UserInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*domain.UserInfo, error)); ok {
		return rf(c, username)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.UserInfo); ok {
		r0 = rf(c, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.UserInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(c, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterUser provides a mock function with given fields: c, payload
func (_m *UserUseCase) RegisterUser(c context.Context, payload *domain.UserCreate) (*domain.UserInfo, error) {
	ret := _m.Called(c, payload)

	if len(ret) == 0 {
		panic("no return value specified for RegisterUser")
	}

	var r0 *domain.UserInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.UserCreate) (*domain.UserInfo, error)); ok {
		return rf(c, payload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *domain.UserCreate) *domain.UserInfo); ok {
		r0 = rf(c, payload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.UserInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *domain.UserCreate) error); ok {
		r1 = rf(c, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserUseCase creates a new instance of UserUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserUseCase {
	mock := &UserUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
