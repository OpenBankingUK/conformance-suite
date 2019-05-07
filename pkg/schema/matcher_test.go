package schema

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatcher(t *testing.T) {
	tcs := []struct {
		pathWithParam string
		path          string
		isMatch       bool
	}{
		{"/accounts/{account_id}", "/accounts/1234567890", true},
	}

	m := NewMatcher()

	for _, tc := range tcs {
		testName := fmt.Sprintf("`%s`%s`%s`", tc.pathWithParam, equalStr[tc.isMatch], tc.path)
		t.Run(testName, func(t *testing.T) {
			assert.Equal(t, tc.isMatch, m.Match(tc.pathWithParam, tc.path))
		})
	}
}

var equalStr = map[bool]string{
	true:  "=",
	false: "!=",
}
