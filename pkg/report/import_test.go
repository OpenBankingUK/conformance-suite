package report

import (
	"io"
	"os"
	"strings"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
	"github.com/stretchr/testify/require"
)

func TestNewZipImporter(t *testing.T) {
	require := test.NewRequire(t)

	reader := strings.NewReader("")
	require.NotNil(NewZipImporter(reader))
}

func Test_zipImporter_Import(t *testing.T) {
	reportValid, err := os.Open("./testdata/report_valid.zip")
	require.NotNil(t, reportValid)
	require.NoError(t, err)
	reportInvalid, err := os.Open("./testdata/report_invalid.zip")
	require.NotNil(t, reportInvalid)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, reportValid.Close())
		require.NoError(t, reportInvalid.Close())
	}()

	type fields struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		want    Report
		wantErr string
	}{
		// valid cases
		{
			name: "valid_case_1",
			fields: fields{
				reader: reportValid,
			},
			want:    Report{},
			wantErr: "",
		},
		// invalid cases
		{
			name: "invalid_case_1",
			fields: fields{
				reader: reportInvalid,
			},
			want:    Report{},
			wantErr: `zipImporter.Import: could not find "report.json" in ZIP archive`,
		},
	}

	// See: "Cleaning up after a group of parallel tests" in https: //blog.golang.org/subtests
	t.Run("group", func(t *testing.T) {
		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				require := test.NewRequire(t)

				importer := NewZipImporter(tt.fields.reader)
				actual, err := importer.Import()

				// expecting error, so assert and return
				if tt.wantErr != "" {
					require.EqualError(err, tt.wantErr)
					return
				}

				require.NoError(err)
				require.Equal(tt.want, actual)
			})
		}
	})
}
