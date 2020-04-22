// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Storer is an autogenerated mock type for the Storer type
type Storer struct {
	mock.Mock
}

// Do provides a mock function with given fields: commandName, args
func (_m *Storer) Do(commandName string, args ...interface{}) (interface{}, error) {
	var _ca []interface{}
	_ca = append(_ca, commandName)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(string, ...interface{}) interface{}); ok {
		r0 = rf(commandName, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, ...interface{}) error); ok {
		r1 = rf(commandName, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
