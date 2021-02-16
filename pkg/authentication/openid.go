package authentication

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/client"
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

type CachedOpenIdConfigGetter struct {
	client *http.Client
	cache  map[string]OpenIDConfiguration
}

func NewOpenIdConfigGetter() *CachedOpenIdConfigGetter {
	return &CachedOpenIdConfigGetter{
		client: client.NewHTTPClient(client.DefaultTimeout),
		cache:  map[string]OpenIDConfiguration{},
	}
}

func (g CachedOpenIdConfigGetter) Get(url string) (OpenIDConfiguration, error) {
	config, ok := g.cache[url]
	if ok {
		logrus.Tracef("Cache hit on getting openid config Uri = %s", url)
		return config, nil
	}

	resp, err := g.client.Get(url)
	if err != nil {
		return OpenIDConfiguration{}, fmt.Errorf("Failed to GET OpenIDConfiguration: url=%+v : %w", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		responseBody, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return OpenIDConfiguration{}, fmt.Errorf("error reading error response from GET OpenIDConfiguration: %w", err)
		}

		return OpenIDConfiguration{}, fmt.Errorf(
			"failed to GET OpenIDConfiguration config: url=%+v, StatusCode=%+v, body=%+v",
			url,
			resp.StatusCode,
			string(responseBody),
		)
	}

	defer resp.Body.Close()
	config = OpenIDConfiguration{}
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return config, fmt.Errorf("Invalid OpenIDConfiguration: url=%+v: %w", url, err)
	}

	logrus.Tracef("JWKS Uri = %s", config.JwksURI)
	jwks_uri_accessor = config.JwksURI
	g.cache[url] = config
	return config, nil
}
