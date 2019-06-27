package discovery

import (
	"errors"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrapperValidateCallsValidatorFuncWithChecker(t *testing.T) {
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
	assert.Equal(t, NoValidationFailures(), failures)
}

func TestWrapperValidateHandlesErrors(t *testing.T) {
	validator := funcWrapperValidator{
		validatorFunc: func(checker model.ConditionalityChecker, discovery *Model) (bool, []ValidationFailure, error) {
			return false, nil, errors.New("some error")
		},
	}

	failures, err := validator.Validate(&Model{})

	assert.Nil(t, failures)
	assert.EqualError(t, err, "some error")
}

func TestValidationNoFailuresReturnsEmpty(t *testing.T) {
	failures := ValidationFailures{}

	assert.True(t, failures.Empty())
}

func TestValidationFailuresReturnsNotEmpty(t *testing.T) {
	failures := ValidationFailures{ValidationFailure{}}

	assert.False(t, failures.Empty())
}

func TestValidationNoValidationFailuresReturnsEmpty(t *testing.T) {
	assert.True(t, NoValidationFailures().Empty())
}
