package executors

import (
	"errors"
	"testing"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/results"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
)

const (
	selectTimeout = 1 * time.Millisecond
)

func TestNewDaemonController(t *testing.T) {
	assert := test.NewAssert(t)

	testChan := make(chan results.TestCase, 100)
	controller := NewDaemonController(testChan)

	assert.NotNil(testChan, controller.Results())
	assert.NotNil(controller.stopLock)
	assert.False(controller.shouldStop)
	assert.False(controller.ShouldStop())
	assert.Empty(controller.AllResults())
}

func TestNewBufferedDaemonController(t *testing.T) {
	assert := test.NewAssert(t)

	controller := NewBufferedDaemonController()

	assert.NotNil(controller.Results())
	assert.NotNil(controller.stopLock)
	assert.False(controller.shouldStop)
	assert.False(controller.ShouldStop())
}

func TestDaemonControllerStops(t *testing.T) {
	assert := test.NewAssert(t)

	testChan := make(chan results.TestCase, 100)
	controller := NewDaemonController(testChan)

	controller.Stop()

	assert.True(controller.shouldStop)
	assert.True(controller.ShouldStop())
}

func TestNewBufferedDaemonControllerResults(t *testing.T) {
	require := test.NewAssert(t)

	controller := NewBufferedDaemonController()

	require.NotNil(controller.Results())
	// initially empty
	select {
	case msg, ok := <-controller.Results():
		require.Nil(msg)
		require.False(ok)
	case <-time.After(selectTimeout):
		break
	}
}

func TestNewBufferedDaemonControllerAllResults(t *testing.T) {
	require := test.NewAssert(t)

	controller := NewBufferedDaemonController()

	require.Empty(controller.AllResults())

	err := errors.New("some error")
	result := results.NewTestCaseResult("123", true, results.NoMetrics(), []error{err}, "endpoint", "api-name", "api-version", "detailed description", "https://openbanking.org.uk/ref/uri", "200")
	controller.AddResult(result)

	// channel contains event
	select {
	case msg, ok := <-controller.Results():
		require.Equal(result, msg)
		require.True(ok)
	case <-time.After(selectTimeout):
		break
	}

	// result has been accumulated
	require.Equal([]results.TestCase{
		result,
	}, controller.AllResults())
}

func TestNewBufferedDaemonControllerSetCompletedAndIsCompleted(t *testing.T) {
	require := test.NewAssert(t)

	controller := NewBufferedDaemonController()

	require.NotNil(controller.IsCompleted())

	// initially not completed
	select {
	case msg, ok := <-controller.IsCompleted():
		require.Nil(msg)
		require.False(ok)
	case <-time.After(selectTimeout):
		break
	}

	// mark as completed
	controller.SetCompleted()

	// should receive completed event on channel
	select {
	case msg, ok := <-controller.IsCompleted():
		require.True(msg)
		require.True(ok)
	case <-time.After(selectTimeout):
		break
	}
}
