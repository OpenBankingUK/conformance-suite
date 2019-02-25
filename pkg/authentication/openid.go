package authentication

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

// OpenIDConfiguration - The OpenID Connect discovery document retrieved by calling /.well-known/openid-configuration.
// https://openid.net/specs/openid-connect-discovery-1_0.html
type OpenIDConfiguration struct {
	TokenEndpoint         string `json:"token_endpoint"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	Issuer                string `json:"issuer"`
}

func OpenIdConfig(url string) (OpenIDConfiguration, error) {
	config := OpenIDConfiguration{}
	resp, err := http.Get(url)
	if err != nil {
		return config, errors.Wrap(err, fmt.Sprintf("Failed to GET OpenID config: %s", url))
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return config, fmt.Errorf("failed to GET OpenID config: %s : HTTP response status: %d", url, resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return config, errors.Wrap(err, fmt.Sprintf("Invalid OpenID config JSON returned: %s ", url))
	}
	return config, nil
}
