package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInstructionIdentificationID(t *testing.T) {
	identifier := instructionIdentificationID()
	// Assert that identifier is alphanumeric between length 1 and 35
	assert.Regexp(t, "^[a-zA-Z0-9]{1,35}$", identifier)
}

func TestExecuteMacro(t *testing.T) {
	macroMap["helloWorld"] = func() (string, error) {
		return "hello world", nil
	}
	macroMap["noReturn"] = func() {}

	tt := []struct {
		name      string
		fnName    string
		params    []string
		expResult string
		expError  string
	}{
		{
			name:      "Hello World",
			fnName:    "helloWorld",
			params:    []string{},
			expResult: "hello world",
		},
		{
			name:      "Run function doesn't exist",
			fnName:    "missingFunction",
			params:    []string{},
			expResult: "",
			expError:  "macro not found",
		},
		{
			name:      "Call existing function with too many params",
			fnName:    "helloWorld",
			params:    []string{"p1", "p2"},
			expResult: "",
			expError:  "the number of params is not adapted",
		},
		{
			name:      "Call function that returns no values",
			fnName:    "noReturn",
			params:    []string{},
			expResult: "",
			expError:  "unable to get result from macro",
		},
	}

	for _, ti := range tt {
		t.Run(ti.name, func(t *testing.T) {
			result, err := ExecuteMacro(ti.fnName, ti.params)
			if err != nil && ti.expError != "" {
				assert.Equal(t, ti.expError, err.Error())

				return
			} else if err != nil {
				assert.Fail(t, "failed with error", err)
			}

			assert.Equal(t, ti.expResult, result)
		})
	}
}
