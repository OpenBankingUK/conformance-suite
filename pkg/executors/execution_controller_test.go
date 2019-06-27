package executors

import (
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/mocks"
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
