package names

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSequentialPrefixedNameGenerate(t *testing.T) {
	generator := NewSententialPrefixedName("#t")

	assert.Equal(t, "#t1001", generator.Generate())
	assert.Equal(t, "#t1002", generator.Generate())

	generator = NewSententialPrefixedName("")

	assert.Equal(t, "1001", generator.Generate())
	assert.Equal(t, "1002", generator.Generate())
}
