package models

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/events"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/results"
	validation "github.com/go-ozzo/ozzo-validation"
)

// ExportRequest - Request to `/api/export`.
type ExportRequest struct {
	Implementer         string `json:"implementer"`           // Implementer/Brand Name
	AuthorisedBy        string `json:"authorised_by"`         // Authorised by
	JobTitle            string `json:"job_title"`             // Job Title
	HasAgreed           bool   `json:"has_agreed"`            // I agree
	AddDigitalSignature bool   `json:"add_digital_signature"` // Sign this report
}

func (e ExportRequest) Validate() error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.Implementer, validation.Required),
		validation.Field(&e.AuthorisedBy, validation.Required),
		validation.Field(&e.JobTitle, validation.Required),
		validation.Field(&e.HasAgreed, validation.Required, validation.In(true)),
	)
}

// ExportResults - Contains `ExportRequest` and results of test run.
type ExportResults struct {
	ExportRequest ExportRequest                `json:"export_request"`
	HasPassed     bool                         `json:"has_passed"`
	Results       []results.TestCase           `json:"results"`
	Tokens        []events.AcquiredAccessToken `json:"tokens"`
}
