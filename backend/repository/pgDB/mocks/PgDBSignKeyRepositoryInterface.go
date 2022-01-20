// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/jerry0420/queue-system/backend/domain"
	mock "github.com/stretchr/testify/mock"
)

// PgDBSignKeyRepositoryInterface is an autogenerated mock type for the PgDBSignKeyRepositoryInterface type
type PgDBSignKeyRepositoryInterface struct {
	mock.Mock
}

// CreateSignKey provides a mock function with given fields: ctx, signKey
func (_m *PgDBSignKeyRepositoryInterface) CreateSignKey(ctx context.Context, signKey *domain.SignKey) error {
	ret := _m.Called(ctx, signKey)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.SignKey) error); ok {
		r0 = rf(ctx, signKey)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetSignKeyByID provides a mock function with given fields: ctx, id, signKeyType
func (_m *PgDBSignKeyRepositoryInterface) GetSignKeyByID(ctx context.Context, id int, signKeyType string) (domain.SignKey, error) {
	ret := _m.Called(ctx, id, signKeyType)

	var r0 domain.SignKey
	if rf, ok := ret.Get(0).(func(context.Context, int, string) domain.SignKey); ok {
		r0 = rf(ctx, id, signKeyType)
	} else {
		r0 = ret.Get(0).(domain.SignKey)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, string) error); ok {
		r1 = rf(ctx, id, signKeyType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveSignKeyByID provides a mock function with given fields: ctx, id, signKeyType
func (_m *PgDBSignKeyRepositoryInterface) RemoveSignKeyByID(ctx context.Context, id int, signKeyType string) (domain.SignKey, error) {
	ret := _m.Called(ctx, id, signKeyType)

	var r0 domain.SignKey
	if rf, ok := ret.Get(0).(func(context.Context, int, string) domain.SignKey); ok {
		r0 = rf(ctx, id, signKeyType)
	} else {
		r0 = ret.Get(0).(domain.SignKey)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, string) error); ok {
		r1 = rf(ctx, id, signKeyType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
