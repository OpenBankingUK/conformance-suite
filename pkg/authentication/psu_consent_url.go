package authentication

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"net/url"
	"path"
)

type PSUConsentClaims struct {
	Aud          string
	Iss          string
	ResponseType string
	Scope        string
	RedirectURI  string
	ConsentId    string
	State        string
}

// PSUURLGenerate generates a PSU Consent URL based on claims
func PSUURLGenerate(claims PSUConsentClaims) (*url.URL, error) {
	token, err := createAlgNoneJWT(claims)
	if err != nil {
		return nil, errors.Wrap(err, "generating psu consent URL")
	}

	consentUrl, err := url.Parse(claims.Aud)
	if err != nil {
		return nil, errors.Wrap(err, "generating psu consent URL")
	}

	consentUrl.Path = path.Join("/auth")

	consentUrlQuery := consentUrl.Query()
	consentUrlQuery.Set("client_id", claims.Iss)
	consentUrlQuery.Set("response_type", claims.ResponseType)
	consentUrlQuery.Set("scope", claims.Scope)
	consentUrlQuery.Set("request", token)
	consentUrlQuery.Set("state", claims.State)
	consentUrl.RawQuery = consentUrlQuery.Encode()

	return consentUrl, nil
}

func createAlgNoneJWT(claims PSUConsentClaims) (string, error) {
	alg := jwt.SigningMethodNone
	if alg == nil {
		return "", fmt.Errorf("no signing method: %v", alg)
	}

	token := &jwt.Token{
		Header: map[string]interface{}{
			"alg": alg.Alg(),
		},
		Claims: makeOpenBankingJWTClaims(claims),
		Method: alg,
	}

	tokenString, err := token.SigningString()
	if err != nil {
		return "", err
	}

	// alg none might cause problems for some JWT parses so we add the "."
	// to terminate the jwt
	tokenString = tokenString + "."

	return tokenString, nil
}

type openBankingClaims struct {
	IdToken idToken `json:"id_token,omitempty"`
}

type idToken struct {
	IntentID intentId `json:"openbanking_intent_id,omitempty"`
}

type intentId struct {
	Essential bool   `json:"essential"`
	Value     string `json:"value"`
}

func makeOpenBankingJWTClaims(claims PSUConsentClaims) jwt.MapClaims {
	return jwt.MapClaims{
		"iss":          claims.Iss,
		"scope":        claims.Scope,
		"aud":          claims.Aud,
		"redirect_uri": claims.RedirectURI,
		"claims": openBankingClaims{
			IdToken: idToken{
				IntentID: intentId{
					Essential: true,
					Value:     claims.ConsentId,
				},
			},
		},
	}
}
