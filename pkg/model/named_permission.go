package model

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/names"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/permissions"
)

type NamedPermission struct {
	Name       string                    `json:"name"`
	CodeSet    permissions.CodeSetResult `json:"codeSet"`
	ConsentUrl string                    `json:"consentUrl"`
}

type NamedPermissions []NamedPermission

func (t *NamedPermissions) Add(token NamedPermission) {
	*t = append(*t, token)
}

// newNamedPermission create a token required to run test cases
// generates a unique name
func newNamedPermission(name string, codeSet permissions.CodeSetResult) NamedPermission {
	return NamedPermission{
		Name:    name,
		CodeSet: codeSet,
	}
}

type SpecConsentRequirements struct {
	Identifier       string           `json:"specIdentifier"`
	NamedPermissions NamedPermissions `json:"namedPermissions"`
}

func NewSpecConsentRequirements(nameGenerator names.Generator, result permissions.CodeSetResultSet, specId string) SpecConsentRequirements {
	namedPermissions := NamedPermissions{}
	for _, resultSet := range result {
		namedPermission := newNamedPermission(nameGenerator.Generate(), resultSet)
		namedPermissions = append(namedPermissions, namedPermission)
	}
	return SpecConsentRequirements{
		Identifier:       specId,
		NamedPermissions: namedPermissions,
	}
}
