package executors

import (
	"errors"
	"sync"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
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

type collector struct {
	requirements []model.SpecConsentRequirements
	tokensLock   *sync.Mutex
	tokens       map[string]string
	doneFunc     func()
	consentIDs   TokenConsentIDs
}

// NewCollector returns a thread safe token collector
func NewCollector(requirements []model.SpecConsentRequirements, doneFunc func()) *collector {
	return &collector{
		requirements: requirements,
		tokensLock:   &sync.Mutex{},
		tokens:       map[string]string{},
		doneFunc:     doneFunc,
	}
}

// Collect receives on token for a permission set
func (c *collector) Collect(setName, token string) error {
	if !c.setNameExists(setName) {
		return errors.New("invalid permission set name")
	}
	c.tokensLock.Lock()
	c.tokens[setName] = token
	if c.isDone() {
		c.doneFunc()
	}
	c.tokensLock.Unlock()
	return nil
}

// Tokens retrieves all collected tokens
func (c *collector) Tokens() []Token {
	c.tokensLock.Lock()
	var result []Token
	for _, spec := range c.requirements {
		for _, np := range spec.NamedPermissions {
			result = append(result, Token{
				Code:            c.tokens[np.Name],
				NamedPermission: np,
			})
		}
	}
	c.tokensLock.Unlock()
	return result
}

// isDone checks if we have as many tokens as permission sets required
func (c *collector) isDone() bool {
	// naive simply count the tokens collected against named permission sets
	tokensRequired := c.countNamedSets()
	return tokensRequired == len(c.tokens)
}

func (c *collector) countNamedSets() int {
	totalTokens := 0
	for _, spec := range c.requirements {
		totalTokens += len(spec.NamedPermissions)
	}
	return totalTokens
}

// setNameExists checks if a setNamed permission exists in the requirements
func (c *collector) setNameExists(setName string) bool {
	for _, spec := range c.requirements {
		for _, np := range spec.NamedPermissions {
			if np.Name == setName {
				return true
			}
		}
	}
	return false
}

type nullCollector struct {
}

// NewNullCollector implements a dummy collector that trigger done immediately and collects nothing
// for using when we don't want to collect or in tests
func NewNullCollector(doneFunc func()) Collector {
	go doneFunc()
	return nullCollector{}
}

// Collect receives on token for a permission set
func (c nullCollector) Collect(setName, token string) error {
	return errors.New("cant collect this is a null collector")
}

// Tokens retrieves all collected tokens
func (c nullCollector) Tokens() []Token {
	return []Token{}
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
}

// TokenCollector - collects tokens
type TokenCollector interface {
	Collect(tokenName, accesstoken string) error
	Tokens() TokenConsentIDs
}

type tokencollector struct {
	tokensLock   *sync.Mutex
	collected    int
	doneFunc     func()
	consentTable TokenConsentIDs
}

// NewTokenCollector -
func NewTokenCollector(consentIds TokenConsentIDs, doneFunc func()) *tokencollector {
	return &tokencollector{
		tokensLock:   &sync.Mutex{},
		collected:    0,
		doneFunc:     doneFunc,
		consentTable: consentIds,
	}
}

// Collect receives an accesstoken to match a named token for which we have a consentid
func (c *tokencollector) Collect(tokenName, accessToken string) error {
	if !c.tokenNameExists(tokenName) {
		return errors.New("invalid token name: " + tokenName)
	}
	c.tokensLock.Lock()
	c.addAccessToken(tokenName, accessToken)
	if c.isDone() {
		c.doneFunc()
	}
	c.tokensLock.Unlock()
	return nil
}

func (c *tokencollector) Tokens() TokenConsentIDs {
	return c.consentTable
}

func (c *tokencollector) tokenNameExists(tokenName string) bool {
	for _, item := range c.consentTable {
		if item.TokenName == tokenName {
			return true
		}
	}
	return false
}

func (c *tokencollector) addAccessToken(tokenName, accessToken string) {
	for k, item := range c.consentTable {
		if tokenName == item.TokenName {
			item.AccessToken = accessToken
			c.consentTable[k] = item
			c.collected++
		}
	}
}

func (c *tokencollector) isDone() bool {
	return c.collected == len(c.consentTable)
}
