package report

import (
	internal_time "bitbucket.org/openbankingteam/conformance-suite/internal/pkg/time"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/server"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

const (
	// Version - version of the `Report`.
	// TODO(mbana): probably need a better way of getting/setting the version of the Report Exporter
	Version = "0.1.0"
)

// Report - The Report.
type Report struct {
	ID             string                  `json:"id"`                               // A unique and immutable identifier used to identify the report. The v4 UUIDs generated conform to RFC 4122.
	Created        string                  `json:"created"`                          // Date and time when the report was created, formatted accorrding to RFC3339 (https://tools.ietf.org/html/rfc3339). Note RFC3339 is derived from ISO 8601 (https://en.wikipedia.org/wiki/ISO_8601).
	Expiration     string                  `json:"expiration,omitempty"`             // Date and time when the report should not longer be accepted, formatted accorrding to RFC3339 (https://tools.ietf.org/html/rfc3339). Note RFC3339 is derived from ISO 8601 (https://en.wikipedia.org/wiki/ISO_8601).
	Version        string                  `json:"version"`                          // The current version of the report model used.
	Status         ReportStatus            `json:"status" validate:"required,max=8"` // A status describing overall condition of the report.
	CertifiedBy    ReportCertifiedBy       `json:"certifiedBy"`                      // The certifier of the report.
	SignatureChain *[]ReportSignatureChain `json:"signatureChain,omitempty"`         // When Add digital signature is set this contains the signature chain.
}

// Validate - called by `github.com/go-ozzo/ozzo-validation` to validate struct.
func (r Report) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUIDv4),
		validation.Field(&r.Created, validation.Required, validation.Date(internal_time.Layout)),
		validation.Field(&r.Expiration, validation.Date(internal_time.Layout)),
		validation.Field(&r.Version, validation.Required),
		validation.Field(&r.Status, validation.Required, validation.In(
			ReportStatusPending,
			ReportStatusComplete,
			ReportStatusError,
		)),
		validation.Field(&r.CertifiedBy, validation.Required),
	)
}

// NewReport - create `Report` from `ExportResponse`.
func NewReport(exportResponse server.ExportResponse) (Report, error) {
	return Report{}, nil
}
