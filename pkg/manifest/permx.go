package manifest

import "bitbucket.org/openbankingteam/conformance-suite/pkg/model"

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

// TokenStore eats tokens
type TokenStore struct {
	store []TokenGatherer
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
	tg := make([]TokenGatherer, 0)
	te := TokenStore{}
	for _, tcp := range tcps { // for each TestCasePermission
		te.createOrUpdate(tcp)
		//tg = append(tg, *newItem)
	}
	dumpJSON(te.store)
	return tg, nil
}

func dumpTG(tg []TokenGatherer) {
	for _, v := range tg {
		fmt.Printf("grouplineitem: %v - %v -  %v\n", v.IDs, v.Perms, v.Permsx)
	}
}

// create or update TokenGethereer
func (te *TokenStore) createOrUpdate(tcp TestCasePermission) {

	if len(te.store) == 0 { // First time - no permissions - just add
		tpg := TokenGatherer{IDs: []string{tcp.ID}, Perms: tcp.Perms, Permsx: tcp.Permsx}
		te.store = append(te.store, tpg)
		return
	}

	for idx, tgItem := range te.store { // loop through each Gathered Item
		tcPermxConflict := false
		tcPermConflict := false

		// Check groupPermissions against testcaseExclusions
		for _, tgperm := range tgItem.Perms { // loop through all
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
		if tcPermxConflict { //move onto next group item
			continue
		}

		// Check groupExclusions against testcasePermissions
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
		if tcPermConflict {
			continue
		}
		newItem := addPermToGathererItem(tcp, tgItem)
		te.store[idx] = newItem
		return
	}
	tpg := TokenGatherer{IDs: []string{tcp.ID}, Perms: tcp.Perms, Permsx: tcp.Permsx}
	te.store = append(te.store, tpg)

	return
}

func addPermToGathererItem(tp TestCasePermission, tg TokenGatherer) TokenGatherer {
	var newPerm, newPermx string
	var permMatch, permxMatch bool

	tg.IDs = append(tg.IDs, tp.ID)
	for _, tgPerm := range tg.Perms { // Loop through all Gathered Permissions
		for _, tpPerm := range tp.Perms { //
			if tpPerm == tgPerm {
				permMatch = true
				newPerm = tpPerm
				break
			}
		}
		if permMatch {
			break
		} else {
			if newPerm != "" {
				tg.Perms = append(tg.Perms, newPerm)
			}
		}
	}

	for _, tgPermx := range tg.Permsx {
		for _, tpPermx := range tp.Permsx {
			if tpPermx == tgPermx {
				permxMatch = true
				newPermx = tpPermx
				break
			}
		}
		if permxMatch {
			break
		} else {
			if newPermx != "" {
				tg.Permsx = append(tg.Permsx, newPermx)
			}
		}
	}
	return tg
}
