package executors

import (
	mocks2 "bitbucket.org/openbankingteam/conformance-suite/pkg/authentication/mocks"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/mocks"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/permissions"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTestCaseRunner(t *testing.T) {
	controller := &mocks.DaemonController{}
	definition := RunDefinition{}
	runner := NewTestCaseRunner(definition, controller)

	assert.NotNil(t, runner.runningLock)
	assert.Equal(t, definition, runner.definition)
	assert.Equal(t, controller, runner.daemonController)
	assert.False(t, runner.running)
}

func TestMakeRuleContext(t *testing.T) {
	controller := &mocks.DaemonController{}
	cert := &mocks2.Certificate{}
	definition := RunDefinition{
		SigningCert: cert,
		DiscoModel: &discovery.Model{
			DiscoveryModel: discovery.ModelDiscovery{
				CustomTests: []discovery.CustomTest{
					{
						Replacements: map[string]string{"key": "value"},
					},
				},
			},
		},
	}
	runner := NewTestCaseRunner(definition, controller)

	ctx := runner.makeRuleCtx()

	value, ok := ctx.Get("SigningCert")
	assert.True(t, ok)
	assert.Equal(t, cert, value)

	value, ok = ctx.Get("key")
	assert.True(t, ok)
	assert.Equal(t, "value", value)
}

func TestPermissionsSetsEmpty(t *testing.T) {
	controller := &mocks.DaemonController{}
	definition := RunDefinition{}
	runner := NewTestCaseRunner(definition, controller)

	results := runner.permissionSets()

	assert.Len(t, results, 0)
}

func TestPermissionsShouldPassAllTestsToResolver(t *testing.T) {
	controller := &mocks.DaemonController{}
	definition := RunDefinition{
		SpecTests: []generation.SpecificationTestCases{
			{
				TestCases: []model.TestCase{
					{ID: "1"}, {ID: "2"},
				},
			},
			{
				TestCases: []model.TestCase{
					{ID: "3"},
				},
			},
		},
	}
	runner := NewTestCaseRunner(definition, controller)
	var resultsGroups []permissions.Group
	called := false
	runner.resolver = func(groups []permissions.Group) permissions.CodeSetResultSet {
		resultsGroups = groups
		called = true
		return nil
	}

	runner.permissionSets()

	assert.True(t, called)
	assert.Len(t, resultsGroups, 3)
	assert.Equal(t, "1", string(resultsGroups[0].TestId))
	assert.Equal(t, "2", string(resultsGroups[1].TestId))
	assert.Equal(t, "3", string(resultsGroups[2].TestId))
}
