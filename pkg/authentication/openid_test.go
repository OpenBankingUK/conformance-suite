package authentication

import (
	"encoding/json"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
)

func TestOpenIDUnmarshal(t *testing.T) {
	require := test.NewRequire(t)

	data := `
{
	"token_endpoint": "https://modelobank2018.o3bank.co.uk:4201/<token_mock>",
	"authorization_endpoint": "https://modelobankauth2018.o3bank.co.uk:4101/<auth_mock>"
}
	`
	expected := OpenIDConfiguration{
		TokenEndpoint:         "https://modelobank2018.o3bank.co.uk:4201/<token_mock>",
		AuthorizationEndpoint: "https://modelobankauth2018.o3bank.co.uk:4101/<auth_mock>",
	}
	actual := OpenIDConfiguration{}
	require.NoError(json.Unmarshal([]byte(data), &actual))

	require.Equal(expected, actual)
}
