package executors

import (
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	"testing"

	mocks2 "bitbucket.org/openbankingteam/conformance-suite/pkg/authentication/mocks"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/mocks"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestNewTestCaseRunner(t *testing.T) {
	controller := &mocks.DaemonController{}
	definition := RunDefinition{}
	runner := NewTestCaseRunner(test.NullLogger(), definition, controller)

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
			DiscoveryModel: discovery.ModelDiscovery{},
		},
	}
	runner := NewTestCaseRunner(test.NullLogger(), definition, controller)

	ctx := runner.makeRuleCtx(&model.Context{})

	value, ok := ctx.Get("SigningCert")
	assert.True(t, ok)
	assert.Equal(t, cert, value)
}
