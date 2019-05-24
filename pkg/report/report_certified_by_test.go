package report

import (
	"strings"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
)

func TestReportCertifiedBy_Validate(t *testing.T) {
	type fields struct {
		Environment  CertifiedByEnvironment
		Brand        string
		AuthorisedBy string
		JobTitle     string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr string
	}{
		// invalid cases
		{
			name:    "blank_all",
			wantErr: "authorisedBy: cannot be blank; brand: cannot be blank; environment: cannot be blank; jobTitle: cannot be blank.",
		},
		{
			name: "blank_brand",
			fields: fields{
				Environment:  CertifiedByEnvironmentTesting,
				JobTitle:     "JobTitle",
				AuthorisedBy: "AuthorisedBy",
			},
			wantErr: "brand: cannot be blank.",
		},
		{
			name: "blank_authorisedBy",
			fields: fields{
				Environment: CertifiedByEnvironmentTesting,
				Brand:       "Brand",
				JobTitle:    "JobTitle",
			},
			wantErr: "authorisedBy: cannot be blank.",
		},
		{
			name: "blank_jobTitle",
			fields: fields{
				Environment:  CertifiedByEnvironmentTesting,
				Brand:        "Brand",
				AuthorisedBy: "AuthorisedBy",
			},
			wantErr: "jobTitle: cannot be blank.",
		},
		// invalid Environment
		{
			name: "invalid_environment",
			fields: fields{
				Environment:  CertifiedByEnvironment(-1),
				Brand:        "Brand",
				AuthorisedBy: "AuthorisedBy",
				JobTitle:     "JobTitle",
			},
			wantErr: "environment: must be a valid value.",
		},
		// check lengths > 60
		{
			name: "length_brand",
			fields: fields{
				Environment:  CertifiedByEnvironmentTesting,
				Brand:        strings.Repeat("b", 61),
				AuthorisedBy: "AuthorisedBy",
				JobTitle:     "JobTitle",
			},
			wantErr: "brand: the length must be between 1 and 60.",
		},
		{
			name: "length_authorisedBy",
			fields: fields{
				Environment:  CertifiedByEnvironmentTesting,
				Brand:        "Brand",
				AuthorisedBy: strings.Repeat("a", 61),
				JobTitle:     "JobTitle",
			},
			wantErr: "authorisedBy: the length must be between 1 and 60.",
		},
		{
			name: "length_jobTitle",
			fields: fields{
				Environment:  CertifiedByEnvironmentTesting,
				Brand:        "Brand",
				AuthorisedBy: "AuthorisedBy",
				JobTitle:     strings.Repeat("j", 61),
			},
			wantErr: "jobTitle: the length must be between 1 and 60.",
		},
		// valid cases
		{
			name: "valid_testing",
			fields: fields{
				Environment:  CertifiedByEnvironmentTesting,
				Brand:        "Brand",
				AuthorisedBy: "AuthorisedBy",
				JobTitle:     "JobTitle",
			},
			wantErr: "",
		},
		{
			name: "valid_production",
			fields: fields{
				Environment:  CertifiedByEnvironmentTesting,
				Brand:        "Brand",
				AuthorisedBy: "AuthorisedBy",
				JobTitle:     "JobTitle",
			},
			wantErr: "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require := test.NewRequire(t)

			r := CertifiedBy{
				Environment:  tt.fields.Environment,
				Brand:        tt.fields.Brand,
				AuthorisedBy: tt.fields.AuthorisedBy,
				JobTitle:     tt.fields.JobTitle,
			}
			err := r.Validate()

			if tt.wantErr == "" {
				require.NoError(err)
			} else {
				require.EqualError(err, tt.wantErr)
			}
		})
	}
}
