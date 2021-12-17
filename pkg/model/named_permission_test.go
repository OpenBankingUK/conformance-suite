package model

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/names"
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
	nameGenerator := names.NewSequentialPrefixedName("to")

	specTokens := NewSpecConsentRequirements(nameGenerator, codeSetResult, "id")

	assert.Equal(t, "id", specTokens.Identifier)
	assert.Len(t, specTokens.NamedPermissions, 2)
}
