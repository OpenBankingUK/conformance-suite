package report

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/results"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/server/models"
	internal_time "bitbucket.org/openbankingteam/conformance-suite/pkg/time"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"time"
)

const (
	// Version - version of the `Report`.
	// TODO(mbana): probably need a better way of getting/setting the version of the Report Exporter
	Version = "0.0.1"
)

// Report - The Report.
type Report struct {
	ID               string             `json:"id"`                       // A unique and immutable identifier used to identify the report. The v4 UUIDs generated conform to RFC 4122.
	Created          string             `json:"created"`                  // Date and time when the report was created, formatted accorrding to RFC3339 (https://tools.ietf.org/html/rfc3339). Note RFC3339 is derived from ISO 8601 (https://en.wikipedia.org/wiki/ISO_8601).
	Expiration       *string            `json:"expiration,omitempty"`     // Date and time when the report should not longer be accepted, formatted accorrding to RFC3339 (https://tools.ietf.org/html/rfc3339). Note RFC3339 is derived from ISO 8601 (https://en.wikipedia.org/wiki/ISO_8601).
	Version          string             `json:"version"`                  // The current version of the report model used.
	Status           Status             `json:"status"`                   // A status describing overall condition of the report.
	CertifiedBy      CertifiedBy        `json:"certifiedBy"`              // The certifier of the report.
	SignatureChain   *[]SignatureChain  `json:"signatureChain,omitempty"` // When Add digital signature is set this contains the signature chain.
	Discovery        discovery.Model    `json:"-"`                        // Original used discovery model
	APISpecification []APISpecification `json:"apiSpecification"`         // API and version tested, along with test cases
}

type APISpecification struct {
	Name    string             `json:"name"`
	Version string             `json:"version"`
	Results []results.TestCase `json:"results"`
}

// Validate - called by `github.com/go-ozzo/ozzo-validation` to validate struct.
func (r Report) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUIDv4),
		validation.Field(&r.Created, validation.Required, validation.Date(internal_time.Layout)),
		validation.Field(&r.Expiration, validation.Date(internal_time.Layout)),
		validation.Field(&r.Version, validation.Required),
		validation.Field(&r.Status, validation.Required, validation.In(
			StatusPending,
			StatusComplete,
			StatusError,
		)),
		validation.Field(&r.CertifiedBy, validation.Required),
	)
}

// NewReport - create `Report` from `ExportResults`.
func NewReport(exportResults models.ExportResults) (Report, error) {
	// Random (Version 4) UUID. NB: `uuid.New()` might panic hence we using this function instead.
	uuid, err := uuid.NewRandom()
	if err != nil {
		return Report{}, err
	}

	created := time.Now().Format(internal_time.Layout)
	expiration := time.Now().AddDate(0, 3, 0).Format(internal_time.Layout) // Expires three (3) months from now
	certifiedBy := CertifiedBy{
		Environment:  CertifiedByEnvironmentTesting, // Hardcode to "testing" for now
		Brand:        exportResults.ExportRequest.Implementer,
		AuthorisedBy: exportResults.ExportRequest.AuthorisedBy,
		JobTitle:     exportResults.ExportRequest.JobTitle,
	}
	signatureChain := []SignatureChain{}

	var apiSpecs []APISpecification

	for k, v := range exportResults.Results {
		apiSpec := APISpecification{
			Name:    k.APIName,
			Version: k.APIVersion,
			Results: v,
		}
		apiSpecs = append(apiSpecs, apiSpec)
	}

	return Report{
		ID:               uuid.String(),
		Created:          created,
		Expiration:       &expiration,
		Version:          Version,
		Status:           StatusComplete,
		CertifiedBy:      certifiedBy,
		SignatureChain:   &signatureChain,
		Discovery:        exportResults.DiscoveryModel,
		APISpecification: apiSpecs,
	}, nil
}
