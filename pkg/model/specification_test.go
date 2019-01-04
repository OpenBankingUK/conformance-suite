package model

import (
	"encoding/json"
	"flag"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "update .golden files")

func TestSpecificationIdentifierFromSchemaVersion(t *testing.T) {
	t.Run("returns specifications identifier when given valid schema version URL", func(t *testing.T) {
		specification, err := SpecificationFromSchemaVersion("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json")
		require.NoError(t, err)
		assert.Equal(t, specification.Identifier, "account-transaction-v3.0")
	})

	t.Run("returns error when given invalid schema version URL", func(t *testing.T) {
		schemaVersion := "https://example.com/invalid"
		specification, err := SpecificationFromSchemaVersion(schemaVersion)
		require.Error(t, err)
		assert.EqualError(t, err, "no specifications found for schema version: "+schemaVersion)
		assert.Equal(t, specification.Identifier, "")
	})
}

func TestOzoneSpecificationHasNotChanged(t *testing.T) {
	expected, err := json.MarshalIndent(specifications, "", "    ")
	require.NoError(t, err)

	goldenFile := filepath.Join("testdata", "ozone_spec.golden")
	if *update {
		t.Log("update golden file")
		err := ioutil.WriteFile(goldenFile, expected, 0644)
		require.NoError(t, err, "failed to update golden file")
	}

	spec, err := ioutil.ReadFile(goldenFile)
	require.NoError(t, err, "failed reading .golden")

	assert.JSONEq(t, string(expected), string(spec))
}
