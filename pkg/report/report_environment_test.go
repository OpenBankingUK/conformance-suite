package report

import (
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
)

func TestReportCertifiedByEnvironment_String(t *testing.T) {
	tests := []struct {
		name string
		r    ReportCertifiedByEnvironment
		want string
	}{
		{
			name: "ReportCertifiedByEnvironmentTesting",
			r:    ReportCertifiedByEnvironmentTesting,
			want: "testing",
		},
		{
			name: "ReportCertifiedByEnvironmentProduction",
			r:    ReportCertifiedByEnvironmentProduction,
			want: "production",
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

func TestReportCertifiedByEnvironment_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		r       ReportCertifiedByEnvironment
		want    []byte
		wantErr string
	}{
		{
			name: "ReportCertifiedByEnvironmentTesting",
			r:    ReportCertifiedByEnvironmentTesting,
			want: []byte(`"testing"`),
		},
		{
			name: "ReportCertifiedByEnvironmentProduction",
			r:    ReportCertifiedByEnvironmentProduction,
			want: []byte(`"production"`),
		},
		{
			name:    "ReportCertifiedByEnvironmentFake",
			r:       ReportCertifiedByEnvironment(-1),
			want:    []byte(nil),
			wantErr: "-1 is an invalid enum for ReportCertifiedByEnvironment",
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

func TestReportCertifiedByEnvironment_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		r       *ReportCertifiedByEnvironment
		args    args
		wantErr string
	}{
		{
			name: "ReportCertifiedByEnvironmentTesting",
			r:    envToEnvPointer(ReportCertifiedByEnvironmentTesting),
			args: args{
				data: []byte(`"testing"`),
			},
		},
		{
			name: "ReportCertifiedByEnvironmentProduction",
			r:    envToEnvPointer(ReportCertifiedByEnvironmentProduction),
			args: args{
				data: []byte(`"production"`),
			},
		},
		{
			name: "ReportCertifiedByEnvironmentFake",
			r:    envToEnvPointer(ReportCertifiedByEnvironment(-1)),
			args: args{
				data: []byte(`"fake"`),
			},
			wantErr: `"fake" is an invalid enum for ReportCertifiedByEnvironment`,
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

func envToEnvPointer(r ReportCertifiedByEnvironment) *ReportCertifiedByEnvironment {
	return &r
}
