package server

import (
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
	"github.com/pkg/errors"
)

func TestNewErrorResponse(t *testing.T) {
	assert := test.NewAssert(t)

	errorResponse := NewErrorResponse(errors.New("standard error"))

	assert.Equal("standard error", errorResponse.Error)
}

func TestNewErrorResponsePreservesStackTrace(t *testing.T) {
	assert := test.NewAssert(t)

	deepError := errors.New("deepError error")
	middleError := errors.Wrap(deepError, "middleError error")
	lastError := errors.Wrap(middleError, "lastError error")
	errorResponse := NewErrorResponse(lastError)

	assert.Equal("lastError error: middleError error: deepError error", errorResponse.Error)
}
