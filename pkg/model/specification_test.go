package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func specsFixture(values map[string]string) []byte {
	schemaVersion, ok := values["schemaVersion"]
	if !ok {
		schemaVersion = "https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json"
	}
	url, ok := values["url"]
	if !ok {
		url = "https://openbanking.atlassian.net/wiki/spaces/DZ/pages/642090641/Account+and+Transaction+API+Specification+-+v3.0"
	}
	version, ok := values["version"]
	if !ok {
		version = "v3.0"
	}
	return []byte(`[
		  {
		    "identifier": "account-transaction-v3.0",
		    "name": "Account and Transaction API Specification",
		    "url": "` + url + `",
		    "version": "` + version + `",
		    "schemaVersion": "` + schemaVersion + `"
		  }
		]`)
}

func assertSpecLoadError(t *testing.T, field string, value string, expected string) {
	config := specsFixture(map[string]string{field: value})
	err := loadSpecifications(config)
	assert.Equal(t, expected, err.Error())
}

func TestLoadSpecifications_invalid(t *testing.T) {
	defer loadDefaultSpecifications() // set config back to a valid state for next test

	expected := "Key: 'Specification.SchemaVersion' Error:Field validation for 'SchemaVersion' failed on the 'url' tag"
	assertSpecLoadError(t, "schemaVersion", "invalid-url", expected)

	expected = "Key: 'Specification.URL' Error:Field validation for 'URL' failed on the 'url' tag"
	assertSpecLoadError(t, "url", "invalid-url", expected)

	expected = "Key: 'Specification.URL' Error:Field validation for 'URL' failed on the 'required' tag"
	assertSpecLoadError(t, "url", "", expected)

	expected = "Key: 'Specification.Version' Error:Field validation for 'Version' failed on the 'required' tag"
	assertSpecLoadError(t, "version", "", expected)
}

func TestSpecificationIdentifierFromSchemaVersion(t *testing.T) {
	t.Run("returns specification identifier when given valid schema version URL", func(t *testing.T) {
		specification, err := SpecificationIdentifierFromSchemaVersion("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json")
		assert.Nil(t, err)
		assert.Equal(t, specification, "account-transaction-v3.0")
	})

	t.Run("returns error when given invalid schema version URL", func(t *testing.T) {
		schemaVersion := "https://example.com/invalid"
		specification, err := SpecificationIdentifierFromSchemaVersion(schemaVersion)
		assert.Equal(t, err.Error(), "No specification found for schema version: "+schemaVersion)
		assert.Equal(t, specification, "")
	})
}
