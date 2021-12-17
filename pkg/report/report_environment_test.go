package report

import (
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
)

func TestReportCertifiedByEnvironment_String(t *testing.T) {
	tests := []struct {
		name string
		r    CertifiedByEnvironment
		want string
	}{
		{
			name: "CertifiedByEnvironmentTesting",
			r:    CertifiedByEnvironmentTesting,
			want: "testing",
		},
		{
			name: "CertifiedByEnvironmentProduction",
			r:    CertifiedByEnvironmentProduction,
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
		r       CertifiedByEnvironment
		want    []byte
		wantErr string
	}{
		{
			name: "CertifiedByEnvironmentTesting",
			r:    CertifiedByEnvironmentTesting,
			want: []byte(`"testing"`),
		},
		{
			name: "CertifiedByEnvironmentProduction",
			r:    CertifiedByEnvironmentProduction,
			want: []byte(`"production"`),
		},
		{
			name:    "CertifiedByEnvironmentFake",
			r:       CertifiedByEnvironment(-1),
			want:    []byte(nil),
			wantErr: "-1 is an invalid enum for CertifiedByEnvironment",
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
		r       *CertifiedByEnvironment
		args    args
		wantErr string
	}{
		{
			name: "CertifiedByEnvironmentTesting",
			r:    envToEnvPointer(CertifiedByEnvironmentTesting),
			args: args{
				data: []byte(`"testing"`),
			},
		},
		{
			name: "CertifiedByEnvironmentProduction",
			r:    envToEnvPointer(CertifiedByEnvironmentProduction),
			args: args{
				data: []byte(`"production"`),
			},
		},
		{
			name: "CertifiedByEnvironmentFake",
			r:    envToEnvPointer(CertifiedByEnvironment(-1)),
			args: args{
				data: []byte(`"fake"`),
			},
			wantErr: `"fake" is an invalid enum for CertifiedByEnvironment`,
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

func envToEnvPointer(r CertifiedByEnvironment) *CertifiedByEnvironment {
	return &r
}
