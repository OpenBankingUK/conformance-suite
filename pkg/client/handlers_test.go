package client

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMsgProcessor(t *testing.T) {
	called := false
	myHandler := func(msg []byte) error {
		if string(msg) == "hello" {
			called = true
		}
		return nil
	}
	processor := newMsgProcessor([]msgHandlerFunc{myHandler})

	err := processor.process([]byte("hello"))

	assert.NoError(t, err)
	assert.True(t, called)
}

func TestMsgHandlersChain(t *testing.T) {
	resultChan := make(chan TestCase, 1)
	endedChan := make(chan struct{}, 1)
	processor := newMsgProcessor(msgHandlersChain(resultChan, endedChan))

	test.WithTimeout(t, time.Second, func(t *testing.T) {
		t.Run("a test case result should send a message on resultChan", func(t *testing.T) {
			msg := []byte(`{"type": "ResultType_TestCaseResult", "test": {}}`)
			err := processor.process(msg)
			assert.NoError(t, err)
			assert.Equal(t, TestCase{}, <-resultChan)
		})

		t.Run("a test cases completed should send a message on endedChan", func(t *testing.T) {
			msg := []byte(`{"type": "ResultType_TestCasesCompleted"}`)
			err := processor.process(msg)
			assert.NoError(t, err)
			assert.Equal(t, struct{}{}, <-endedChan)
		})
	})
}
