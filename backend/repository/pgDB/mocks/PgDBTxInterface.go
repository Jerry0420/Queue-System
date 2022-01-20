// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	pgDB "github.com/jerry0420/queue-system/backend/repository/pgDB"
	mock "github.com/stretchr/testify/mock"
)

// PgDBTxInterface is an autogenerated mock type for the PgDBTxInterface type
type PgDBTxInterface struct {
	mock.Mock
}

// BeginTx provides a mock function with given fields:
func (_m *PgDBTxInterface) BeginTx() (pgDB.PgDBInterface, error) {
	ret := _m.Called()

	var r0 pgDB.PgDBInterface
	if rf, ok := ret.Get(0).(func() pgDB.PgDBInterface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pgDB.PgDBInterface)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CommitTx provides a mock function with given fields: pgDbTx
func (_m *PgDBTxInterface) CommitTx(pgDbTx pgDB.PgDBInterface) error {
	ret := _m.Called(pgDbTx)

	var r0 error
	if rf, ok := ret.Get(0).(func(pgDB.PgDBInterface) error); ok {
		r0 = rf(pgDbTx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RollbackTx provides a mock function with given fields: pgDbTx
func (_m *PgDBTxInterface) RollbackTx(pgDbTx pgDB.PgDBInterface) {
	_m.Called(pgDbTx)
}