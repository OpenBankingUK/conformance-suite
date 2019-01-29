package executors

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/results"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDaemonController(t *testing.T) {
	testChan := make(chan results.TestCase, 100)
	errorChan := make(chan error, 100)

	controller := NewDaemonController(testChan, errorChan)

	assert.Equal(t, testChan, controller.Results())
	assert.Equal(t, errorChan, controller.Errors())
	assert.NotNil(t, controller.stopLock)
	assert.False(t, controller.shouldStop)
	assert.False(t, controller.ShouldStop())
}

func TestNewBufferedDaemonController(t *testing.T) {
	controller := NewBufferedDaemonController()

	assert.NotNil(t, controller.Results())
	assert.NotNil(t, controller.Errors())
	assert.NotNil(t, controller.stopLock)
	assert.False(t, controller.shouldStop)
	assert.False(t, controller.ShouldStop())
}

func TestDaemonControllerStops(t *testing.T) {
	testChan := make(chan results.TestCase, 100)
	errorChan := make(chan error, 100)
	controller := NewDaemonController(testChan, errorChan)

	controller.Stop()

	assert.True(t, controller.shouldStop)
	assert.True(t, controller.ShouldStop())
}
