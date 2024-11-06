package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// AuthRequest - Request to `/api/login` POST.
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate - used by github.com/go-ozzo/ozzo-validation to validate struct.
func (r AuthRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required),
		validation.Field(&r.Password, validation.Required),
	)
}
