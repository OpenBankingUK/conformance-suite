// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import results "bitbucket.org/openbankingteam/conformance-suite/pkg/executors/results"

// DaemonController is an autogenerated mock type for the DaemonController type
type DaemonController struct {
	mock.Mock
}

// Errors provides a mock function with given fields:
func (_m *DaemonController) Errors() chan error {
	ret := _m.Called()

	var r0 chan error
	if rf, ok := ret.Get(0).(func() chan error); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan error)
		}
	}

	return r0
}

// Results provides a mock function with given fields:
func (_m *DaemonController) Results() chan results.TestCase {
	ret := _m.Called()

	var r0 chan results.TestCase
	if rf, ok := ret.Get(0).(func() chan results.TestCase); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan results.TestCase)
		}
	}

	return r0
}

// ShouldStop provides a mock function with given fields:
func (_m *DaemonController) ShouldStop() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Stop provides a mock function with given fields:
func (_m *DaemonController) Stop() {
	_m.Called()
}