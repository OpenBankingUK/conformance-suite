package executors

import (
	mocks2 "bitbucket.org/openbankingteam/conformance-suite/pkg/authentication/mocks"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTestCaseRunner(t *testing.T) {
	controller := &mocks.DaemonController{}
	definition := RunDefinition{}
	runner := NewTestCaseRunner(definition, controller)

	assert.NotNil(t, runner.mux)
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
