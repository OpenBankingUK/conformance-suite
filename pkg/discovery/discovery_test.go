package discovery

import (
	"io/ioutil"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"

	"github.com/stretchr/testify/assert"
)

// invalidTestCase
// name - name of the test case.
// config - the discovery model config.
// err - the expected err
type invalidTestCase struct {
	name        string
	config      string
	expectedErr string
}

// invalidTestCase
// discoveryJSON - the discovery model JSON
// failures - list of failures
// err - the expected err
type invalidTest struct {
	discoveryJSON string
	success       bool
	failures      []string
	err           error
}

// conditionalityCheckerMock - implements model.ConditionalityChecker interface for tests
type conditionalityCheckerMock struct {
	isPresent        bool
	missingMandatory []model.Input
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
	return c.isPresent, nil
}

// MissingMandatory - returns stubbed array of missing endpoints
func (c conditionalityCheckerMock) MissingMandatory(endpoints []model.Input, specification string) ([]model.Input, error) {
	return c.missingMandatory, nil
}

// unmarshalDiscoveryJSON - returns discovery model
func testUnmarshalDiscoveryJSON(t *testing.T, discoveryJSON string) *Model {
	discovery, err := unmarshalDiscoveryJSON(discoveryJSON)
	assert.NoError(t, err)
	return discovery
}

// discoveryStub - returns discovery JSON with given field stubbed with given value
func discoveryStub(field string, value string) string {
	version := "v0.0.1"
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
	}

	discoveryItems := `, "discoveryItems": [
		{
			"apiSpecification": {
				"name": "Account and Transaction API Specification",
				"url": "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.0",
				"version": "v3.0",
				"schemaVersion": "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json"
			},
			"openidConfigurationUri": "https://as.aspsp.ob.forgerock.financial/oauth2/.well-known/openid-configuration",
			"resourceBaseUri": "https://rs.aspsp.ob.forgerock.financial:443/"` + endpoints + `
		}
	]`
	if field == "discoveryItems" {
		if value == "" {
			discoveryItems = ""
		} else {
			discoveryItems = `, "discoveryItems": ` + value
		}
	}
	return `
		{
			"discoveryModel": {
				"version": "` + version + `"` +
		discoveryItems + `
			}
	}`
}

// testValidateFailures - run Validate, and test validation failure expectations
func testValidateFailures(t *testing.T, checker model.ConditionalityChecker, expected *invalidTest) {
	discovery := testUnmarshalDiscoveryJSON(t, expected.discoveryJSON)
	result, failures, err := Validate(checker, discovery)
	assert.Equal(t, expected.success, result)
	assert.Equal(t, expected.err, err)
	assert.Equal(t, expected.failures, failures)
}

// TestValidate - test Validate function
func TestValidate(t *testing.T) {

	t.Run("when version missing returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{}, &invalidTest{
			discoveryJSON: discoveryStub("version", ""),
			failures: []string{
				`Key: 'Model.DiscoveryModel.Version' Error:Field validation for 'Version' failed on the 'required' tag`,
			},
		})
	})

	t.Run("when version not in discovery.SupportedVersions() returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{isPresent: true}, &invalidTest{
			discoveryJSON: discoveryStub("version", "v9.9.9"),
			failures: []string{
				`Key: 'Model.DiscoveryModel.Version' Error:Version v9.9.9 not in list of supported versions`,
			},
		})
	})

	t.Run("when discoveryItems missing returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{}, &invalidTest{
			discoveryJSON: discoveryStub("discoveryItems", ""),
			failures: []string{
				`Key: 'Model.DiscoveryModel.DiscoveryItems' Error:Field validation for 'DiscoveryItems' failed on the 'required' tag`,
			},
		})
	})

	t.Run("when discoveryItems is empty array returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{}, &invalidTest{
			discoveryJSON: discoveryStub("discoveryItems", "[]"),
			failures: []string{
				`Key: 'Model.DiscoveryModel.DiscoveryItems' Error:Field validation for 'DiscoveryItems' failed on the 'gt' tag`,
			},
		})
	})

	t.Run("when discoveryItems has empty endpoints array returns failure", func(t *testing.T) {
		testValidateFailures(t, conditionalityCheckerMock{}, &invalidTest{
			discoveryJSON: discoveryStub("endpoints", "[]"),
			failures: []string{
				`Key: 'Model.DiscoveryModel.DiscoveryItems[0].Endpoints' Error:Field validation for 'Endpoints' failed on the 'gt' tag`,
			},
		})
	})

	t.Run("when conditionality checker reports endpoints not present returns failures", func(t *testing.T) {
		stubAllNotPresent := conditionalityCheckerMock{isPresent: false}
		testValidateFailures(t, stubAllNotPresent, &invalidTest{
			discoveryJSON: discoveryStub("", ""),
			failures: []string{
				`discoveryItemIndex=0, invalid endpoint Method=POST, Path=/account-access-consents`,
				`discoveryItemIndex=0, invalid endpoint Method=GET, Path=/accounts/{AccountId}/balances`,
			},
		})
	})

	t.Run("when conditionality checker reports endpoints present returns no failures", func(t *testing.T) {
		stubAllPresent := conditionalityCheckerMock{isPresent: true}
		testValidateFailures(t, stubAllPresent, &invalidTest{
			discoveryJSON: discoveryStub("", ""),
			success:       true,
			failures:      []string{},
		})
	})

	t.Run("when conditionality checker reports missing mandatory endpoints returns failures", func(t *testing.T) {
		stubMissingMandatory := conditionalityCheckerMock{
			isPresent: true,
			missingMandatory: []model.Input{
				model.Input{Method: "GET", Endpoint: "/account-access-consents/{ConsentId}"},
				model.Input{Method: "DELETE", Endpoint: "/account-access-consents/{ConsentId}"},
			},
		}
		testValidateFailures(t, stubMissingMandatory, &invalidTest{
			discoveryJSON: discoveryStub("", ""),
			failures: []string{
				`discoveryItemIndex=0, missing mandatory endpoint Method=GET, Path=/account-access-consents/{ConsentId}`,
				`discoveryItemIndex=0, missing mandatory endpoint Method=DELETE, Path=/account-access-consents/{ConsentId}`,
			},
		})
	})
}

func TestDiscovery_FromJSONString_Valid(t *testing.T) {
	assert := assert.New(t)

	discoveryExample, err := ioutil.ReadFile("../../docs/discovery-example.json")
	assert.NoError(err)
	assert.NotNil(discoveryExample)
	config := string(discoveryExample)

	accountAPIDiscoveryItem := ModelDiscoveryItem{
		APISpecification: ModelAPISpecification{
			Name:          "Account and Transaction API Specification",
			URL:           "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.0",
			Version:       "v3.0",
			SchemaVersion: "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json",
		},
		OpenidConfigurationURI: "https://as.aspsp.ob.forgerock.financial/oauth2/.well-known/openid-configuration",
		ResourceBaseURI:        "https://rs.aspsp.ob.forgerock.financial:443/",
		Endpoints: []ModelEndpoint{
			ModelEndpoint{
				Method:                "POST",
				Path:                  "/account-access-consents",
				ConditionalProperties: []ModelConditionalProperties(nil),
			},
			ModelEndpoint{
				Method:                "GET",
				Path:                  "/account-access-consents/{ConsentId}",
				ConditionalProperties: []ModelConditionalProperties(nil),
			},
			ModelEndpoint{Method: "DELETE",
				Path:                  "/account-access-consents/{ConsentId}",
				ConditionalProperties: []ModelConditionalProperties(nil),
			},
			ModelEndpoint{Method: "GET",
				Path:                  "/accounts/{AccountId}/product",
				ConditionalProperties: []ModelConditionalProperties(nil),
			},
			ModelEndpoint{Method: "GET",
				Path: "/accounts/{AccountId}/transactions",
				ConditionalProperties: []ModelConditionalProperties{
					ModelConditionalProperties{
						Schema:   "OBTransaction3Detail",
						Property: "Balance",
						Path:     "Data.Transaction[*].Balance",
					},
					ModelConditionalProperties{
						Schema:   "OBTransaction3Detail",
						Property: "MerchantDetails",
						Path:     "Data.Transaction[*].MerchantDetails",
					},
					ModelConditionalProperties{
						Schema:   "OBTransaction3Basic",
						Property: "TransactionReference",
						Path:     "Data.Transaction[*].TransactionReference",
					},
					ModelConditionalProperties{
						Schema:   "OBTransaction3Detail",
						Property: "TransactionReference",
						Path:     "Data.Transaction[*].TransactionReference",
					},
				},
			},
			ModelEndpoint{
				Method:                "GET",
				Path:                  "/accounts",
				ConditionalProperties: []ModelConditionalProperties(nil),
			},
			ModelEndpoint{
				Method:                "GET",
				Path:                  "/accounts/{AccountId}",
				ConditionalProperties: []ModelConditionalProperties(nil),
			},
			ModelEndpoint{
				Method:                "GET",
				Path:                  "/accounts/{AccountId}/balances",
				ConditionalProperties: []ModelConditionalProperties(nil),
			},
		},
	}

	modelActual, err := unmarshalDiscoveryJSON(config)
	assert.NoError(err)
	assert.NotNil(modelActual.DiscoveryModel)
	discoveryModel := modelActual.DiscoveryModel

	t.Run("model has a version", func(t *testing.T) {
		assert.Equal(discoveryModel.Version, "v0.0.1")
	})

	t.Run("model has correct number of discovery items", func(t *testing.T) {
		assert.Equal(len(discoveryModel.DiscoveryItems), 2)
	})

	t.Run("model has correct discovery item contents", func(t *testing.T) {
		assert.Equal(accountAPIDiscoveryItem, discoveryModel.DiscoveryItems[0])
	})
}

func TestDiscovery_Version(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(Version(), "v0.0.1")
}
