package executors

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/reporting"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAllPass(t *testing.T) {
	testResults := []reporting.Test{
		{Pass: true},
		{Pass: true},
		{Pass: true},
	}

	assert.True(t, allPass(testResults))
}

func TestAllPassFalse(t *testing.T) {
	testResults := []reporting.Test{
		{Pass: true},
		{Pass: false},
		{Pass: true},
	}

	assert.False(t, allPass(testResults))
}

func TestAllPassEmpty(t *testing.T) {
	testResults := []reporting.Test{}

	assert.True(t, allPass(testResults))
}

func TestMakeSpecResult(t *testing.T) {
	tests := []reporting.Test{{Pass: false}}
	spec := discovery.ModelAPISpecification{
		Name:          "Name",
		SchemaVersion: "SchemaVersion",
		URL:           "URL",
		Version:       "Version",
	}

	reportSpec := makeSpecResult(spec, tests)

	assert.Equal(t, "Name", reportSpec.Name)
	assert.Equal(t, "SchemaVersion", reportSpec.SchemaVersion)
	assert.Equal(t, "URL", reportSpec.URL)
	assert.Equal(t, "Version", reportSpec.Version)
	assert.Len(t, reportSpec.Tests, 1)
	assert.False(t, reportSpec.Pass)
}
