package executors

import (
	"testing"

	"github.com/OpenBankingUK/conformance-suite/pkg/test"

	"github.com/OpenBankingUK/conformance-suite/pkg/executors/mocks"
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
