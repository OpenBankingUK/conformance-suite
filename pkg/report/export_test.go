package report

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
	"github.com/stretchr/testify/require"
)

func TestNewZipExporter(t *testing.T) {
	require := test.NewRequire(t)

	report := Report{}
	writer := bytes.NewBuffer([]byte{})
	require.NotNil(NewZipExporter(report, writer))
}

func Test_zipExporter_Export(t *testing.T) {
	t.Skip()
	tempDir, err := ioutil.TempDir("", "Test_zipExporter_Export")
	require.NoError(t, err)

	t.Log("tempDir:", tempDir)

	type fields struct {
		report   Report
		filename string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr string
	}{
		// invalid cases
		{
			name: "report_invalid_1.zip",
			fields: fields{
				report:   Report{},
				filename: filepath.Join(tempDir, "report_invalid_1.zip"),
			},
			wantErr: `zipExporter.Export: json.MarshalIndent failed, report={ID: Created: Expiration:<nil> Version: Status: CertifiedBy:{Environment: Brand: AuthorisedBy: JobTitle:} SignatureChain:<nil>}: json: error calling MarshalJSON for type report.Status: 0 is an invalid enum for Status`,
		},
		// valid cases
		{
			name: "report_valid_1.zip",
			fields: fields{
				report: Report{
					ID:         "f47ac10b-58cc-4372-8567-0e02b2c3d479",
					Created:    time.Now().Format(time.RFC3339),
					Expiration: stringToPointer(time.Now().Format(time.RFC3339)),
					Version:    "Version",
					Status:     StatusComplete,
					CertifiedBy: CertifiedBy{
						Environment:  CertifiedByEnvironmentTesting,
						Brand:        "Brand",
						AuthorisedBy: "AuthorisedBy",
						JobTitle:     "JobTitle",
					},
				},
				filename: filepath.Join(tempDir, "report_valid_1.zip"),
			},
			wantErr: "",
		},
	}

	// See: "Cleaning up after a group of parallel tests" in https: //blog.golang.org/subtests
	t.Run("group", func(t *testing.T) {
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				require := test.NewRequire(t)

				// Create zip file named `tt.fields.filename`
				zipFilename := tt.fields.filename
				writer, err := os.Create(zipFilename)
				require.NoError(err)
				defer func() {
					require.NoError(writer.Close())
				}()
				exporter := NewZipExporter(tt.fields.report, writer)

				// If we are expecting an error, check that is is present then return early.
				// Else, allow the tests to continue.
				if tt.wantErr != "" {
					require.EqualError(exporter.Export(), tt.wantErr)
					return
				}

				// Do the export
				require.NoError(exporter.Export())

				// Calculate what report.json we expect to see.
				expectedReport, err := json.MarshalIndent(tt.fields.report, marshalIndentPrefix, marshalIndent)
				require.NoError(err)

				// Read zip file and check contents match.
				// Open a zip archive for reading.
				zipReader, err := zip.OpenReader(zipFilename)
				require.NoError(err)
				defer func() {
					require.NoError(zipReader.Close())
				}()

				// Check that a file named `reportFilename` exists in the ZIP archive and store its index.
				hasReportFileName := false
				reportFileIndex := 0
				for zipFileIndex, zipFile := range zipReader.File {
					if zipFile.Name == reportFilename {
						hasReportFileName = true
						reportFileIndex = zipFileIndex
						break
					} else {
						// ignore anything that isn't the report.json
						t.Logf("ignoring non %q report file (file name is not equal to %q)\n", zipFile.Name, reportFilename)
					}
				}
				require.True(hasReportFileName)

				// Store the report file
				reportFile := zipReader.File[reportFileIndex]
				t.Logf("reportFile=%+v\n", reportFile)

				// Check uncompress size against expectedReport
				require.EqualValues(len(expectedReport), reportFile.UncompressedSize)

				// Open so we can read it's content
				reportFileReader, err := reportFile.Open()
				defer func() {
					require.NoError(reportFileReader.Close())
				}()
				require.NoError(err)

				buf := bytes.NewBuffer([]byte{})
				n, err := buf.ReadFrom(reportFileReader)
				require.NoError(err)
				require.NotZero(n)
				actualReport := buf.String()

				// Compare contents to what we expect
				t.Logf("expectedReport=%+v\n", string(expectedReport))
				t.Logf("actualReport=%+v\n", actualReport)
				require.JSONEq(string(expectedReport), actualReport)
			})
		}
	})

	t.Log("cleanup, tempDir:", tempDir)
	require.NoError(t, os.RemoveAll(tempDir))
}
