package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	discovery_mocks "bitbucket.org/openbankingteam/conformance-suite/pkg/discovery/mocks"
	gmocks "bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/report"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/server/models"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
	internal_time "bitbucket.org/openbankingteam/conformance-suite/pkg/time"
	version_mocks "bitbucket.org/openbankingteam/conformance-suite/pkg/version/mocks"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

func TestServerPostExport(t *testing.T) {
	require := test.NewRequire(t)

	discoveryModel := &discovery.Model{}
	validator := &discovery_mocks.Validator{}
	validator.On("Validate", discoveryModel).Return(discovery.NoValidationFailures(), nil)
	generator := &gmocks.MockGenerator{}
	journey := NewJourney(nullLogger(), generator, validator, discovery.NewNullTLSValidator(), false)

	failures, err := journey.SetDiscoveryModel(discoveryModel)
	require.NoError(err)
	require.Equal(discovery.NoValidationFailures(), failures)

	server := NewServer(journey, nullLogger(), &version_mocks.Version{})
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	expected := models.ExportRequest{
		Environment:  "sandbox",
		Implementer:  "implementer",
		AuthorisedBy: "authorised_by",
		JobTitle:     "job_title",
		Products: []string{
			"Business",
			"Personal",
			"Cards",
		},
		HasAgreed:           true,
		AddDigitalSignature: false,
	}
	requestJSON, err := json.MarshalIndent(expected, ``, `  `)
	require.NoError(err)
	require.NotNil(requestJSON)

	// make the request and record time the export operation was started
	exportTime := time.Now()
	// we can remove this sleep if we change `bitbucket.org/openbankingteam/conformance-suite/pkg/time`.Layout to `time.RFC3339Nano`.
	time.Sleep(1 * time.Second)
	code, body, headers := request(http.MethodPost, "/api/export", bytes.NewReader(requestJSON), server)

	require.Equal(http.StatusOK, code, body.String())
	require.Equal(http.Header{
		echo.HeaderContentType: []string{
			MIMEApplicationZIP,
		},
		// echo.HeaderContentDisposition: []string{
		// 	`attachment; filename="report.zip"`,
		// },
	}, headers, body.String())

	// Verify returned `report.zip` is correct.
	require.NotNil(body)
	importer := report.NewZipImporter(body)
	actual, err := importer.Import()
	require.NoError(err)

	// `actual` contains someething like below
	// actual := report.Report{
	// 	ID:         "80581bb3-fdff-4f37-ba11-b688ecb20b73",
	// 	Created:    "2019-03-21T13:00:11Z",
	// 	Expiration: (*string)(0xc0012b6830),
	// 	Version:    "0.0.1",
	// 	Status:     2,
	// 	CertifiedBy: report.CertifiedBy{
	// 		Environment:  1,
	// 		Brand:        "implementer",
	// 		AuthorisedBy: "authorised_by",
	// 		JobTitle:     "job_title",
	// 	},
	// 	SignatureChain: (*[]report.SignatureChain)(0xc0012e07c0),
	// }

	// check uuids is version 4 and is RFC4122 variant
	id, err := uuid.Parse(actual.ID)
	require.NoError(err)
	require.Equal(uuid.RFC4122, id.Variant())
	require.Equal(uuid.Version(4), id.Version())

	// check created time is before or equal the time we called `/api/export/`
	reportCreatedTime, err := time.Parse(internal_time.Layout, actual.Created)
	require.NoError(err)
	require.True(reportCreatedTime.Sub(exportTime) >= 0)
	require.True(exportTime.Equal(reportCreatedTime) || reportCreatedTime.After(exportTime))

	require.Equal(report.StatusComplete, actual.Status)
	require.Equal(report.CertifiedByEnvironmentSandbox, actual.CertifiedBy.Environment)
	require.Equal(expected.Implementer, actual.CertifiedBy.Brand)
	require.Equal(expected.AuthorisedBy, actual.CertifiedBy.AuthorisedBy)
	require.Equal(expected.JobTitle, actual.CertifiedBy.JobTitle)
}

func TestServerPostExport_InvalidRequest(t *testing.T) {
	require := test.NewRequire(t)

	server := NewServer(testJourney(), nullLogger(), &version_mocks.Version{})
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	type TestCase struct {
		request models.ExportRequest
		err     string
	}
	testCases := []TestCase{
		{
			request: models.ExportRequest{
				Environment:         "sandbox",
				Implementer:         "implementer",
				AuthorisedBy:        "authorised_by",
				JobTitle:            "job_title",
				Products:            []string{},
				HasAgreed:           true,
				AddDigitalSignature: false,
			},
			err: `{"error":"products: cannot be blank."}`,
		},
		{
			request: models.ExportRequest{
				Environment:  "sandbox",
				Implementer:  "implementer",
				AuthorisedBy: "authorised_by",
				JobTitle:     "job_title",
				Products: []string{
					"Business",
					"Business",
				},
				HasAgreed:           true,
				AddDigitalSignature: false,
			},
			err: `{"error":"products: pkg/server/models.ExportRequest: 'products' ([\"Business\" \"Business\"]) contains duplicate value (\"Business\")."}`,
		},
		{
			request: models.ExportRequest{
				Environment:  "sandbox",
				Implementer:  "implementer",
				AuthorisedBy: "authorised_by",
				JobTitle:     "job_title",
				Products: []string{
					"Business",
					"Personal",
					"Cards",
					"More_Than_Allowed_Values",
				},
				HasAgreed:           true,
				AddDigitalSignature: false,
			},
			err: `{"error":"products: pkg/server/models.ExportRequest: 'products' (4) contains more than supported values (3)."}`,
		},
		{
			request: models.ExportRequest{
				Environment:  "sandbox",
				Implementer:  "implementer",
				AuthorisedBy: "authorised_by",
				JobTitle:     "job_title",
				Products: []string{
					"Invalid_Product",
				},
				HasAgreed:           true,
				AddDigitalSignature: false,
			},
			err: `{"error":"products: pkg/server/models.ExportRequest: 'products' ([\"Invalid_Product\"]) invalid value provided (\"Invalid_Product\")."}`,
		},
		{
			request: models.ExportRequest{
				Environment:  "sandbox",
				Implementer:  "implementer",
				AuthorisedBy: "authorised_by",
				JobTitle:     "job_title",
				Products: []string{
					"Business",
					"Invalid_Product",
				},
				HasAgreed:           true,
				AddDigitalSignature: false,
			},
			err: `{"error":"products: pkg/server/models.ExportRequest: 'products' ([\"Business\" \"Invalid_Product\"]) invalid value provided (\"Invalid_Product\")."}`,
		},
	}

	for _, testCase := range testCases {
		requestJSON, err := json.MarshalIndent(testCase.request, ``, `  `)
		require.NoError(err)
		require.NotNil(requestJSON)

		// we can remove this sleep if we change `bitbucket.org/openbankingteam/conformance-suite/pkg/time`.Layout to `time.RFC3339Nano`.
		time.Sleep(1 * time.Second)
		code, body, headers := request(http.MethodPost, "/api/export", bytes.NewReader(requestJSON), server)

		// Body contains error, so it cannot be nil.
		require.Equal(testCase.err, body.String(), body.String())
		require.Equal(http.StatusBadRequest, code, body.String())
		require.Equal(http.Header{
			echo.HeaderContentType: []string{
				echo.MIMEApplicationJSONCharsetUTF8,
			},
		}, headers, body.String())
	}
}
