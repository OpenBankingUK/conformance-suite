package reporting

import (
	"encoding/json"
	"flag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "update .golden files")

func TestResultsMarshalHasNotChanged(t *testing.T) {
	result := Result{
		Specifications: []Specification{
			{
				Name:          "spec name",
				Version:       "spec version",
				SchemaVersion: "spec schema version",
				URL:           "url",
				Pass:          true,
				Tests: []Test{
					{
						Name:     "test name",
						Id:       "test id",
						Endpoint: "test endpoint",
						Pass:     true,
					},
				},
			},
		},
	}

	expected, err := json.MarshalIndent(result, "", "    ")
	require.NoError(t, err)

	goldenFile := filepath.Join("testdata", "test_results.golden")
	if *update {
		t.Log("update golden file")
		err := ioutil.WriteFile(goldenFile, expected, 0644)
		require.NoError(t, err, "failed to update golden file")
	}

	testResults, err := ioutil.ReadFile(goldenFile)
	require.NoError(t, err, "failed reading .golden")

	assert.JSONEq(t, string(expected), string(testResults))
}
