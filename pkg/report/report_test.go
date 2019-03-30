package report

import (
	"reflect"
	"testing"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/server"
)

func TestReport_Validate(t *testing.T) {
	type fields struct {
		ID             string
		Created        string
		Expiration     string
		Version        string
		Status         ReportStatus
		CertifiedBy    ReportCertifiedBy
		SignatureChain *[]ReportSignatureChain
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr string
	}{
		// invalid cases
		{
			name:    "blank_all",
			fields:  fields{},
			wantErr: "certifiedBy: (authorisedBy: cannot be blank; brand: cannot be blank; environment: cannot be blank; jobTitle: cannot be blank.); created: cannot be blank; id: cannot be blank; status: cannot be blank; version: cannot be blank.",
		},
		{
			name: "blank_id",
			fields: fields{
				Created:    time.Now().Format(time.RFC3339),
				Expiration: time.Now().Format(time.RFC3339),
				Version:    "Version",
				Status:     ReportStatusPending,
				CertifiedBy: ReportCertifiedBy{
					Environment:  ReportCertifiedByEnvironmentTesting,
					Brand:        "Brand",
					AuthorisedBy: "AuthorisedBy",
					JobTitle:     "JobTitle",
				},
			},
			wantErr: "id: cannot be blank.",
		},
		{
			name: "invalid_id_format",
			fields: fields{
				ID:         "id_invalid",
				Created:    time.Now().Format(time.RFC3339),
				Expiration: time.Now().Format(time.RFC3339),
				Version:    "Version",
				Status:     ReportStatusPending,
				CertifiedBy: ReportCertifiedBy{
					Environment:  ReportCertifiedByEnvironmentTesting,
					Brand:        "Brand",
					AuthorisedBy: "AuthorisedBy",
					JobTitle:     "JobTitle",
				},
			},
			wantErr: "id: must be a valid UUID v4.",
		},
		{
			name: "blank_created",
			fields: fields{
				ID:         "f47ac10b-58cc-4372-8567-0e02b2c3d479",
				Expiration: time.Now().Format(time.RFC3339),
				Version:    "Version",
				Status:     ReportStatusPending,
				CertifiedBy: ReportCertifiedBy{
					Environment:  ReportCertifiedByEnvironmentTesting,
					Brand:        "Brand",
					AuthorisedBy: "AuthorisedBy",
					JobTitle:     "JobTitle",
				},
			},
			wantErr: "created: cannot be blank.",
		},
		{
			name: "invalid_created_time_format",
			fields: fields{
				ID:         "f47ac10b-58cc-4372-8567-0e02b2c3d479",
				Created:    time.Now().Format(time.ANSIC),
				Expiration: time.Now().Format(time.RFC3339),
				Version:    "Version",
				Status:     ReportStatusPending,
				CertifiedBy: ReportCertifiedBy{
					Environment:  ReportCertifiedByEnvironmentTesting,
					Brand:        "Brand",
					AuthorisedBy: "AuthorisedBy",
					JobTitle:     "JobTitle",
				},
			},
			wantErr: "created: must be a valid date.",
		},
		{
			name: "invalid_expiration_time_format",
			fields: fields{
				ID:         "f47ac10b-58cc-4372-8567-0e02b2c3d479",
				Created:    time.Now().Format(time.RFC3339),
				Expiration: time.Now().Format(time.ANSIC),
				Version:    "Version",
				Status:     ReportStatusPending,
				CertifiedBy: ReportCertifiedBy{
					Environment:  ReportCertifiedByEnvironmentTesting,
					Brand:        "Brand",
					AuthorisedBy: "AuthorisedBy",
					JobTitle:     "JobTitle",
				},
			},
			wantErr: "expiration: must be a valid date.",
		},
		{
			name: "invalid_status_value",
			fields: fields{
				ID:         "f47ac10b-58cc-4372-8567-0e02b2c3d479",
				Created:    time.Now().Format(time.RFC3339),
				Expiration: time.Now().Format(time.RFC3339),
				Version:    "Version",
				Status:     ReportStatus(-1),
				CertifiedBy: ReportCertifiedBy{
					Environment:  ReportCertifiedByEnvironmentTesting,
					Brand:        "Brand",
					AuthorisedBy: "AuthorisedBy",
					JobTitle:     "JobTitle",
				},
			},
			wantErr: "status: must be a valid value.",
		},
		// valid cases
		{
			name: "valid",
			fields: fields{
				ID:         "f47ac10b-58cc-4372-8567-0e02b2c3d479",
				Created:    time.Now().Format(time.RFC3339),
				Expiration: time.Now().Format(time.RFC3339),
				Version:    "Version",
				Status:     ReportStatusComplete,
				CertifiedBy: ReportCertifiedBy{
					Environment:  ReportCertifiedByEnvironmentTesting,
					Brand:        "Brand",
					AuthorisedBy: "AuthorisedBy",
					JobTitle:     "JobTitle",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require := test.NewRequire(t)

			r := Report{
				ID:             tt.fields.ID,
				Created:        tt.fields.Created,
				Expiration:     tt.fields.Expiration,
				Version:        tt.fields.Version,
				Status:         tt.fields.Status,
				CertifiedBy:    tt.fields.CertifiedBy,
				SignatureChain: tt.fields.SignatureChain,
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

func TestNewReport(t *testing.T) {
	t.Parallel()
	// TODO: add test cases once functionality is read. Intentionally skipping test for now.
	t.Skip()

	type args struct {
		exportResponse server.ExportResponse
	}
	tests := []struct {
		name    string
		args    args
		want    Report
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewReport(tt.args.exportResponse)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewReport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReport() = %v, want %v", got, tt.want)
			}
		})
	}
}
