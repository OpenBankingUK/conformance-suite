package events

import (
	"testing"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
)

func TestEvents(t *testing.T) {
	selectTimeout := 1 * time.Millisecond
	require := test.NewRequire(t)

	events := NewEvents()
	{
		// initially empty
		select {
		case event, ok := <-events.Tokens():
			require.Nil(event)
			require.False(ok)
		case event, ok := <-events.AllTokens():
			require.Nil(event)
			require.False(ok)
		case <-time.After(selectTimeout):
			break
		}
	}
	{
		// put acquired token event
		tokenName := "to1001"
		acquiredAccessTokenEvent := NewAcquiredAccessToken(tokenName)
		events.Tokens() <- acquiredAccessTokenEvent

		select {
		case event, ok := <-events.Tokens():
			require.Equal(acquiredAccessTokenEvent, event)
			require.True(ok)
		case event, ok := <-events.AllTokens():
			require.Nil(event)
			require.False(ok)
		case <-time.After(selectTimeout):
			require.FailNow("expected NewAcquiredAccessToken")
		}
	}
	{
		// put all tokens acquired event
		tokenNames := []string{"to1001"}
		acquiredAllAccessTokensEvent := NewAcquiredAllAccessTokens(tokenNames)
		events.AllTokens() <- acquiredAllAccessTokensEvent

		select {
		case event, ok := <-events.Tokens():
			require.Nil(event)
			require.False(ok)
		case event, ok := <-events.AllTokens():
			require.Equal(acquiredAllAccessTokensEvent, event)
			require.True(ok)
		case <-time.After(selectTimeout):
			require.FailNow("expected NewAcquiredAllAccessTokens")
		}
	}
}
