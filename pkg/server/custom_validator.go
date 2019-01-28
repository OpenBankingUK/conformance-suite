package server

import "gopkg.in/go-playground/validator.v9"

// CustomValidator used to validate incoming payloads (for now).
// https://echo.labstack.com/guide/request#validate-data
type CustomValidator struct {
	validator *validator.Validate
}

func newCustomValidator() *CustomValidator {
	return &CustomValidator{
		validator.New(),
	}
}

// Validate incoming payloads (for now) that contain the struct tag `validate:"required"`.
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
