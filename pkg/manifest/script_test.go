package manifest

import (
	"fmt"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestGenerateTestCases(t *testing.T) {
	tests, err := GenerateTestCases("TestSpec", "http://mybaseurl", &model.Context{})
	assert.Nil(t, err)

	perms, err := GetPermissions(tests)
	assert.Nil(t, err)
	m := make(map[string]string, 0)
	for _, v := range perms {
		fmt.Printf("perms: %s %-50.50s %s\n", v.ID, v.Path, v.Permissions)
		m[v.Path] = v.ID
	}
	fmt.Println("----------------------==")
	for k := range m {
		fmt.Println(k)
	}

	for _, v := range tests {
		dumpJSON(v)
	}

}
