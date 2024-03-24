package mocks

import (
	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"
)

type GormPostgres struct {
	mock.Mock
}

func (_m *GormPostgres) GetConnection() *gorm.DB {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetConnection")
	}

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func() *gorm.DB); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}

func NewGormPostgres(t interface {
	mock.TestingT
	Cleanup(func())
}) *GormPostgres {
	mock := &GormPostgres{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
