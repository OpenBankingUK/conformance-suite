package schema

import (
	"time"

	"github.com/dlclark/regexp2"
	"github.com/getkin/kin-openapi/openapi3"
)

// matchTimeout is the maximum duration allowed for a single regexp2 match
// operation. A finite timeout prevents a pathological backtracking pattern
// from blocking a goroutine indefinitely (ReDoS).
const matchTimeout = 5 * time.Second

// regexp2Compiler is an openapi3.RegexCompilerFunc that uses the regexp2 engine,
// which supports PCRE features such as lookahead and lookbehind. Because RE2 is a
// strict subset of PCRE, all existing RE2-compatible patterns continue to compile
// and match identically.
func regexp2Compiler(expr string) (openapi3.RegexMatcher, error) {
	re, err := regexp2.Compile(expr, 0)
	if err != nil {
		return nil, err
	}
	re.MatchTimeout = matchTimeout
	return &regexp2Matcher{re: re}, nil
}

// regexp2Matcher wraps a *regexp2.Regexp to satisfy openapi3.RegexMatcher.
type regexp2Matcher struct {
	re *regexp2.Regexp
}

func (m *regexp2Matcher) MatchString(s string) bool {
	matched, err := m.re.MatchString(s)
	if err != nil {
		return false
	}
	return matched
}
