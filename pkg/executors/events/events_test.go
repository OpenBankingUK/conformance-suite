package events

import (
	"testing"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
)

const (
	selectTimeout = 1 * time.Millisecond
)

func TestEvents(t *testing.T) {
	require := test.NewRequire(t)

	events := NewEvents()
	// initially empty
	select {
	case msg, ok := <-events.TokensChannel():
		require.Nil(msg)
		require.False(ok)
	case msg, ok := <-events.AllTokensChannel():
		require.Nil(msg)
		require.False(ok)
	case <-time.After(selectTimeout):
		break
	}

	require.Empty(events.AllAcquiredAccessToken())
	require.Empty(events.AllAcquiredAllAccessTokens())

	// put acquired token event
	tokenName := "to1001"
	acquiredAccessToken := NewAcquiredAccessToken(tokenName)
	events.AddAcquiredAccessToken(acquiredAccessToken)

	select {
	case msg, ok := <-events.TokensChannel():
		require.Equal(acquiredAccessToken, msg)
		require.True(ok)
	case msg, ok := <-events.AllTokensChannel():
		require.Nil(msg)
		require.False(ok)
	case <-time.After(selectTimeout):
		require.FailNow("expected NewAcquiredAccessToken")
	}

	require.Equal([]AcquiredAccessToken{
		acquiredAccessToken,
	}, events.AllAcquiredAccessToken())
	require.Empty(events.AllAcquiredAllAccessTokens())

	// put all tokens acquired event
	tokenNames := []string{"to1001"}
	acquiredAllAccessTokens := NewAcquiredAllAccessTokens(tokenNames)
	events.AddAcquiredAllAccessTokens(acquiredAllAccessTokens)

	select {
	case msg, ok := <-events.TokensChannel():
		require.Nil(msg)
		require.False(ok)
	case msg, ok := <-events.AllTokensChannel():
		require.Equal(acquiredAllAccessTokens, msg)
		require.True(ok)
	case <-time.After(selectTimeout):
		require.FailNow("expected NewAcquiredAllAccessTokens")
	}

	require.Equal([]AcquiredAccessToken{
		acquiredAccessToken,
	}, events.AllAcquiredAccessToken())
	require.Equal([]AcquiredAllAccessTokens{
		acquiredAllAccessTokens,
	}, events.AllAcquiredAllAccessTokens())
}
