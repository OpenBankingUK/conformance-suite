package report

import validation "github.com/go-ozzo/ozzo-validation"

// CertifiedBy - contains details of who certified the `Report`.
type CertifiedBy struct {
	Environment  CertifiedByEnvironment `json:"environment"`  // Name of the environment tested
	Brand        string                 `json:"brand"`        // Name of the brand tested
	AuthorisedBy string                 `json:"authorisedBy"` // Name of the Authoriser
	JobTitle     string                 `json:"jobTitle"`     // Job Title of the Authorisee
}

// Validate - called by `github.com/go-ozzo/ozzo-validation` to validate struct.
func (r CertifiedBy) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Environment, validation.Required, validation.In(
			CertifiedByEnvironmentTesting,
			CertifiedByEnvironmentProduction,
			CertifiedByEnvironmentSandbox,
		)),
		validation.Field(&r.Brand, validation.Required, validation.Length(1, 60)),
		validation.Field(&r.AuthorisedBy, validation.Required, validation.Length(1, 60)),
		validation.Field(&r.JobTitle, validation.Required, validation.Length(1, 60)),
	)
}
