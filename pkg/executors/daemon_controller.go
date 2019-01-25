package executors

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/results"
	"sync"
)

type DaemonController interface {
	Stop()
	ShouldStop() bool
	Results() chan results.Test
	Errors() chan error
}

// daemonController manages routine running tests
// allowing to stop and collect results/errors
type daemonController struct {
	resultChan chan results.Test
	errorsChan chan error
	mx         *sync.Mutex
	shouldStop bool
}

// NewDaemonController new instance to control a background routine
func NewDaemonController(resultChan chan results.Test, errorsChan chan error) *daemonController {
	return &daemonController{
		resultChan: resultChan,
		errorsChan: errorsChan,
		mx:         &sync.Mutex{},
		shouldStop: false,
	}
}

// Stop tell the daemon to stop
func (rc *daemonController) Stop() {
	rc.mx.Lock()
	rc.shouldStop = true
	rc.mx.Unlock()
}

// ShouldStop indicates that the daemon should stop
// this should be invoked often my the background routine and stop
// if this true
func (rc *daemonController) ShouldStop() bool {
	rc.mx.Lock()
	shouldStop := rc.shouldStop
	rc.mx.Unlock()
	return shouldStop
}

func (rc *daemonController) Results() chan results.Test {
	return rc.resultChan
}

func (rc *daemonController) Errors() chan error {
	return rc.errorsChan
}
