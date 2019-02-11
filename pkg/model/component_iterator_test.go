package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateComponentIterator(t *testing.T) {
	var param map[string]interface{}
	ci := NewComponentIterator(nil, param)
	assert.NotNil(t, ci)
}

func TestIterateOverComponentIterator(t *testing.T) {

	sites := [][]string{{"a", "b", "c"}, {"d", "e", "f"}, {"h", "i", "j"}, {"d"}}
	tokenNames := []string{"TokA", "TokB", "TokD", "TokE", "TokF"}
	dummy := []string{"x", "y", "z", "a", "b", "c", "d", "e"}

	params := make(map[string]interface{})
	params["permissions"] = sites
	params["token_names"] = tokenNames
	params["dummy"] = dummy

	comp, err := LoadComponent("../../templates/tokenProviderComponent.json")
	assert.Nil(t, err)

	ci := NewComponentIterator(&comp, params)
	ci.Iterate(*emptyContext)

}
