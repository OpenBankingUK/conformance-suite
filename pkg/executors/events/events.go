package events

// Events -
type Events interface {
	Tokens() chan AcquiredAccessToken
	AllTokens() chan AcquiredAllAccessTokens
}

// NewEvents -
func NewEvents() Events {
	const size = 100
	return events{
		acquiredAccessTokensChan:   make(chan AcquiredAccessToken, size),
		aquiredAllAccessTokensChan: make(chan AcquiredAllAccessTokens, size),
	}
}

type events struct {
	acquiredAccessTokensChan   chan AcquiredAccessToken
	aquiredAllAccessTokensChan chan AcquiredAllAccessTokens
}

func (e events) Tokens() chan AcquiredAccessToken {
	return e.acquiredAccessTokensChan
}

func (e events) AllTokens() chan AcquiredAllAccessTokens {
	return e.aquiredAllAccessTokensChan
}
