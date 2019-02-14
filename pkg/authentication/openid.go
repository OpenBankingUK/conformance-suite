package authentication

// OpenIDConfiguration - The OpenID Connect discovery document retrieved by calling /.well-known/openid-configuration.
// https://openid.net/specs/openid-connect-discovery-1_0.html
type OpenIDConfiguration struct {
	TokenEndpoint string `json:"token_endpoint"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
}
