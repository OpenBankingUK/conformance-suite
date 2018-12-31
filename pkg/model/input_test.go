package model

import (
	"os"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/tracer"
	"github.com/stretchr/testify/assert"
)

var tc = &TestCase{}
var ctx = &Context{}

func TestMain(m *testing.M) {
	tracer.Silent = false
	m.Run()
	//time.Sleep(1 * time.Second)
	os.Exit(0)
}

func TestCreateRequestEmptyEndpointOrMethod(t *testing.T) {
	i := &Input{}
	req, err := i.CreateRequest(tc, ctx)
	assert.NotNil(t, err)
	assert.Nil(t, req)
	t.Log(err.Error())

	i = &Input{Endpoint: "http://google.com"}
	req, err = i.CreateRequest(tc, ctx)
	assert.NotNil(t, err)
	assert.Nil(t, req)
	t.Log(err.Error())

	i = &Input{Method: "GET"}
	req, err = i.CreateRequest(tc, ctx)
	assert.NotNil(t, err)
	assert.Nil(t, req)
	t.Log(err.Error())
}

func TestInputGetValuesMissingContextVariable(t *testing.T) {
	match := Match{Description: "simple match test", ContextName: "GetValueToFind"}
	accessor := ContextAccessor{Matches: []Match{match}}
	i := &Input{Method: "GET", Endpoint: "http://google.com", ContextGet: accessor}
	req, err := i.CreateRequest(tc, ctx)
	assert.NotNil(t, err)
	assert.Nil(t, req)
	t.Log(err.Error())

}
