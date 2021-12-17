package report

import (
	"reflect"
	"testing"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/events"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/results"
	"github.com/stretchr/testify/assert"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/server/models"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
)

func TestReport_Validate(t *testing.T) {
	type fields struct {
		ID             string
		Created        string
		Expiration     string
		Version        string
		Status         Status
		CertifiedBy    CertifiedBy
		SignatureChain *[]SignatureChain
		AgreedTC       bool
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
				Status:     StatusPending,
				CertifiedBy: CertifiedBy{
					Environment:  CertifiedByEnvironmentTesting,
					Brand:        "Brand",
					AuthorisedBy: "AuthorisedBy",
					JobTitle:     "JobTitle",
				},
			},
			wantErr: "id: cannot be blank.",
		},
		{
			name: "has_not_agreed_production", // non production envs are tested implicitly with all other cases
			fields: fields{
				ID:         "f47ac10b-58cc-4372-8567-0e02b2c3d479",
				Created:    time.Now().Format(time.RFC3339),
				Expiration: time.Now().Format(time.RFC3339),
				Version:    "Version",
				Status:     StatusPending,
				CertifiedBy: CertifiedBy{
					Environment:  CertifiedByEnvironmentProduction,
					Brand:        "Brand",
					AuthorisedBy: "AuthorisedBy",
					JobTitle:     "JobTitle",
				},
				AgreedTC: false,
			},
			wantErr: "agreedTermsConditions: cannot be blank.",
		},
		{
			name: "invalid_id_format",
			fields: fields{
				ID:         "id_invalid",
				Created:    time.Now().Format(time.RFC3339),
				Expiration: time.Now().Format(time.RFC3339),
				Version:    "Version",
				Status:     StatusPending,
				CertifiedBy: CertifiedBy{
					Environment:  CertifiedByEnvironmentTesting,
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
				Status:     StatusPending,
				CertifiedBy: CertifiedBy{
					Environment:  CertifiedByEnvironmentTesting,
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
				Status:     StatusPending,
				CertifiedBy: CertifiedBy{
					Environment:  CertifiedByEnvironmentTesting,
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
				Status:     StatusPending,
				CertifiedBy: CertifiedBy{
					Environment:  CertifiedByEnvironmentTesting,
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
				Status:     Status(-1),
				CertifiedBy: CertifiedBy{
					Environment:  CertifiedByEnvironmentTesting,
					Brand:        "Brand",
					AuthorisedBy: "AuthorisedBy",
					JobTitle:     "JobTitle",
				},
			},
			wantErr: "status: must be a valid value.",
		},
		// valid cases
		{
			name: "valid_testing_env",
			fields: fields{
				ID:         "f47ac10b-58cc-4372-8567-0e02b2c3d479",
				Created:    time.Now().Format(time.RFC3339),
				Expiration: time.Now().Format(time.RFC3339),
				Version:    "Version",
				Status:     StatusComplete,
				CertifiedBy: CertifiedBy{
					Environment:  CertifiedByEnvironmentTesting,
					Brand:        "Brand",
					AuthorisedBy: "AuthorisedBy",
					JobTitle:     "JobTitle",
				},
			},
		},
		{
			name: "valid_sandbox_env",
			fields: fields{
				ID:         "f47ac10b-58cc-4372-8567-0e02b2c3d479",
				Created:    time.Now().Format(time.RFC3339),
				Expiration: time.Now().Format(time.RFC3339),
				Version:    "Version",
				Status:     StatusComplete,
				CertifiedBy: CertifiedBy{
					Environment:  CertifiedByEnvironmentSandbox,
					Brand:        "Brand",
					AuthorisedBy: "AuthorisedBy",
					JobTitle:     "JobTitle",
				},
			},
		},
		{
			name: "valid_production_env",
			fields: fields{
				ID:         "f47ac10b-58cc-4372-8567-0e02b2c3d479",
				Created:    time.Now().Format(time.RFC3339),
				Expiration: time.Now().Format(time.RFC3339),
				Version:    "Version",
				Status:     StatusComplete,
				CertifiedBy: CertifiedBy{
					Environment:  CertifiedByEnvironmentProduction,
					Brand:        "Brand",
					AuthorisedBy: "AuthorisedBy",
					JobTitle:     "JobTitle",
				},
				AgreedTC: true,
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
				Expiration:     stringToPointer(tt.fields.Expiration),
				Version:        tt.fields.Version,
				Status:         tt.fields.Status,
				CertifiedBy:    tt.fields.CertifiedBy,
				SignatureChain: tt.fields.SignatureChain,
				AgreedTC:       tt.fields.AgreedTC,
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

func stubExportResults() models.ExportResults {
	return models.ExportResults{
		Tokens: []events.AcquiredAccessToken{
			{
				TokenName: "token-name",
			},
		},
		DiscoveryModel: discovery.Model{
			DiscoveryModel: discovery.ModelDiscovery{
				Name:             "Name",
				Description:      "Description",
				DiscoveryVersion: "Discovery version",
				TokenAcquisition: "Token Acquisition",
				DiscoveryItems: []discovery.ModelDiscoveryItem{
					{
						APISpecification: discovery.ModelAPISpecification{
							Name:          "Name",
							SchemaVersion: "schema-version",
							Version:       "version",
							Manifest:      "manifest",
							URL:           "url",
							SpecType:      "specType",
						},
						ResourceBaseURI:        "resource-base-uri",
						OpenidConfigurationURI: "open-id-configuration-id",
						ResourceIds: discovery.ResourceIds{
							"foo": "bar",
						},
						Endpoints: []discovery.ModelEndpoint{
							{
								Method:                "GET",
								Path:                  "/",
								ConditionalProperties: nil,
							},
						},
					},
				},
				CustomTests: []discovery.CustomTest{},
			},
		},
		ExportRequest: models.ExportRequest{
			Implementer:         "Implemented",
			AuthorisedBy:        "Authorised by",
			JobTitle:            "Job title",
			HasAgreed:           false,
			AddDigitalSignature: false,
		},
		HasPassed: false,
		Results:   map[results.ResultKey][]results.TestCase{},
	}
}

func TestReport_GetFails_Fails_Zero(t *testing.T) {
	require := test.NewRequire(t)

	specs := stubResults(true, true, true)
	expected := 0
	actual := GetFails(specs)

	require.Equal(expected, actual)
}

func TestReport_GetFails_Fails_Three(t *testing.T) {
	require := test.NewRequire(t)

	specs := stubResults(false, false, false)
	expected := 3
	actual := GetFails(specs)

	require.Equal(expected, actual)
}

func TestReport_GetFails_Fails_Four(t *testing.T) {
	require := test.NewRequire(t)

	specs := stubResults(false, false, false)
	spec1 := results.ResultKey{
		APIVersion: "APIVersion1",
		APIName:    "APIName1",
	}
	specs[spec1][1].Pass = false
	expected := 4
	actual := GetFails(specs)

	require.Equal(expected, actual)
}

func TestNewReport(t *testing.T) {
	t.Parallel()
	// TODO: add test cases once functionality is read. Intentionally skipping test for now.
	t.Skip()

	type args struct {
		exportResults models.ExportResults
	}
	tests := []struct {
		name string
		args args
		want Report
		err  error // wantErr can be inferred by this being nil or not
	}{
		{
			name: "Valid report",
			args: args{
				exportResults: stubExportResults(),
			},
			want: Report{},
			err:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewReport(tt.args.exportResults, "Testing")
			if tt.err != nil {
				assert.New(t).Equal(tt.err, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func stubResults(pass1, pass2, pass3 bool) map[results.ResultKey][]results.TestCase {
	specs := map[results.ResultKey][]results.TestCase{}

	spec1 := results.ResultKey{
		APIVersion: "APIVersion1",
		APIName:    "APIName1",
	}
	specs[spec1] = []results.TestCase{
		results.NewTestCaseResult("1.1", pass1, results.NoMetrics(), nil, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri", "200 0K"),
		results.NewTestCaseResult("1.2", true, results.NoMetrics(), nil, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri", "200 0K"),
		results.NewTestCaseResult("1.3", true, results.NoMetrics(), nil, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri", "200 0K"),
	}

	spec2 := results.ResultKey{
		APIVersion: "APIVersion2",
		APIName:    "APIName2",
	}
	specs[spec2] = []results.TestCase{
		results.NewTestCaseResult("2.1", true, results.NoMetrics(), nil, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri", "200 0K"),
		results.NewTestCaseResult("2.2", pass2, results.NoMetrics(), nil, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri", "200 0K"),
		results.NewTestCaseResult("2.3", true, results.NoMetrics(), nil, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri", "200 0K"),
	}

	spec3 := results.ResultKey{
		APIVersion: "APIVersion3",
		APIName:    "APIName3",
	}
	specs[spec3] = []results.TestCase{
		results.NewTestCaseResult("3.1", true, results.NoMetrics(), nil, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri", "200 0K"),
		results.NewTestCaseResult("3.2", true, results.NoMetrics(), nil, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri", "200 0K"),
		results.NewTestCaseResult("3.3", pass3, results.NoMetrics(), nil, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri", "200 0K"),
	}

	spec4 := results.ResultKey{
		APIVersion: "APIVersion4",
		APIName:    "APIName4",
	}
	specs[spec4] = []results.TestCase{
		results.NewTestCaseResult("4.1", true, results.NoMetrics(), nil, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri", "200 0K"),
		results.NewTestCaseResult("4.2", true, results.NoMetrics(), nil, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri", "200 0K"),
		results.NewTestCaseResult("4.3", true, results.NoMetrics(), nil, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri", "200 0K"),
	}

	return specs
}

func stringToPointer(str string) *string {
	return &str
}
