package authentication

import (
	"encoding/json"
	"net/http"
	"fmt"
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
		return config, fmt.Errorf(fmt.Sprintf("Failed to GET OpenID config: %s - %s", url, err))
	} else {
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return config, fmt.Errorf("Failed to GET OpenID config: %s - HTTP response status: %d", url, resp.StatusCode)
		}
		if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
			return config, fmt.Errorf(fmt.Sprintf("Invalid OpenID config JSON returned: %s - %s", url, err))
		}
		return config, nil
	}
}
