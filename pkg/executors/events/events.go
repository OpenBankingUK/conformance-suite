package events

// Events -
type Events interface {
	AddAcquiredAccessToken(acquiredAccessToken AcquiredAccessToken)
	TokensChannel() <-chan AcquiredAccessToken
	AllAcquiredAccessToken() []AcquiredAccessToken

	AddAcquiredAllAccessTokens(acquiredAllAccessTokens AcquiredAllAccessTokens)
	AllTokensChannel() <-chan AcquiredAllAccessTokens
	AllAcquiredAllAccessTokens() []AcquiredAllAccessTokens
}

// NewEvents -
func NewEvents() Events {
	const size = 100
	return &events{
		acquiredAccessTokens:       []AcquiredAccessToken{},
		acquiredAccessTokensChan:   make(chan AcquiredAccessToken, size),
		acquiredAllAccessTokens:    []AcquiredAllAccessTokens{},
		aquiredAllAccessTokensChan: make(chan AcquiredAllAccessTokens, size),
	}
}

type events struct {
	acquiredAccessTokens       []AcquiredAccessToken
	acquiredAccessTokensChan   chan AcquiredAccessToken
	acquiredAllAccessTokens    []AcquiredAllAccessTokens
	aquiredAllAccessTokensChan chan AcquiredAllAccessTokens
}

func (e *events) AddAcquiredAccessToken(acquiredAccessToken AcquiredAccessToken) {
	e.acquiredAccessTokens = append(e.acquiredAccessTokens, acquiredAccessToken)
	e.acquiredAccessTokensChan <- acquiredAccessToken
}

func (e *events) TokensChannel() <-chan AcquiredAccessToken {
	return e.acquiredAccessTokensChan
}

func (e *events) AllAcquiredAccessToken() []AcquiredAccessToken {
	return e.acquiredAccessTokens
}

func (e *events) AddAcquiredAllAccessTokens(acquiredAllAccessTokens AcquiredAllAccessTokens) {
	e.acquiredAllAccessTokens = append(e.acquiredAllAccessTokens, acquiredAllAccessTokens)
	e.aquiredAllAccessTokensChan <- acquiredAllAccessTokens
}

func (e *events) AllTokensChannel() <-chan AcquiredAllAccessTokens {
	return e.aquiredAllAccessTokensChan
}

func (e *events) AllAcquiredAllAccessTokens() []AcquiredAllAccessTokens {
	return e.acquiredAllAccessTokens
}
