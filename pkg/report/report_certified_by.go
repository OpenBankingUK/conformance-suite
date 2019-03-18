package report

import validation "github.com/go-ozzo/ozzo-validation"

// ReportCertifiedBy - contains details of who certified the `Report`.
type ReportCertifiedBy struct {
	Environment  ReportCertifiedByEnvironment `json:"environment"`  // Name of the environment tested
	Brand        string                       `json:"brand"`        // Name of the brand tested
	AuthorisedBy string                       `json:"authorisedBy"` // Name of the Authoriser
	JobTitle     string                       `json:"jobTitle"`     // Job Title of the Authorisee
}

// Validate - called by `github.com/go-ozzo/ozzo-validation` to validate struct.
func (r ReportCertifiedBy) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Environment, validation.Required, validation.In(
			ReportCertifiedByEnvironmentTesting,
			ReportCertifiedByEnvironmentProduction,
		)),
		validation.Field(&r.Brand, validation.Required, validation.Length(1, 60)),
		validation.Field(&r.AuthorisedBy, validation.Required, validation.Length(1, 60)),
		validation.Field(&r.JobTitle, validation.Required, validation.Length(1, 60)),
	)
}
