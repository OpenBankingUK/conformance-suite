package model

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/permissions"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNamedPermission(t *testing.T) {
	codeSetResult := permissions.CodeSetResult{
		CodeSet: permissions.CodeSet{"a"},
	}

	token := newNamedPermission("first", codeSetResult)

	assert.Equal(t, "first", token.Name)
	assert.Equal(t, codeSetResult, token.CodeSet)
}

func TestNewNamedPermissionsAdd(t *testing.T) {
	tokens := NamedPermissions{}
	token := newNamedPermission("first", permissions.CodeSetResult{})

	tokens.Add(token)

	assert.Len(t, tokens, 1)
	assert.Equal(t, token, tokens[0])
}

func TestNewSpecConsentRequirements(t *testing.T) {
	codeSetResult := permissions.CodeSetResultSet{
		{
			CodeSet: permissions.CodeSet{"a"},
		},
		{
			CodeSet: permissions.CodeSet{"b"},
		},
	}

	specTokens := NewSpecConsentRequirements(codeSetResult, "id")

	assert.Equal(t, "id", specTokens.Identifier)
	assert.Len(t, specTokens.NamedPermissions, 2)
}

func TestPrefixedNumber(t *testing.T) {
	assert.Equal(t, "to0", prefixedNumber(0, "to"))
	assert.Equal(t, "to1", prefixedNumber(1, "to"))
	assert.Equal(t, "to99", prefixedNumber(99, "to"))
	assert.Equal(t, "99", prefixedNumber(99, ""))
}
func TestRandomString(t *testing.T) {
	randomStringResult := randomString(5)
	assert.Len(t, randomStringResult, 5)

	randomStringResult2 := randomString(5)
	assert.NotEqual(t, randomStringResult, randomStringResult2)
}
