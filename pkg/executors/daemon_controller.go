package executors

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/results"
	"sync"
)

type DaemonController interface {
	Stop()
	ShouldStop() bool
	Results() chan results.TestCase
	Errors() chan error
}

// daemonController manages routine running tests
// allowing to stop and collect results/errors
type daemonController struct {
	resultChan chan results.TestCase
	errorsChan chan error
	stopLock   *sync.Mutex
	shouldStop bool
}

// NewDaemonController new instance to control a background routine
func NewDaemonController(resultChan chan results.TestCase, errorsChan chan error) *daemonController {
	return &daemonController{
		resultChan: resultChan,
		errorsChan: errorsChan,
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

func (rc *daemonController) Errors() chan error {
	return rc.errorsChan
}
