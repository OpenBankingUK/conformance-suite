package generation

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTestCaseName(t *testing.T) {
	assert.Equal(t, "#t0001", testCaseName(1))
	assert.Equal(t, "#t0000", testCaseName(0))
	assert.Equal(t, "#t10001", testCaseName(10001))
}

func TestGetGoodResponseCode(t *testing.T) {
	tcs := []struct {
		codes        []int
		expectedCode int
		err          error
	}{
		{
			[]int{200},
			200,
			nil,
		},
		{
			[]int{300, 250},
			250,
			nil,
		},
		{
			[]int{300},
			0,
			errors.New("Cannot find good response code between 200 and 299"),
		},
	}

	for _, tc := range tcs {
		t.Run("", func(t *testing.T) {
			code, err := getGoodResponseCode(tc.codes)
			assert.Equal(t, tc.expectedCode, code)
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetResponseCodes(t *testing.T) {
	op := &spec.Operation{
		OperationProps: spec.OperationProps{
			Responses: &spec.Responses{
				ResponsesProps: spec.ResponsesProps{
					StatusCodeResponses: map[int]spec.Response{
						200: {},
						300: {},
					},
				},
			},
		},
	}

	result := getResponseCodes(op)

	assert.EqualValues(t, []int{300, 200}, result)
}

func TestGetResourceIds(t *testing.T) {
	item := &discovery.ModelDiscoveryItem{ResourceIds: map[string]string{"hello": "world"}}

	result := getResourceIds(item, "/{hello}")

	assert.Equal(t, "/world", result)
}

func TestGetResourceIdsNoMatch(t *testing.T) {
	item := &discovery.ModelDiscoveryItem{ResourceIds: map[string]string{}}

	result := getResourceIds(item, "/{hello}")

	assert.Equal(t, "/{hello}", result)
}

func TestGetOperationsEmpty(t *testing.T) {
	props := &spec.PathItem{}

	results := getOperations(props)

	expected := map[string]*spec.Operation{}
	assert.Equal(t, expected, results)
}

func TestGetOperations(t *testing.T) {
	props := &spec.PathItem{
		PathItemProps: spec.PathItemProps{
			Get:  &spec.Operation{},
			Post: &spec.Operation{},
		},
	}

	results := getOperations(props)

	assert.Len(t, results, 2)
	assert.NotNil(t, results["GET"])
	assert.NotNil(t, results["POST"])
}
