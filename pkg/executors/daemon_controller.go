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

	AddResult(result results.TestCase)
	AllResults() []results.TestCase
	AllResultsGrouped() map[results.ResultKey][]results.TestCase
	AddResponseFields(string)
	ResponseFieldsJSON() string

	Results() <-chan results.TestCase

	SetCompleted()
	IsCompleted() <-chan bool
}

// daemonController manages routine running tests
// allowing to stop and collect results/errors
type daemonController struct {
	results         []results.TestCase
	resultsGrouped  map[results.ResultKey][]results.TestCase
	resultChan      chan results.TestCase
	responseFields  string
	stopLock        *sync.Mutex
	shouldStop      bool
	isCompletedChan chan bool
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
		results:         []results.TestCase{},
		resultChan:      resultChan,
		stopLock:        &sync.Mutex{},
		shouldStop:      false,
		isCompletedChan: make(chan bool, 1),
		resultsGrouped:  make(map[results.ResultKey][]results.TestCase),
	}
}

// Stop tell the daemon to stop
func (rc *daemonController) Stop() {
	rc.stopLock.Lock()
	defer rc.stopLock.Unlock()
	rc.shouldStop = true
}

// Stopped tell the daemon service has stopped
func (rc *daemonController) Stopped() {
	rc.stopLock.Lock()
	defer rc.stopLock.Unlock()
	rc.shouldStop = false
}

// ShouldStop indicates that the daemon should stop
// this should be invoked often by the background routine and stop
// if this true
func (rc *daemonController) ShouldStop() bool {
	rc.stopLock.Lock()
	defer rc.stopLock.Unlock()
	shouldStop := rc.shouldStop
	return shouldStop
}

// AddResult - add result.
func (rc *daemonController) AddResult(result results.TestCase) {
	rc.results = append(rc.results, result)
	mpKey := results.ResultKey{
		APIVersion: result.APIVersion,
		APIName:    result.API,
	}
	if _, ok := rc.resultsGrouped[mpKey]; !ok {
		rc.resultsGrouped[mpKey] = make([]results.TestCase, 0)
	}
	rc.resultsGrouped[mpKey] = append(rc.resultsGrouped[mpKey], result)
	rc.resultChan <- result
}

// AllResults - returns all the accumulated results.
func (rc *daemonController) AllResults() []results.TestCase {
	return rc.results
}

// AllResultsGrouped - returns all the accumulated results Grouped by the type `ResultKey`.
func (rc *daemonController) AllResultsGrouped() map[results.ResultKey][]results.TestCase {
	return rc.resultsGrouped
}

func (rc *daemonController) AddResponseFields(f string) {
	rc.responseFields = f
}

func (rc *daemonController) ResponseFieldsJSON() string {
	return rc.responseFields
}

// ResultsChannel - return channel for receiving results.
func (rc *daemonController) Results() <-chan results.TestCase {
	return rc.resultChan
}

// SetCompleted - mark the tests as completed.
func (rc *daemonController) SetCompleted() {
	rc.isCompletedChan <- true
}

// IsCompleted - channel to subscribe to completed event.
func (rc *daemonController) IsCompleted() <-chan bool {
	return rc.isCompletedChan
}
