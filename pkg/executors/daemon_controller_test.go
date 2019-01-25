package executors

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/results"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDaemonController(t *testing.T) {
	testChan := make(chan results.Test, 100)
	errorChan := make(chan error, 100)

	controller := NewDaemonController(testChan, errorChan)

	assert.Equal(t, testChan, controller.Results())
	assert.Equal(t, errorChan, controller.Errors())
	assert.NotNil(t, controller.mx)
	assert.False(t, controller.shouldStop)
	assert.False(t, controller.ShouldStop())
}

func TestDaemonControllerStops(t *testing.T) {
	testChan := make(chan results.Test, 100)
	errorChan := make(chan error, 100)
	controller := NewDaemonController(testChan, errorChan)

	controller.Stop()

	assert.True(t, controller.shouldStop)
	assert.True(t, controller.ShouldStop())
}
