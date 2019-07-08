package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContextGetPutString(t *testing.T) {
	context := Context{}
	context.PutString("string", "a string")

	get, err := context.GetString("string")

	assert.NoError(t, err)
	assert.Equal(t, "a string", get)
}

func TestContextGetPutStringNotFound(t *testing.T) {
	context := Context{}

	get, err := context.GetString("non existing key")

	assert.Error(t, err)
	assert.IsType(t, ErrNotFound, err)
	assert.Equal(t, "", get)
}

func TestContextGetStringHandlesCast(t *testing.T) {
	context := Context{}
	context.Put("wrong", 1)

	get, err := context.GetString("wrong")

	assert.EqualError(t, err, "error casting key to string")
	assert.Equal(t, "", get)
}

func TestContextGetInt(t *testing.T) {
	context := Context{}
	context.Put("int", 99)

	get, ok := context.Get("int")

	assert.True(t, ok)
	assert.IsType(t, int(0), get)
	assert.Equal(t, 99, get)
}

func TestContextGetStringSlice(t *testing.T) {
	context := Context{}
	var values []interface{}
	values = append(values, "read", "write")
	context.Put("array", values)

	get, err := context.GetStringSlice("array")

	assert.NoError(t, err)
	assert.IsType(t, []string{}, get)
}

func TestContextGetStringSliceHandlesCastError(t *testing.T) {
	context := Context{}
	context.Put("array", []string{"a"})

	get, err := context.GetStringSlice("array")

	assert.EqualError(t, err, "cast error can't get string slice from context")
	assert.Nil(t, get)
}

func TestContextGetStringSliceHandlesCastErrorFromElement(t *testing.T) {
	context := Context{}
	var values []interface{}
	values = append(values, 0)
	context.Put("array", values)

	get, err := context.GetStringSlice("array")

	assert.EqualError(t, err, "element cast error can't get string slice from context")
	assert.Nil(t, get)
}

func TestContextPutStringSlice(t *testing.T) {
	context := Context{}

	context.PutStringSlice("array", []string{"a", "b"})
	get, err := context.GetStringSlice("array")

	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, get)
}

func TestIsSet(t *testing.T) {
	context := Context{}
	context.Put("a", "foobar")
	context.Put("b", struct{}{})
	context.Put("c", "")
	context.Put("d", nil)

	assert.True(t, context.IsSet("a"))
	assert.True(t, context.IsSet("b"))
	assert.False(t, context.IsSet("c"))
	assert.False(t, context.IsSet("d"))
}
