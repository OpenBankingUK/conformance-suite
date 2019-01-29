// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import io "io"
import mock "github.com/stretchr/testify/mock"

// Generator is an autogenerated mock type for the Generator type
type Generator struct {
	mock.Mock
}

// Generate provides a mock function with given fields: input, output
func (_m *Generator) Generate(input io.Reader, output io.Writer) error {
	ret := _m.Called(input, output)

	var r0 error
	if rf, ok := ret.Get(0).(func(io.Reader, io.Writer) error); ok {
		r0 = rf(input, output)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}