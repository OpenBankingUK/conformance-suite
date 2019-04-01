package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	internal_time "bitbucket.org/openbankingteam/conformance-suite/internal/pkg/time"
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version/mocks"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/report"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/server/models"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

func TestServerPostExport(t *testing.T) {
	require := test.NewRequire(t)

	server := NewServer(testJourney(), nullLogger(), &mocks.Version{})
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	expected := models.ExportRequest{
		Implementer:  "implementer",
		AuthorisedBy: "authorised_by",
		JobTitle:     "job_title",
		HasAgreed:    true,
	}
	requestJSON, err := json.MarshalIndent(expected, ``, `  `)
	require.NoError(err)
	require.NotNil(requestJSON)

	// make the request and record time the export operation was started
	exportTime := time.Now()
	// we can remove this sleep if we change `bitbucket.org/openbankingteam/conformance-suite/internal/pkg/time`.Layout to `time.RFC3339Nano`.
	time.Sleep(1 * time.Second)
	code, body, headers := request(http.MethodPost, "/api/export", bytes.NewReader(requestJSON), server)

	require.Equal(http.StatusOK, code)
	require.Equal(http.Header{
		echo.HeaderContentType: []string{
			MIMEApplicationZIP,
		},
		// echo.HeaderContentDisposition: []string{
		// 	`attachment; filename="report.zip"`,
		// },
	}, headers)

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
	require.Equal(report.CertifiedByEnvironmentTesting, actual.CertifiedBy.Environment)
	require.Equal(expected.Implementer, actual.CertifiedBy.Brand)
	require.Equal(expected.AuthorisedBy, actual.CertifiedBy.AuthorisedBy)
	require.Equal(expected.JobTitle, actual.CertifiedBy.JobTitle)
}
