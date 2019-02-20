package authentication

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestClientSecretBasicGeneration tests the generation of `client secret basic` value as a product of
// `client_id` and `client_secret` as per https://tools.ietf.org/html/rfc7617
func TestClientSecretBasicGeneration(t *testing.T) {
	assert := assert.New(t)

	tt := []struct {
		clientID      string
		clientSecret  string
		tokenExpected string
		errorExpected error
		label         string
	}{
		{
			clientID:      "dc3a363e-2cc3-4187-b6df-579f21bad6c8",
			clientSecret:  "e648104b-f52a-43e1-a2e0-fe3a047497cf",
			tokenExpected: "ZGMzYTM2M2UtMmNjMy00MTg3LWI2ZGYtNTc5ZjIxYmFkNmM4OmU2NDgxMDRiLWY1MmEtNDNlMS1hMmUwLWZlM2EwNDc0OTdjZg==",
			label:         "valid credentials",
		},
		{
			clientID:      "",
			clientSecret:  "foobar",
			tokenExpected: "",
			errorExpected: errors.New("clientID cannot be empty"),
			label:         "empty client id",
		},
		{
			clientID:      "foobar",
			clientSecret:  "",
			tokenExpected: "",
			errorExpected: errors.New("clientSecret cannot be empty"),
			label:         "empty client secret",
		},
	}

	for _, ti := range tt {
		tokenActual, err := CalculateClientSecretBasicToken(ti.clientID, ti.clientSecret)
		if ti.errorExpected != nil {
			assert.Equal(ti.errorExpected.Error(), err.Error())
		}

		assert.Equal(ti.tokenExpected, tokenActual, ti.label)
	}
}
