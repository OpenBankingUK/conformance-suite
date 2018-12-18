package server

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewErrorResponse(t *testing.T) {
	errorResponse := NewErrorResponse(errors.New("standard error"))

	assert.Equal(t, "standard error", errorResponse.Error)
}

func TestNewErrorResponsePreservesStackTrace(t *testing.T) {
	deepError := errors.New("deepError error")
	middleError := errors.Wrap(deepError, "middleError error")
	lastError := errors.Wrap(middleError, "lastError error")
	errorResponse := NewErrorResponse(lastError)

	assert.Equal(t, "lastError error: middleError error: deepError error", errorResponse.Error)
}
