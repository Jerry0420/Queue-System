// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/jerry0420/queue-system/backend/domain"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// IntegrationUseCaseInterface is an autogenerated mock type for the IntegrationUseCaseInterface type
type IntegrationUseCaseInterface struct {
	mock.Mock
}

// CloseStore provides a mock function with given fields: ctx, store
func (_m *IntegrationUseCaseInterface) CloseStore(ctx context.Context, store domain.Store) error {
	ret := _m.Called(ctx, store)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Store) error); ok {
		r0 = rf(ctx, store)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CloseStoreRoutine provides a mock function with given fields: ctx
func (_m *IntegrationUseCaseInterface) CloseStoreRoutine(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateCustomers provides a mock function with given fields: ctx, session, oldStatus, newStatus, customers
func (_m *IntegrationUseCaseInterface) CreateCustomers(ctx context.Context, session *domain.StoreSession, oldStatus string, newStatus string, customers []domain.Customer) error {
	ret := _m.Called(ctx, session, oldStatus, newStatus, customers)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.StoreSession, string, string, []domain.Customer) error); ok {
		r0 = rf(ctx, session, oldStatus, newStatus, customers)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateStore provides a mock function with given fields: ctx, store, queues
func (_m *IntegrationUseCaseInterface) CreateStore(ctx context.Context, store *domain.Store, queues []domain.Queue) error {
	ret := _m.Called(ctx, store, queues)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Store, []domain.Queue) error); ok {
		r0 = rf(ctx, store, queues)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ForgetPassword provides a mock function with given fields: ctx, email
func (_m *IntegrationUseCaseInterface) ForgetPassword(ctx context.Context, email string) (domain.Store, error) {
	ret := _m.Called(ctx, email)

	var r0 domain.Store
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.Store); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(domain.Store)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStoreWithQueuesAndCustomersById provides a mock function with given fields: ctx, storeId
func (_m *IntegrationUseCaseInterface) GetStoreWithQueuesAndCustomersById(ctx context.Context, storeId int) (domain.StoreWithQueues, error) {
	ret := _m.Called(ctx, storeId)

	var r0 domain.StoreWithQueues
	if rf, ok := ret.Get(0).(func(context.Context, int) domain.StoreWithQueues); ok {
		r0 = rf(ctx, storeId)
	} else {
		r0 = ret.Get(0).(domain.StoreWithQueues)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, storeId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RefreshToken provides a mock function with given fields: ctx, encryptedRefreshToken
func (_m *IntegrationUseCaseInterface) RefreshToken(ctx context.Context, encryptedRefreshToken string) (domain.Store, string, string, time.Time, error) {
	ret := _m.Called(ctx, encryptedRefreshToken)

	var r0 domain.Store
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.Store); ok {
		r0 = rf(ctx, encryptedRefreshToken)
	} else {
		r0 = ret.Get(0).(domain.Store)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(context.Context, string) string); ok {
		r1 = rf(ctx, encryptedRefreshToken)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 string
	if rf, ok := ret.Get(2).(func(context.Context, string) string); ok {
		r2 = rf(ctx, encryptedRefreshToken)
	} else {
		r2 = ret.Get(2).(string)
	}

	var r3 time.Time
	if rf, ok := ret.Get(3).(func(context.Context, string) time.Time); ok {
		r3 = rf(ctx, encryptedRefreshToken)
	} else {
		r3 = ret.Get(3).(time.Time)
	}

	var r4 error
	if rf, ok := ret.Get(4).(func(context.Context, string) error); ok {
		r4 = rf(ctx, encryptedRefreshToken)
	} else {
		r4 = ret.Error(4)
	}

	return r0, r1, r2, r3, r4
}

// SigninStore provides a mock function with given fields: ctx, email, password
func (_m *IntegrationUseCaseInterface) SigninStore(ctx context.Context, email string, password string) (domain.Store, string, time.Time, error) {
	ret := _m.Called(ctx, email, password)

	var r0 domain.Store
	if rf, ok := ret.Get(0).(func(context.Context, string, string) domain.Store); ok {
		r0 = rf(ctx, email, password)
	} else {
		r0 = ret.Get(0).(domain.Store)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(context.Context, string, string) string); ok {
		r1 = rf(ctx, email, password)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 time.Time
	if rf, ok := ret.Get(2).(func(context.Context, string, string) time.Time); ok {
		r2 = rf(ctx, email, password)
	} else {
		r2 = ret.Get(2).(time.Time)
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(context.Context, string, string) error); ok {
		r3 = rf(ctx, email, password)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// UpdatePassword provides a mock function with given fields: ctx, passwordToken, newPassword
func (_m *IntegrationUseCaseInterface) UpdatePassword(ctx context.Context, passwordToken string, newPassword string) (domain.Store, error) {
	ret := _m.Called(ctx, passwordToken, newPassword)

	var r0 domain.Store
	if rf, ok := ret.Get(0).(func(context.Context, string, string) domain.Store); ok {
		r0 = rf(ctx, passwordToken, newPassword)
	} else {
		r0 = ret.Get(0).(domain.Store)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, passwordToken, newPassword)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VerifyNormalToken provides a mock function with given fields: ctx, normalToken
func (_m *IntegrationUseCaseInterface) VerifyNormalToken(ctx context.Context, normalToken string) (domain.TokenClaims, error) {
	ret := _m.Called(ctx, normalToken)

	var r0 domain.TokenClaims
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.TokenClaims); ok {
		r0 = rf(ctx, normalToken)
	} else {
		r0 = ret.Get(0).(domain.TokenClaims)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, normalToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VerifySessionToken provides a mock function with given fields: ctx, sessionToken
func (_m *IntegrationUseCaseInterface) VerifySessionToken(ctx context.Context, sessionToken string) (domain.Store, error) {
	ret := _m.Called(ctx, sessionToken)

	var r0 domain.Store
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.Store); ok {
		r0 = rf(ctx, sessionToken)
	} else {
		r0 = ret.Get(0).(domain.Store)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, sessionToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
