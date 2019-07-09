package authentication

import (
	"fmt"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
)

func TestCalculateCHash(t *testing.T) {
	require := test.NewRequire(t)

	tt := []struct {
		label         string
		code          string
		alg           string
		expectedHash  string
		expectedError error
	}{
		{
			label:        "ES256 empty code",
			code:         "",
			alg:          "ES256",
			expectedHash: "47DEQpj8HBSa-_TImW-5JA",
		},
		{
			label:        "ES256 code valid",
			code:         "80bf17a3-e617-4983-9d62-b50bd8e6fce4",
			alg:          "ES256",
			expectedHash: "EE_Bf-grXWv5GGhs5FZ0ug",
		},
		{
			label:        "PS256 code valid",
			code:         "80bf17a3-e617-4983-9d62-b50bd8e6fce4",
			alg:          "PS256",
			expectedHash: "EE_Bf-grXWv5GGhs5FZ0ug",
		},
		{
			label:         "algorithm not supported",
			code:          "80bf17a3-e617-4983-9d62-b50bd8e6fce4",
			alg:           "bad-algorithm",
			expectedHash:  "",
			expectedError: fmt.Errorf(`authentication.CalculateCHash: "bad-algorithm" algorithm not supported`),
		},
	}

	for _, tti := range tt {
		cHash, err := CalculateCHash(tti.alg, tti.code)
		require.Equal(tti.expectedHash, cHash, tti.label)

		if tti.expectedError != nil {
			require.Equal(tti.expectedError, err, tti.label)
		}

	}
}
