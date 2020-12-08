package discovery

import (
	"errors"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"

	"github.com/stretchr/testify/assert"
)

// conditionalityCheckerMock - implements model.ConditionalityChecker interface for tests
type conditionalityCheckerMock struct {
	isPresent           bool
	isPresentErr        error
	missingMandatory    []model.Input
	missingMandatoryErr error
}

// IsOptional - not used in discovery test
func (c conditionalityCheckerMock) IsOptional(method, endpoint string, specification string) (bool, error) {
	return false, nil
}

// Returns IsMandatory true for POST /account-access-consents, false for all other endpoint/methods.
func (c conditionalityCheckerMock) IsMandatory(method, endpoint string, specification string) (bool, error) {
	if method == "POST" && endpoint == "/account-access-consents" {
		return true, nil
	}
	return false, nil
}

// IsConditional - not used in discovery test
func (c conditionalityCheckerMock) IsConditional(method, endpoint string, specification string) (bool, error) {
	return false, nil
}

// IsPresent - returns stubbed isPresent boolean value
func (c conditionalityCheckerMock) IsPresent(method, endpoint string, specification string) (bool, error) {
	return c.isPresent, c.isPresentErr
}

// MissingMandatory - returns stubbed array of missing endpoints
func (c conditionalityCheckerMock) MissingMandatory(endpoints []model.Input, specification string) ([]model.Input, error) {
	return c.missingMandatory, c.missingMandatoryErr
}

// UnmarshalDiscoveryJSON - returns discovery model
func testUnmarshalDiscoveryJSON(t *testing.T, discoveryJSON string) *Model {
	t.Helper()

	discovery, err := UnmarshalDiscoveryJSON(discoveryJSON)
	assert.NoError(t, err)
	return discovery
}

// discoveryStub - returns discovery JSON with given field stubbed with given value
func discoveryStub(field string, value string) string {
	name := "ob-v3.1-generic"
	description := "An Open Banking UK generic discovery template for v3.1 of Accounts and Payments."
	version := "v0.4.0"
	tokenAcquisition := "psu"
	specName := "Account and Transaction API Specification"
	specURL := "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937820271/Account+and+Transaction+API+Specification+-+v3.1"
	specVersion := "v3.1.0"
	schemaVersion := "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json"
	manifest := "https://www.example.com"
	endpoints := `, "endpoints": [
		{
			"method": "POST",
			"path": "/account-access-consents"
		},
		{
			"method": "GET",
			"path": "/accounts/{AccountId}/balances"
		}
	]`

	switch field {
	case "version":
		version = value
	case "endpoints":
		if value == "" {
			endpoints = ""
		} else {
			endpoints = `, "endpoints": ` + value
		}
	case "specName":
		specName = value
	case "schemaVersion":
		schemaVersion = value
	case "manifest":
		manifest = value
	case "tokenAcquisition":
		tokenAcquisition = value
	case "specURL":
		specURL = value
	case "specVersion":
		specVersion = value
	case "name":
		name = value
	case "description":
		description = value
	}

	apiSpecification := apiSpecificationStub(specName, specURL, specVersion, schemaVersion, manifest, field, value)

	discoveryItems := discoveryItemsStub(apiSpecification, endpoints, field, value)

	return `
	{
		"discoveryModel": {` +
		`"name": "` + name + `",` +
		`"description": "` + description + `",` +
		`"discoveryVersion": "` + version + `",` +
		`"tokenAcquisition": "` + tokenAcquisition + `"` +
		discoveryItems + `
		}
	}`
}

func apiSpecificationStub(specName string, specURL string, specVersion string, schemaVersion string, manifest, field string, value string) string {
	apiSpecification := `"apiSpecification": {
			"name": "` + specName + `",
			"url": "` + specURL + `",
			"version": "` + specVersion + `",
			"schemaVersion": "` + schemaVersion + `",
			"manifest":	"` + manifest + `"
		},`
	if field == "apiSpecification" {
		if value == "" {
			apiSpecification = ""
		} else {
			apiSpecification = `"apiSpecification": ` + value + `,`
		}
	}
	return apiSpecification
}

func discoveryItemsStub(apiSpecification string, endpoints string, field string, value string) string {
	discoveryItems := `, "discoveryItems": [
			{
				` + apiSpecification + `
				"openidConfigurationUri": "https://ob19-auth1-ui.o3bank.co.uk/.well-known/openid-configuration",
				"resourceBaseUri": "https://ob19-rs1.o3bank.co.uk:4501/open-banking/v3.1/aisp"` + endpoints + `
			}
		]`
	if field == "discoveryItems" {
		if value == "" {
			discoveryItems = ""
		} else {
			discoveryItems = `, "discoveryItems": ` + value
		}
	}
	return discoveryItems
}

// TestValidate - test Validate function
func TestValidate(t *testing.T) {

	// invalidTestCase
	// discoveryJSON - the discovery model JSON
	// failures - list of ValidationFailure structs
	// err - the expected err
	type invalidTest struct {
		discoveryJSON string
		success       bool
		failures      []ValidationFailure
		err           error
	}

	// testValidateFailures - run Validate, and test validation failure expectations
	testValidateFailures := func(t *testing.T, checker model.ConditionalityChecker, expected *invalidTest) {
		t.Helper()

		discovery := testUnmarshalDiscoveryJSON(t, expected.discoveryJSON)
		result, failures, err := Validate(checker, discovery)
		assert.Equal(t, expected.success, result)
		assert.Equal(t, expected.err, err)
		assert.Equal(t, expected.failures, failures)
	}

	t.Run("name missing returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{}, &invalidTest{
			discoveryJSON: discoveryStub("name", ""),
			failures: []ValidationFailure{
				{
					Key:   "DiscoveryModel.Name",
					Error: "Field 'DiscoveryModel.Name' is required",
				},
			},
		})
	})

	t.Run("description missing returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{}, &invalidTest{
			discoveryJSON: discoveryStub("description", ""),
			failures: []ValidationFailure{
				{
					Key:   "DiscoveryModel.Description",
					Error: "Field 'DiscoveryModel.Description' is required",
				},
			},
		})
	})

	t.Run("when discoveryVersion missing returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{}, &invalidTest{
			discoveryJSON: discoveryStub("version", ""),
			failures: []ValidationFailure{
				{
					Key:   "DiscoveryModel.DiscoveryVersion",
					Error: "Field 'DiscoveryModel.DiscoveryVersion' is required",
				},
			}})
	})

	t.Run("when version not in discovery.SupportedVersions() returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{isPresent: true}, &invalidTest{
			discoveryJSON: discoveryStub("version", "v9.9.9"),
			failures: []ValidationFailure{
				{
					Key:   "DiscoveryModel.DiscoveryVersion",
					Error: "DiscoveryVersion 'v9.9.9' not in list of supported versions",
				},
			}})
	})

	t.Run("ensure that `psu`, `headless`, `store` are supported tokenAcquisition values", func(t *testing.T) {
		methods := []string{"psu", "headless", "store"}
		for _, method := range methods {
			testValidateFailures(t, conditionalityCheckerMock{isPresent: true}, &invalidTest{
				discoveryJSON: discoveryStub("tokenAcquisition", method),
				success:       true,
				failures:      []ValidationFailure{}})
		}
	})

	t.Run("when tokenAcquisition provided is not in SupportedTokenAcquisitions() returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{isPresent: true}, &invalidTest{
			discoveryJSON: discoveryStub("tokenAcquisition", "foo"),
			failures: []ValidationFailure{
				{
					Key:   "DiscoveryModel.TokenAcquisition",
					Error: "TokenAcquisition 'foo' not in list of supported methods",
				},
			}})
	})

	t.Run("when discoveryItems missing returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{}, &invalidTest{
			discoveryJSON: discoveryStub("discoveryItems", ""),
			failures: []ValidationFailure{
				{
					Key:   "DiscoveryModel.DiscoveryItems",
					Error: "Field 'DiscoveryModel.DiscoveryItems' is required",
				},
			}})
	})

	t.Run("when discoveryItem missing apiSpecification returns failures", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{}, &invalidTest{
			discoveryJSON: discoveryStub("apiSpecification", ""),
			failures: []ValidationFailure{
				{
					Key:   "DiscoveryModel.DiscoveryItems[0].APISpecification.Name",
					Error: "Field 'DiscoveryModel.DiscoveryItems[0].APISpecification.Name' is required",
				},
				{
					Key:   "DiscoveryModel.DiscoveryItems[0].APISpecification.URL",
					Error: "Field 'DiscoveryModel.DiscoveryItems[0].APISpecification.URL' is required",
				},
				{
					Key:   "DiscoveryModel.DiscoveryItems[0].APISpecification.Version",
					Error: "Field 'DiscoveryModel.DiscoveryItems[0].APISpecification.Version' is required",
				},
				{
					Key:   "DiscoveryModel.DiscoveryItems[0].APISpecification.SchemaVersion",
					Error: "Field 'DiscoveryModel.DiscoveryItems[0].APISpecification.SchemaVersion' is required",
				},
				{
					Key:   "DiscoveryModel.DiscoveryItems[0].APISpecification.Manifest",
					Error: "Field 'DiscoveryModel.DiscoveryItems[0].APISpecification.Manifest' is required",
				},
			}})
	})

	t.Run("when discoveryItem apiSpecification schemaVersion not in suite config returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{}, &invalidTest{
			discoveryJSON: discoveryStub("schemaVersion", "http://example.com/bad-schema"),
			failures: []ValidationFailure{
				{
					Key:   "DiscoveryModel.DiscoveryItems[0].APISpecification.SchemaVersion",
					Error: "'SchemaVersion' not supported by suite 'http://example.com/bad-schema'",
				},
			}})
	})

	t.Run("when discoveryItem apiSpecification schemaVersion in suite config but Name field not matching returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{isPresent: true}, &invalidTest{
			discoveryJSON: discoveryStub("specName", "Bad Spec Name"),
			failures: []ValidationFailure{
				{
					Key:   "DiscoveryModel.DiscoveryItems[0].APISpecification.Name",
					Error: "'Name' should be 'Account and Transaction API Specification' when schemaVersion is 'https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json'",
				},
			}})
	})

	t.Run("when discoveryItem apiSpecification schemaVersion in suite config but Version field not matching returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{isPresent: true}, &invalidTest{
			discoveryJSON: discoveryStub("specVersion", "v9.9.9"),
			failures: []ValidationFailure{
				{
					Key:   "DiscoveryModel.DiscoveryItems[0].APISpecification.Version",
					Error: "'Version' should be 'v3.1.0' when schemaVersion is 'https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json'",
				},
			}})
	})

	t.Run("when discoveryItem apiSpecification schemaVersion in suite config but URL field not matching returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{isPresent: true}, &invalidTest{
			discoveryJSON: discoveryStub("specURL", "http://example.com/bad-url"),
			failures: []ValidationFailure{
				{
					Key:   "DiscoveryModel.DiscoveryItems[0].APISpecification.URL",
					Error: "'URL' should be 'https://openbanking.atlassian.net/wiki/spaces/DZ/pages/937820271/Account+and+Transaction+API+Specification+-+v3.1' when schemaVersion is 'https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json'",
				},
			}})
	})

	t.Run("when discoveryItems has empty endpoints array returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{}, &invalidTest{
			discoveryJSON: discoveryStub("endpoints", "[]"),
			failures: []ValidationFailure{
				{
					Key:   "DiscoveryModel.DiscoveryItems[0].Endpoints",
					Error: "Field 'DiscoveryModel.DiscoveryItems[0].Endpoints' cannot be empty",
				},
			}})
	})

	t.Run("when conditionality checker isPresent throws error returns failures", func(t *testing.T) {
		stubIsPresentError := conditionalityCheckerMock{
			isPresent:    false,
			isPresentErr: errors.New("some error message"),
		}
		testValidateFailures(t, stubIsPresentError, &invalidTest{
			discoveryJSON: discoveryStub("", ""),
			failures: []ValidationFailure{
				{
					Key:   "DiscoveryModel.DiscoveryItems[0].Endpoints[0]",
					Error: "some error message",
				},
				{
					Key:   "DiscoveryModel.DiscoveryItems[0].Endpoints[1]",
					Error: "some error message",
				},
			},
		},
		)
	})

	t.Run("when conditionality checker reports endpoints not present returns failures", func(t *testing.T) {
		stubAllNotPresent := conditionalityCheckerMock{
			isPresent: false,
		}
		testValidateFailures(t, stubAllNotPresent, &invalidTest{
			discoveryJSON: discoveryStub("", ""),
			failures: []ValidationFailure{
				{
					Key:   "DiscoveryModel.DiscoveryItems[0].Endpoints[0]",
					Error: "Invalid endpoint Method='POST', Path='/account-access-consents'",
				},
				{
					Key:   "DiscoveryModel.DiscoveryItems[0].Endpoints[1]",
					Error: "Invalid endpoint Method='GET', Path='/accounts/{AccountId}/balances'",
				},
			},
		})
	})

	t.Run("when conditionality checker reports endpoints present returns no failures", func(t *testing.T) {
		stubAllPresent := conditionalityCheckerMock{
			isPresent: true,
		}
		testValidateFailures(t, stubAllPresent, &invalidTest{
			discoveryJSON: discoveryStub("", ""),
			success:       true,
			failures:      []ValidationFailure{},
		})
	})

	t.Run("when conditionality checker missingMandatory throws error returns failures", func(t *testing.T) {
		stubIsPresentError := conditionalityCheckerMock{
			isPresent:           true,
			missingMandatory:    []model.Input{},
			missingMandatoryErr: errors.New("the error message"),
		}
		testValidateFailures(t, stubIsPresentError, &invalidTest{
			discoveryJSON: discoveryStub("", ""),
			failures: []ValidationFailure{
				{
					Key:   "DiscoveryModel.DiscoveryItems[0].Endpoints",
					Error: "the error message",
				},
			},
		},
		)
	})

	t.Run("when conditionality checker reports missing mandatory endpoints returns failures", func(t *testing.T) {
		stubMissingMandatory := conditionalityCheckerMock{
			isPresent: true,
			missingMandatory: []model.Input{
				{Method: "GET", Endpoint: "/account-access-consents/{ConsentId}"},
				{Method: "DELETE", Endpoint: "/account-access-consents/{ConsentId}"},
			},
		}
		testValidateFailures(t, stubMissingMandatory, &invalidTest{
			discoveryJSON: discoveryStub("", ""),
			failures: []ValidationFailure{
				{
					Key:   "DiscoveryModel.DiscoveryItems[0].Endpoints",
					Error: "Missing mandatory endpoint Method='GET', Path='/account-access-consents/{ConsentId}'",
				},
				{
					Key:   "DiscoveryModel.DiscoveryItems[0].Endpoints",
					Error: "Missing mandatory endpoint Method='DELETE', Path='/account-access-consents/{ConsentId}'",
				},
			},
		})
	})

	t.Run("Validation should fail when `manifest` is empty", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{}, &invalidTest{
			discoveryJSON: discoveryStub("manifest", ""),
			failures: []ValidationFailure{
				{
					Key:   "DiscoveryModel.DiscoveryItems[0].APISpecification.Manifest",
					Error: "Field 'DiscoveryModel.DiscoveryItems[0].APISpecification.Manifest' is required",
				},
			},
		})
	})

	t.Run("Validation should fail when `manifest` is a normal string, not URL", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{}, &invalidTest{
			discoveryJSON: discoveryStub("manifest", "some-string"),
			failures: []ValidationFailure{
				{
					Key:   "DiscoveryModel.DiscoveryItems[0].APISpecification.Manifest",
					Error: "Field 'DiscoveryModel.DiscoveryItems[0].APISpecification.Manifest' must be 'file://' or 'https://'",
				},
			},
		})
	})

	t.Run("Validation should pass when `manifest` is a https URL", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{}, &invalidTest{
			discoveryJSON: discoveryStub("manifest", "https://www.example.com"),
			failures: []ValidationFailure{
				{
					Key:   "DiscoveryModel.DiscoveryItems[0].Endpoints[0]",
					Error: "Invalid endpoint Method='POST', Path='/account-access-consents'",
				},
				{
					Key:   "DiscoveryModel.DiscoveryItems[0].Endpoints[1]",
					Error: "Invalid endpoint Method='GET', Path='/accounts/{AccountId}/balances'",
				},
			},
		})
	})

	t.Run("Validation should fail when `manifest` is a http URL instead of https", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{}, &invalidTest{
			discoveryJSON: discoveryStub("manifest", "http://www.example.com"),
			failures: []ValidationFailure{
				{
					Key:   "DiscoveryModel.DiscoveryItems[0].APISpecification.Manifest",
					Error: "Field 'DiscoveryModel.DiscoveryItems[0].APISpecification.Manifest' must be 'file://' or 'https://'",
				},
			},
		})
	})
}

func TestDiscovery_Version(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(Version(), "v0.4.0")
}
