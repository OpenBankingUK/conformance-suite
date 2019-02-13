package generation

import (
	"fmt"
	"sort"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
)

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

	for index, tc := range tcs {
		tc := tc
		t.Run(fmt.Sprintf("TestGetGoodResponseCode/%d", index), func(t *testing.T) {
			assert := test.NewAssert(t)

			code, err := getGoodResponseCode(tc.codes)
			assert.Equal(tc.expectedCode, code)
			if tc.err != nil {
				assert.EqualError(err, tc.err.Error())
			} else {
				assert.NoError(err)
			}
		})
	}
}

func TestGetResponseCodes(t *testing.T) {
	assert := test.NewAssert(t)

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
	sort.Ints(result)

	expected := []int{200, 300}
	assert.Equal(expected, result)
}

func TestGetResourceIds(t *testing.T) {
	assert := test.NewAssert(t)

	item := &discovery.ModelDiscoveryItem{ResourceIds: map[string]string{"hello": "world"}}

	result := getResourceIds(item, "/{hello}")

	assert.Equal("/world", result)
}

func TestGetResourceIdsNoMatch(t *testing.T) {
	assert := test.NewAssert(t)

	item := &discovery.ModelDiscoveryItem{ResourceIds: map[string]string{}}

	result := getResourceIds(item, "/{hello}")

	assert.Equal("/{hello}", result)
}

func TestGetOperationsEmpty(t *testing.T) {
	assert := test.NewAssert(t)

	props := &spec.PathItem{}

	results := getOperations(props)

	expected := map[string]*spec.Operation{}
	assert.Equal(expected, results)
}

func TestGetOperations(t *testing.T) {
	assert := test.NewAssert(t)

	props := &spec.PathItem{
		PathItemProps: spec.PathItemProps{
			Get:  &spec.Operation{},
			Post: &spec.Operation{},
		},
	}

	results := getOperations(props)

	assert.Len(results, 2)
	assert.NotNil(results["GET"])
	assert.NotNil(results["POST"])
}
