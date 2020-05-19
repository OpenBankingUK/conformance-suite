package authentication

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/client"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// OpenIDConfiguration - The OpenID Connect discovery document retrieved by calling /.well-known/openid-configuration.
// https://openid.net/specs/openid-connect-discovery-1_0.html
type OpenIDConfiguration struct {
	TokenEndpoint                          string   `json:"token_endpoint,omitempty"`
	TokenEndpointAuthMethodsSupported      []string `json:"token_endpoint_auth_methods_supported,omitempty"`
	RequestObjectSigningAlgValuesSupported []string `json:"request_object_signing_alg_values_supported,omitempty"`
	AuthorizationEndpoint                  string   `json:"authorization_endpoint,omitempty"`
	Issuer                                 string   `json:"issuer,omitempty"`
	ResponseTypesSupported                 []string `json:"response_types_supported,omitempty"`
	AcrValuesSupported                     []string `json:"acr_values_supported,omitempty"`
	JwksURI                                string   `json:"jwks_uri,omitempty"`
}

var jwks_uri_accessor = ""

func GetJWKSUri() string {
	return jwks_uri_accessor
}

func OpenIdConfig(url string) (OpenIDConfiguration, error) {
	resp, err := client.NewHTTPClient(client.DefaultTimeout).Get(url)
	if err != nil {
		return OpenIDConfiguration{}, errors.Wrapf(err, "Failed to GET OpenIDConfiguration: url=%+v", url)
	}

	if resp.StatusCode != http.StatusOK {
		responseBody, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return OpenIDConfiguration{}, errors.Wrap(err, "error reading error response from GET OpenIDConfiguration")
		}

		return OpenIDConfiguration{}, fmt.Errorf(
			"failed to GET OpenIDConfiguration config: url=%+v, StatusCode=%+v, body=%+v",
			url,
			resp.StatusCode,
			string(responseBody),
		)
	}

	defer resp.Body.Close()
	config := OpenIDConfiguration{}
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return config, errors.Wrap(err, fmt.Sprintf("Invalid OpenIDConfiguration: url=%+v", url))
	}

	logrus.Tracef("JWKS Uri = %s", config.JwksURI)
	jwks_uri_accessor = config.JwksURI
	return config, nil
}
