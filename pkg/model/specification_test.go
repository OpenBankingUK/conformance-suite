package model

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update .golden files")

func TestSpecificationIdentifierFromSchemaVersion(t *testing.T) {
	t.Run("returns specifications identifier when given valid schema version URL", func(t *testing.T) {
		specification, err := SpecificationFromSchemaVersion("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.1.0/dist/account-info-swagger.json")
		require.NoError(t, err)
		assert.Equal(t, "account-transaction-v3.1", specification.Identifier)
	})

	t.Run("returns error when given invalid schema version URL", func(t *testing.T) {
		schemaVersion := "https://example.com/invalid"
		specification, err := SpecificationFromSchemaVersion(schemaVersion)
		require.Error(t, err)
		assert.EqualError(t, err, "no specifications found for schema version: "+schemaVersion)
		assert.Equal(t, "", specification.Identifier)
	})
}

func TestSpecificationHasNotChanged(t *testing.T) {
	expected, err := json.MarshalIndent(specifications, "", "    ")
	require.NoError(t, err)
	goldenFile := filepath.Join("testdata", "spec-config.golden.json")

	if *update {
		t.Log("update golden file")
		require.NoError(t, ioutil.WriteFile(goldenFile, expected, 0644), "failed to update golden file")
	}

	spec, err := ioutil.ReadFile(goldenFile)
	require.NoError(t, err, "failed reading .golden")
	assert.JSONEq(t, string(expected), string(spec))
}
