package executors

import (
	"errors"
	"sync"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/executors/events"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"

	"github.com/sirupsen/logrus"
)

// Collector collects tokens for a set or permissions requirements and calls a
// subscribed function when it has received all
type Collector interface {
	Collect(setName, token string) error
	Tokens() []Token
}

// Token represents a token acquisition for one of the permission set requirement
type Token struct {
	Code            string
	NamedPermission model.NamedPermission
}

// TokenConsentIDs captures the token/consentIds awaiting authorisation
type TokenConsentIDs []TokenConsentIDItem

// TokenConsentIDItem is a single consentId mapping to token name
type TokenConsentIDItem struct {
	TokenName   string
	ConsentID   string
	Permissions string
	AccessToken string
	ConsentURL  string
	Error       string
}

// TokenCollector - collects tokens
type TokenCollector interface {
	Collect(tokenName, accesstoken string) error
	Tokens() TokenConsentIDs
}

type tokenCollector struct {
	tokensLock   *sync.Mutex
	collected    int
	doneFunc     func()
	consentTable TokenConsentIDs
	log          *logrus.Entry
	events       events.Events
}

// NewTokenCollector -
func NewTokenCollector(log *logrus.Entry, consentIds TokenConsentIDs, doneFunc func(), events events.Events) TokenCollector {
	return &tokenCollector{
		tokensLock:   &sync.Mutex{},
		collected:    0,
		doneFunc:     doneFunc,
		consentTable: consentIds,
		log:          log.WithField("module", "tokenCollector"),
		events:       events,
	}
}

// Collect receives an accesstoken to match a named token for which we have a consentid
func (c *tokenCollector) Collect(tokenName, accessToken string) error {
	c.log.Debug("tokenCollector.Collect, tokensLock=false")
	c.tokensLock.Lock()
	c.log.Debug("tokenCollector.Collect, tokensLock=true")
	defer func() {
		c.log.Debug("tokenCollector.Collect, tokensLock=false")
		c.tokensLock.Unlock()
	}()

	c.log.Debugf("entry with tokenname %s: accessToken: %s", tokenName, accessToken)
	if !c.tokenNameExists(tokenName) {
		return errors.New("invalid token name: " + tokenName)
	}

	c.addAccessToken(tokenName, accessToken)
	c.log.Debugf("collected: %d, total to be collected: %d", c.collected, len(c.consentTable))
	if c.isDone() {
		c.log.Debug("calling done!")
		tokenNames := []string{}
		for _, item := range c.consentTable {
			tokenNames = append(tokenNames, item.TokenName)
		}

		acquiredAllAccessTokens := events.NewAcquiredAllAccessTokens(tokenNames)
		c.events.AddAcquiredAllAccessTokens(acquiredAllAccessTokens)

		c.doneFunc()
	}

	return nil
}

func (c *tokenCollector) Tokens() TokenConsentIDs {
	c.log.Debug("tokenCollector.Tokens, tokensLock=false")
	c.tokensLock.Lock()
	c.log.Debug("tokenCollector.Collect, tokensLock=true")
	defer func() {
		c.log.Debug("tokenCollector.Tokens, tokensLock=false")
		c.tokensLock.Unlock()
	}()

	return c.consentTable
}

func (c *tokenCollector) tokenNameExists(tokenName string) bool {
	for _, item := range c.consentTable {
		if item.TokenName == tokenName {
			return true
		}
	}
	return false
}

func (c *tokenCollector) addAccessToken(tokenName, accessToken string) {
	for k, item := range c.consentTable {
		if tokenName == item.TokenName {
			item.AccessToken = accessToken
			c.consentTable[k] = item
			c.collected++

			acquiredAccessToken := events.NewAcquiredAccessToken(tokenName)
			c.events.AddAcquiredAccessToken(acquiredAccessToken)
		}
	}
}

func (c *tokenCollector) isDone() bool {
	return c.collected == len(c.consentTable)
}
