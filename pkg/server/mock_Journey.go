// Code generated by mockery v1.0.0. DO NOT EDIT.

package server

import discovery "bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
import executors "bitbucket.org/openbankingteam/conformance-suite/pkg/executors"
import generation "bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
import mock "github.com/stretchr/testify/mock"

// MockJourney is an autogenerated mock type for the Journey type
type MockJourney struct {
	mock.Mock
}

// AllTokenCollected provides a mock function with given fields:
func (_m *MockJourney) AllTokenCollected() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// CollectToken provides a mock function with given fields: code, state, scope
func (_m *MockJourney) CollectToken(code string, state string, scope string) error {
	ret := _m.Called(code, state, scope)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(code, state, scope)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Results provides a mock function with given fields:
func (_m *MockJourney) Results() executors.DaemonController {
	ret := _m.Called()

	var r0 executors.DaemonController
	if rf, ok := ret.Get(0).(func() executors.DaemonController); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(executors.DaemonController)
		}
	}

	return r0
}

// RunTests provides a mock function with given fields:
func (_m *MockJourney) RunTests() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetConfig provides a mock function with given fields: config
func (_m *MockJourney) SetConfig(config JourneyConfig) error {
	ret := _m.Called(config)

	var r0 error
	if rf, ok := ret.Get(0).(func(JourneyConfig) error); ok {
		r0 = rf(config)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetDiscoveryModel provides a mock function with given fields: discoveryModel
func (_m *MockJourney) SetDiscoveryModel(discoveryModel *discovery.Model) (discovery.ValidationFailures, error) {
	ret := _m.Called(discoveryModel)

	var r0 discovery.ValidationFailures
	if rf, ok := ret.Get(0).(func(*discovery.Model) discovery.ValidationFailures); ok {
		r0 = rf(discoveryModel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(discovery.ValidationFailures)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*discovery.Model) error); ok {
		r1 = rf(discoveryModel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StopTestRun provides a mock function with given fields:
func (_m *MockJourney) StopTestRun() {
	_m.Called()
}

// TestCases provides a mock function with given fields:
func (_m *MockJourney) TestCases() (generation.TestCasesRun, error) {
	ret := _m.Called()

	var r0 generation.TestCasesRun
	if rf, ok := ret.Get(0).(func() generation.TestCasesRun); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(generation.TestCasesRun)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}