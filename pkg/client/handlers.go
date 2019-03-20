package client

import (
	"encoding/json"
	"github.com/pkg/errors"
)

// msgProcessor passes all messages/event through a list of handler
// so each has an opportunity to handle that message
type msgProcessor struct {
	handlers []msgHandlerFunc
}

func newMsgProcessor(handlers []msgHandlerFunc) msgProcessor {
	return msgProcessor{handlers: handlers}
}

// process just passes msg into all handler, aborts if any errors
func (p msgProcessor) process(msg []byte) error {
	for _, handler := range p.handlers {
		err := handler(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

type msgHandlerFunc func(msg []byte) error

// msgHandlersChain build a list of message handlers
func msgHandlersChain(resultChan chan<- TestCase, endedChan chan<- struct{}) []msgHandlerFunc {
	return []msgHandlerFunc{
		handlerTestCaseResult(resultChan),
		handlerRunEnded(endedChan),
	}
}

// handlerTestCaseResult looks for messages with TestCaseResult schema and
// returns the result thru a channel if its a result message
func handlerTestCaseResult(resultChan chan<- TestCase) msgHandlerFunc {
	return func(msg []byte) error {
		tcResult := TestCaseResult{}
		err := json.Unmarshal(msg, &tcResult)
		if err != nil {
			return errors.Wrap(err, "test case result handler")
		}
		if tcResult.Type == "ResultType_TestCaseResult" {
			resultChan <- tcResult.Test
		}
		return nil
	}
}

type TestCaseResult struct {
	Type string   `json:"type"`
	Test TestCase `json:"test"`
}

// TestCase result for a run
type TestCase struct {
	Id   string `json:"id"`
	Pass bool   `json:"pass"`
	Fail string `json:"fail,omitempty"`
}

type event struct {
	Type string `json:"type"`
}

// handlerRunEnded checks if a message is a ended running test and if so signals the ended channel
func handlerRunEnded(endedChan chan<- struct{}) msgHandlerFunc {
	return func(msg []byte) error {
		aEvent := event{}
		err := json.Unmarshal(msg, &aEvent)
		if err != nil {
			return err
		}
		if aEvent.Type == "ResultType_TestCasesCompleted" {
			endedChan <- struct{}{}
		}
		return nil
	}
}
