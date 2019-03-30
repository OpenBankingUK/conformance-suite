package events

// AcquiredAccessToken - When `code` has been exchanged for an `access_token`.
type AcquiredAccessToken struct {
	TokenName string `json:"token_name"`
}

// AcquiredAllAccessTokens - When all `code`s have been exchanged for `access_token`s.
type AcquiredAllAccessTokens struct {
	TokenNames []string `json:"token_names"`
}

// NewAcquiredAccessToken -
func NewAcquiredAccessToken(tokenName string) AcquiredAccessToken {
	return AcquiredAccessToken{
		TokenName: tokenName,
	}
}

// NewAcquiredAccessToken -
func NewAcquiredAllAccessTokens(tokenNames []string) AcquiredAllAccessTokens {
	return AcquiredAllAccessTokens{
		TokenNames: tokenNames,
	}
}
