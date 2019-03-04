//go:generate mockery -name DaemonController
package executors

import (
	"sync"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/results"
)

type DaemonController interface {
	Stop()
	ShouldStop() bool
	Stopped()
	Results() chan results.TestCase
}

// daemonController manages routine running tests
// allowing to stop and collect results/errors
type daemonController struct {
	resultChan chan results.TestCase
	stopLock   *sync.Mutex
	shouldStop bool
}

// NewBufferedDaemonController new instance to control a background routine with 100 objects
// buffer in result and error channels
func NewBufferedDaemonController() *daemonController {
	const chanBufferSize = 100
	return NewDaemonController(make(chan results.TestCase, chanBufferSize))

}

// NewDaemonController new instance to control a background routine
func NewDaemonController(resultChan chan results.TestCase) *daemonController {
	return &daemonController{
		resultChan: resultChan,
		stopLock:   &sync.Mutex{},
		shouldStop: false,
	}
}

// Stop tell the daemon to stop
func (rc *daemonController) Stop() {
	rc.stopLock.Lock()
	rc.shouldStop = true
	rc.stopLock.Unlock()
}

// Stopped tell the daemon service has stopped
func (rc *daemonController) Stopped() {
	rc.stopLock.Lock()
	rc.shouldStop = false
	rc.stopLock.Unlock()
}

// ShouldStop indicates that the daemon should stop
// this should be invoked often by the background routine and stop
// if this true
func (rc *daemonController) ShouldStop() bool {
	rc.stopLock.Lock()
	shouldStop := rc.shouldStop
	rc.stopLock.Unlock()
	return shouldStop
}

func (rc *daemonController) Results() chan results.TestCase {
	return rc.resultChan
}
