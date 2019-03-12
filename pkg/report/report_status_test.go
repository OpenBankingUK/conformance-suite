package report

import (
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
)

func TestReportStatus_String(t *testing.T) {
	tests := []struct {
		name string
		r    ReportStatus
		want string
	}{
		{
			name: "ReportStatusPending",
			r:    ReportStatusPending,
			want: "Pending",
		},
		{
			name: "ReportStatusComplete",
			r:    ReportStatusComplete,
			want: "Complete",
		},
		{
			name: "ReportStatusError",
			r:    ReportStatusError,
			want: "Error",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require := test.NewRequire(t)
			require.Equal(tt.want, tt.r.String())
		})
	}
}

func TestReportStatus_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		r       ReportStatus
		want    []byte
		wantErr string
	}{
		{
			name: "ReportStatusPending",
			r:    ReportStatusPending,
			want: []byte(`"Pending"`),
		},
		{
			name: "ReportStatusComplete",
			r:    ReportStatusComplete,
			want: []byte(`"Complete"`),
		},
		{
			name: "ReportStatusError",
			r:    ReportStatusError,
			want: []byte(`"Error"`),
		},
		{
			name:    "ReportStatusFake",
			r:       ReportStatus(-1),
			wantErr: "-1 is an invalid enum for ReportStatus",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require := test.NewRequire(t)
			got, err := tt.r.MarshalJSON()

			if tt.wantErr == "" {
				require.NoError(err)
			} else {
				require.EqualError(err, tt.wantErr)
			}
			require.Equal(tt.want, got)
		})
	}
}

func TestReportStatus_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		r       *ReportStatus
		args    args
		wantErr string
	}{
		{
			name: "ReportStatusPending",
			r:    statusToStatusPointer(ReportStatusPending),
			args: args{
				data: []byte(`"Pending"`),
			},
		},
		{
			name: "ReportStatusComplete",
			r:    statusToStatusPointer(ReportStatusComplete),
			args: args{
				data: []byte(`"Complete"`),
			},
		},
		{
			name: "ReportStatusError",
			r:    statusToStatusPointer(ReportStatusError),
			args: args{
				data: []byte(`"Error"`),
			},
		},
		{
			name: "ReportStatusFake",
			r:    statusToStatusPointer(ReportStatus(-1)),
			args: args{
				data: []byte(`"fake"`),
			},
			wantErr: `"fake" is an invalid enum for ReportStatus`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require := test.NewRequire(t)
			err := tt.r.UnmarshalJSON(tt.args.data)

			if tt.wantErr == "" {
				require.NoError(err)
			} else {
				require.EqualError(err, tt.wantErr)
			}
		})
	}
}

func statusToStatusPointer(r ReportStatus) *ReportStatus {
	return &r
}
