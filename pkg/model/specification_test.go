package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpecificationIdentifierFromSchemaVersion(t *testing.T) {
	t.Run("returns specification identifier when valid schema version URL", func(t *testing.T) {
		specification, err := SpecificationIdentifierFromSchemaVersion("https://raw.githubusercontent.com/OpenBankingUK/read-write-api-specs/v3.0.0/dist/account-info-swagger.json")
		assert.Nil(t, err)
		assert.Equal(t, specification, "account-transaction-v3.0")
	})

	t.Run("returns error when invalid schema version URL", func(t *testing.T) {
		specification, err := SpecificationIdentifierFromSchemaVersion("https://example.com/invalid")
		assert.NotNil(t, err)
		assert.Equal(t, specification, "")
	})
}
