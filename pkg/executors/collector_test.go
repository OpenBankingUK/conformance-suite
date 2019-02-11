package executors

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/permissions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewCollector(t *testing.T) {
	specRequirements := []model.SpecConsentRequirements{
		{
			Identifier: "spec1",
			NamedPermissions: model.NamedPermissions{
				{
					Name: "namedSet1",
					CodeSet: permissions.CodeSetResult{
						CodeSet: permissions.CodeSet{"Read", "Write"},
						TestIds: []permissions.TestId{"test1"},
					},
				},
			},
		},
	}

	called := false
	collector := NewCollector(specRequirements, func() { called = true })

	assert.False(t, called)
	assert.Equal(t, 1, collector.countNamedSets())
	assert.True(t, collector.setNameExists("namedSet1"))
}

func TestNewCollectorNamedSetDoesntExist(t *testing.T) {
	specRequirements := []model.SpecConsentRequirements{}

	collector := NewCollector(specRequirements, func() {})

	assert.Equal(t, 0, collector.countNamedSets())
	assert.False(t, collector.setNameExists("namedSet1"))
}

func TestNewCollectorCollectsAndDone(t *testing.T) {
	specRequirements := []model.SpecConsentRequirements{
		{
			Identifier: "spec1",
			NamedPermissions: model.NamedPermissions{
				{
					Name: "namedSet1",
					CodeSet: permissions.CodeSetResult{
						CodeSet: permissions.CodeSet{"Read", "Write"},
						TestIds: []permissions.TestId{"test1"},
					},
				},
				{
					Name: "namedSet2",
					CodeSet: permissions.CodeSetResult{
						CodeSet: permissions.CodeSet{"Read", "Write"},
						TestIds: []permissions.TestId{"test2"},
					},
				},
			},
		},
	}
	called := false
	collector := NewCollector(specRequirements, func() { called = true })

	err := collector.Collect("namedSet1", "token from named set 1")
	require.NoError(t, err)
	err = collector.Collect("namedSet2", "token from named set 2")
	require.NoError(t, err)

	tokens := collector.Tokens()

	assert.True(t, called)
	assert.Len(t, tokens, 2)
	assert.Equal(t, "token from named set 1", tokens[0].Code)
	assert.Equal(t, "token from named set 2", tokens[1].Code)
}

func TestNewCollectorCollectsReturnsErrorForNonExistingNamedSet(t *testing.T) {
	specRequirements := []model.SpecConsentRequirements{}
	collector := NewCollector(specRequirements, func() {})

	err := collector.Collect("namedSet1", "token from named set 1")

	assert.EqualError(t, err, "invalid permission set name")
}
