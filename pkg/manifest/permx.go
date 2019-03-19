package manifest

import (
	"fmt"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
)

// TestCasePermission -
type TestCasePermission struct {
	ID     string   `json:"id,omitempty"`
	Perms  []string `json:"perms,omitempty"`
	Permsx []string `json:"permsx,omitempty"`
}

// TokenGatherer -
type TokenGatherer struct {
	Name   string   `json:"name,omitempty"`
	IDs    []string `json:"ids,omitempty"`
	Perms  []string `json:"perms,omitempty"`
	Permsx []string `json:"permsx,omitempty"`
}

// GetTestCasePermissions -
func GetTestCasePermissions(tcs []model.TestCase) ([]TestCasePermission, error) {
	tcps := []TestCasePermission{}
	for _, tc := range tcs {
		ctx := tc.Context
		perms, _ := ctx.GetStringSlice("permissions")
		permsx, _ := ctx.GetStringSlice("permissions-excluded")
		tcp := TestCasePermission{ID: tc.ID, Perms: perms, Permsx: permsx}
		tcps = append(tcps, tcp)
	}
	return tcps, nil
}

// GatherTokens - gathers all tokens
func GatherTokens(tcps []TestCasePermission) ([]TokenGatherer, error) {
	tg := []TokenGatherer{}
	for _, tcp := range tcps {
		mixer(&tcp, tg)
	}
	return tg, nil
}

func mixer(tcp *TestCasePermission, tg []TokenGatherer) {

	fmt.Printf("testcasepermissions: %#v\n", tcp)
	fmt.Printf("testcasepermissions: %#v\n", tcp)

	for _, tgItem := range tg {
		tcPermxConflict := false
		tcPermConflict := false
		for _, tgperm := range tgItem.Perms {
			for _, tcpermx := range tcp.Permsx {
				if tgperm == tcpermx {
					tcPermxConflict = true
					break
				}
			}
			if tcPermxConflict {
				break
			}
		}
		for _, tgpermx := range tgItem.Permsx {
			for _, tcperm := range tcp.Perms {
				if tgpermx == tcperm {
					tcPermConflict = true
					break
				}
			}
			if tcPermConflict {
				break
			}
		}

	}
}

// check included doesn't include excluded
// check exluded doesn't include included
// if ok add
// if not ok, check next
// if not added at end add
