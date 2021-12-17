package model

import (
	"github.com/OpenBankingUK/conformance-suite/pkg/names"
	"github.com/OpenBankingUK/conformance-suite/pkg/permissions"
)

// NamedPermission - permission structure
type NamedPermission struct {
	Name       string                    `json:"name"`
	CodeSet    permissions.CodeSetResult `json:"codeSet"`
	ConsentURL string                    `json:"consentUrl"`
}

// NamedPermissions - permission structure
type NamedPermissions []NamedPermission

// Add - to named permissions
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

// SpecConsentRequirements -
type SpecConsentRequirements struct {
	Identifier       string           `json:"specIdentifier"`
	NamedPermissions NamedPermissions `json:"namedPermissions"`
}

// NewSpecConsentRequirements - create a new SpecConsentRequirements
func NewSpecConsentRequirements(nameGenerator names.Generator, result permissions.CodeSetResultSet, specID string) SpecConsentRequirements {
	namedPermissions := NamedPermissions{}
	for _, resultSet := range result {
		namedPermission := newNamedPermission(nameGenerator.Generate(), resultSet)
		namedPermissions = append(namedPermissions, namedPermission)
	}
	return SpecConsentRequirements{
		Identifier:       specID,
		NamedPermissions: namedPermissions,
	}
}
