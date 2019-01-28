package server

import "gopkg.in/go-playground/validator.v9"

// echoValidatorAdapter adapts an go-playground Validator to a echo Validator
type echoValidatorAdapter struct {
	validator *validator.Validate
}

func newEchoValidatorAdapter() *echoValidatorAdapter {
	return &echoValidatorAdapter{
		validator.New(),
	}
}

func (cv *echoValidatorAdapter) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
