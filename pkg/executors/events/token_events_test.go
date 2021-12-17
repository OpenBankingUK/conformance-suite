package events

import (
	"encoding/json"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"

	"testing"
)

func TestAcquiredAccessTokenResultJsonMarshal(t *testing.T) {
	require := test.NewRequire(t)

	expected := `
{
    "token_name": "to1001"
}
	`
	tokenName := "to1001"
	result := NewAcquiredAccessToken(tokenName)

	actual, err := json.Marshal(result)
	require.NoError(err)
	require.NotEmpty(actual)

	require.JSONEq(expected, string(actual))
}

func TestAcquiredAllAccessTokensResultJsonMarshal(t *testing.T) {
	require := test.NewRequire(t)

	expected := `
{
    "token_names": ["to1001"]
}
	`
	tokenNames := []string{"to1001"}
	result := NewAcquiredAllAccessTokens(tokenNames)

	actual, err := json.Marshal(result)
	require.NoError(err)
	require.NotEmpty(actual)

	require.JSONEq(expected, string(actual))
}
