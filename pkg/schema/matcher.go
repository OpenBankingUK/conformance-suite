package schema

import (
	"regexp"
)

// Matcher exposes a path comparison interface
// Expected usage to compare a path that has placeholder and a
// real path. Ex: /accounts/{account_id} == /accounts/1234567890
type Matcher interface {
	Match(pathWithParams, path2 string) bool
}

type paramMatcher struct{}

func NewMatcher() Matcher {
	return paramMatcher{}
}

// matches a param in a URL in format `{AccountId}`
var r = regexp.MustCompile("{[a-zA-Z0-9_]+}")

func (m paramMatcher) Match(pathWithParams, path string) bool {

	pathWithParams = r.ReplaceAllString(pathWithParams, `[^/]+`) + `$`

	rr, err := regexp.Compile(pathWithParams)
	if err != nil {
		return false
	}

	return rr.MatchString(path)
}
