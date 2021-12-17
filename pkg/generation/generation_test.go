package generation

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"fmt"
	"sort"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
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
	genConfig := GeneratorConfig{
		ResourceIDs: model.ResourceIDs{
			AccountIDs: []model.ResourceAccountID{
				{AccountID: "12345"},
			},
			StatementIDs: []model.ResourceStatementID{
				{StatementID: "6789"},
			},
		},
	}

	result := getResourceIds(item, "/{AccountId}", genConfig)
	assert.Equal("/12345", result)
	result = getResourceIds(item, "/{StatementId}", genConfig)
	assert.Equal("/6789", result)
	result = getResourceIds(item, "/{AccountId}/{StatementId}", genConfig)
	assert.Equal("/12345/6789", result)
}

func TestGetResourceIdsNoMatch(t *testing.T) {
	assert := test.NewAssert(t)

	item := &discovery.ModelDiscoveryItem{ResourceIds: map[string]string{}}

	result := getResourceIds(item, "/{hello}", GeneratorConfig{})

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
