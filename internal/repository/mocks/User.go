// Code generated by mockery v2.7.5. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	entity "gitlab.com/nodefluxio/vanilla-dashboard/internal/entity"

	util "gitlab.com/nodefluxio/vanilla-dashboard/internal/pkg/util"
)

// User is an autogenerated mock type for the User type
type User struct {
	mock.Mock
}

// Count provides a mock function with given fields: ctx, paging
func (_m *User) Count(ctx context.Context, paging *util.Pagination) (int, error) {
	ret := _m.Called(ctx, paging)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, *util.Pagination) int); ok {
		r0 = rf(ctx, paging)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *util.Pagination) error); ok {
		r1 = rf(ctx, paging)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: ctx, data
func (_m *User) Create(ctx context.Context, data *entity.User) (*entity.User, error) {
	ret := _m.Called(ctx, data)

	var r0 *entity.User
	if rf, ok := ret.Get(0).(func(context.Context, *entity.User) *entity.User); ok {
		r0 = rf(ctx, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *entity.User) error); ok {
		r1 = rf(ctx, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, ID
func (_m *User) Delete(ctx context.Context, ID uint64) error {
	ret := _m.Called(ctx, ID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) error); ok {
		r0 = rf(ctx, ID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByEmail provides a mock function with given fields: ctx, email
func (_m *User) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	ret := _m.Called(ctx, email)

	var r0 *entity.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *entity.User); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByUsername provides a mock function with given fields: ctx, username
func (_m *User) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	ret := _m.Called(ctx, username)

	var r0 *entity.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *entity.User); ok {
		r0 = rf(ctx, username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDetail provides a mock function with given fields: ctx, ID
func (_m *User) GetDetail(ctx context.Context, ID uint64) (*entity.User, error) {
	ret := _m.Called(ctx, ID)

	var r0 *entity.User
	if rf, ok := ret.Get(0).(func(context.Context, uint64) *entity.User); ok {
		r0 = rf(ctx, ID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetList provides a mock function with given fields: ctx, paging
func (_m *User) GetList(ctx context.Context, paging *util.Pagination) ([]*entity.User, error) {
	ret := _m.Called(ctx, paging)

	var r0 []*entity.User
	if rf, ok := ret.Get(0).(func(context.Context, *util.Pagination) []*entity.User); ok {
		r0 = rf(ctx, paging)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entity.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *util.Pagination) error); ok {
		r1 = rf(ctx, paging)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateBasicData provides a mock function with given fields: ctx, data
func (_m *User) UpdateBasicData(ctx context.Context, data *entity.User) error {
	ret := _m.Called(ctx, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.User) error); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdatePassword provides a mock function with given fields: ctx, password, ID
func (_m *User) UpdatePassword(ctx context.Context, password string, ID uint64) error {
	ret := _m.Called(ctx, password, ID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, uint64) error); ok {
		r0 = rf(ctx, password, ID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}