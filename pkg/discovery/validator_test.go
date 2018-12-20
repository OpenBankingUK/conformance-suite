package discovery

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWrapperValidateCallsValidatorFunWithChecker(t *testing.T) {
	calledWith := false
	conditionalityChecker := model.NewConditionalityChecker()
	validator := funcWrapperValidator{
		validatorFunc: func(checker model.ConditionalityChecker, discovery *Model) (bool, []ValidationFailure, error) {
			if checker == conditionalityChecker {
				calledWith = true
			}
			return true, []ValidationFailure{{}}, nil
		},
		conditionalityChecker: conditionalityChecker,
	}

	_, err := validator.Validate(&Model{})

	require.NoError(t, err)
	assert.True(t, calledWith)
}

func TestWrapperValidateCastsToValidationFailures(t *testing.T) {
	validator := funcWrapperValidator{
		validatorFunc: func(checker model.ConditionalityChecker, discovery *Model) (bool, []ValidationFailure, error) {
			return true, []ValidationFailure{{}}, nil
		},
	}

	failures, err := validator.Validate(&Model{})

	require.NoError(t, err)
	assert.IsType(t, ValidationFailures{}, failures)
}

func TestWrapperValidateReturnsNoValidationFailures(t *testing.T) {
	validator := funcWrapperValidator{
		validatorFunc: func(checker model.ConditionalityChecker, discovery *Model) (bool, []ValidationFailure, error) {
			return false, nil, nil
		},
	}

	failures, err := validator.Validate(&Model{})

	require.NoError(t, err)
	assert.Equal(t, NoValidationFailures, failures)
}

func TestWrapperValidateHandlesErrors(t *testing.T) {
	validator := funcWrapperValidator{
		validatorFunc: func(checker model.ConditionalityChecker, discovery *Model) (bool, []ValidationFailure, error) {
			return false, nil, errors.New("some error")
		},
	}

	failures, err := validator.Validate(&Model{})

	assert.Nil(t, failures)
	assert.Error(t, err)
}

func TestValidatorNoFailuresReturnsEmpty(t *testing.T) {
	failures := ValidationFailures{}

	assert.True(t, failures.Empty())
}

func TestValidatorFailuresReturnsNotEmpty(t *testing.T) {
	failures := ValidationFailures{ValidationFailure{}}

	assert.False(t, failures.Empty())
}

func TestValidatorNoValidationFailuresReturnsEmpty(t *testing.T) {
	assert.True(t, NoValidationFailures.Empty())
}
