package model

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/permissions"
	"math/rand"
	"strconv"
	"time"
)

type NamedPermission struct {
	Name    string
	CodeSet permissions.CodeSetResult
}

type NamedPermissions []NamedPermission

func (t *NamedPermissions) Add(token NamedPermission) {
	*t = append(*t, token)
}

// newNamedPermission create a token required to run test cases
// generates a unique name
func newNamedPermission(name string, codeSet permissions.CodeSetResult) NamedPermission {
	return NamedPermission{
		name,
		codeSet,
	}
}

type SpecConsentRequirements struct {
	Identifier       string           `json:"specIdentifier"`
	NamedPermissions NamedPermissions `json:"namedPermissions"`
}

func NewSpecConsentRequirements(result permissions.CodeSetResultSet, specId string) SpecConsentRequirements {
	tokens := NamedPermissions{}
	tokenNamePrefix := "to-" + randomString(3) + "-"
	for i, resultSet := range result {
		token := newNamedPermission(prefixedNumber(i, tokenNamePrefix), resultSet)
		tokens = append(tokens, token)
	}
	return SpecConsentRequirements{
		Identifier:       specId,
		NamedPermissions: tokens,
	}
}

// prefixedNumber returns a number prefixed by a string
func prefixedNumber(next int, prefix string) string {
	return prefix + strconv.Itoa(next)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func randomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
