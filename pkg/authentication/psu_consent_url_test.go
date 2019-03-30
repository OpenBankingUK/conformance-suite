package authentication

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPSUURLGenerate(t *testing.T) {
	claims := PSUConsentClaims{
		Aud:          "https://server",
		Iss:          "iss",
		Scope:        "scope",
		ResponseType: "responseType",
		RedirectURI:  "redirect_uri",
		ConsentId:    "123",
		State:        "state",
	}

	url, err := PSUURLGenerate(claims)

	require.NoError(t, err)
	assert.Equal(t, claims.Iss, url.Query().Get("client_id"))
	assert.Equal(t, claims.ResponseType, url.Query().Get("response_type"))
	assert.Equal(t, claims.Scope, url.Query().Get("scope"))
	assert.Equal(t, claims.State, url.Query().Get("state"))
	token, err := createAlgNoneJWT(claims)
	require.NoError(t, err)
	assert.Equal(t, token, url.Query().Get("request"))
}

func TestCreateAlgNoneJWTEmpty(t *testing.T) {
	claims := PSUConsentClaims{}

	jwtString, err := createAlgNoneJWT(claims)

	require.NoError(t, err)
	expected := "eyJhbGciOiJub25lIn0.eyJhdWQiOiIiLCJjbGFpbXMiOnsiaWRfdG9rZW4iOnsib3BlbmJhbmtpbmdfaW50ZW50X2lkIjp7ImVzc2VudGlhbCI6dHJ1ZSwidmFsdWUiOiIifX19LCJpc3MiOiIiLCJyZWRpcmVjdF91cmkiOiIiLCJzY29wZSI6IiJ9."
	assert.Equal(t, expected, jwtString)
}

func TestCreateAlgNoneJWTUsesNoneAlg(t *testing.T) {
	claims := PSUConsentClaims{}
	jwtString, err := createAlgNoneJWT(claims)
	require.NoError(t, err)

	jwt, err := jwt.Parse(jwtString, nil)

	assert.Error(t, err) // error expected as we using no keyFunc to parse
	assert.Equal(t, "none", jwt.Header["alg"])
}

func TestCreateAlgNoneJWTUsesClaims(t *testing.T) {
	claims := PSUConsentClaims{
		Aud:         "http://server",
		Iss:         "iss",
		Scope:       "scope",
		RedirectURI: "redirect_uri",
		ConsentId:   "123",
	}
	jwtString, err := createAlgNoneJWT(claims)
	require.NoError(t, err)

	token, err := jwt.Parse(jwtString, nil)

	assert.Error(t, err) // error expected as we using no keyFunc to parse
	c, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)
	assert.Equal(t, "http://server", c["aud"])
	assert.Equal(t, "iss", c["iss"])
	assert.Equal(t, "scope", c["scope"])
	assert.Equal(t, "redirect_uri", c["redirect_uri"])

	// extraction consent id
	cc, ok := c["claims"].(map[string]interface{})
	require.True(t, ok)
	idt, ok := cc["id_token"].(map[string]interface{})
	require.True(t, ok)
	iid, ok := idt["openbanking_intent_id"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "123", iid["value"])
}
