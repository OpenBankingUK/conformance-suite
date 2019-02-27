package authentication

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/client"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

// OpenIDConfiguration - The OpenID Connect discovery document retrieved by calling /.well-known/openid-configuration.
// https://openid.net/specs/openid-connect-discovery-1_0.html
type OpenIDConfiguration struct {
	TokenEndpoint                     string   `json:"token_endpoint"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
	AuthorizationEndpoint             string   `json:"authorization_endpoint"`
	Issuer                            string   `json:"issuer"`
}

var AUTH_METHODS_SORTED_MOST_SECURE_FIRST = []string{
	"tls_client_auth", // most secure
	"private_key_jwt",
	"client_secret_jwt",
	"client_secret_post",
	"client_secret_basic", // least secure
}

func OpenIdConfig(url string) (OpenIDConfiguration, error) {
	body, e := retrieveConfig(url)
	if body != nil {
		defer body.Close()
	}
	if e != nil {
		return OpenIDConfiguration{}, e
	}

	config := OpenIDConfiguration{}
	if err := json.NewDecoder(body).Decode(&config); err != nil {
		return config, errors.Wrap(err, fmt.Sprintf("Invalid OpenID config JSON returned: %s ", url))
	}
	config.TokenEndpointAuthMethodsSupported = sortAuthMethodsMostSecureFirst(config.TokenEndpointAuthMethodsSupported)
	return config, nil
}

func retrieveConfig(url string) (io.ReadCloser, error) {
	resp, err := client.NewHTTPClient(client.DefaultTimeout).Get(url)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to GET OpenID config: %s", url))
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to GET OpenID config: %s : HTTP response status: %d", url, resp.StatusCode)
	}
	return resp.Body, nil
}

func sortAuthMethodsMostSecureFirst(methods []string) []string {
	sorted := make([]string, len(methods))
	i := 0
	for _, a := range AUTH_METHODS_SORTED_MOST_SECURE_FIRST {
		for index, m := range methods {
			if a == m {
				sorted[i] = m
				methods[index] = ""
				i = i + 1
			}
		}
	}
	for _, m := range methods {
		if m != "" {
			logger := logrus.StandardLogger().WithField("authentication", "OpenIDConfiguration")
			logger.Infof("Invalid token endpoint auth method in OpenID config: %s", m)
		}
	}

	return sorted
}
