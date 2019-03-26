package permissions

type TestId string

// CodeSetResult represents one set of permissions that are valid for a set of test ids
type CodeSetResult struct {
	CodeSet CodeSet  `json:"codes"`
	TestIds []TestId `json:"testIds"`
}

// CodeSetResultSet represents all permissions sets and their respective test id
type CodeSetResultSet []CodeSetResult

// Resolver find minimal codeSet required to satisfy a set Group of permissions (endpoints)
func Resolver(groups []Group) CodeSetResultSet {
	if len(groups) == 0 {
		return nil
	}

	var groupsFound groupSet
	for _, config := range groups {

		if len(groupsFound) == 0 {
			newGroup := config
			groupsFound = append(groupsFound, &newGroup)
			continue
		}

		if config.isSatisfiedByAnyOf(groupsFound) {
			continue
		}

		group, found := groupsFound.firstCompatible(&config)
		if !found {
			newGroup := config
			groupsFound = append(groupsFound, &newGroup)
			continue
		}

		group.add(&config)
	}

	return mapToCodeSets(groups, groupsFound)
}

// StringSliceToTestID -
func StringSliceToTestID(s []string) []TestId {
	tids := make([]TestId, 0)
	for _, v := range s {
		tids = append(tids, TestId(v))
	}
	return tids
}

// StringSliceToCodeSet -
func StringSliceToCodeSet(s []string) CodeSet {
	var cs CodeSet
	for _, v := range s {
		cs = append(cs, Code(v))
	}
	return cs

}

// mapToCodeSets maps all permission groups found to results that include test id list
// for each group found
func mapToCodeSets(groups []Group, groupsFound []*Group) []CodeSetResult {
	var codeSets CodeSetResultSet
	for _, groupFound := range groupsFound {
		// find tests that are satisfied by this Group
		for _, group := range groups {
			if group.isSatisfiedBy(groupFound) {
				codeSets.addTestToGroupFound(group, groupFound)
			}
		}
	}
	return codeSets
}

// addTestToGroupFound finds a permission set for a test and adds it to the test id list
// if doesnt find adds a new permission set
func (cs *CodeSetResultSet) addTestToGroupFound(group Group, groupFound *Group) {
	for key, codeSet := range *cs {
		if codeSet.CodeSet.Equals(groupFound.Included) {
			thisCs := *cs
			thisCs[key] = CodeSetResult{
				CodeSet: groupFound.Included,
				TestIds: append(codeSet.TestIds, group.TestId),
			}
			return
		}
	}
	// not found codeSet doesnt have this Group
	*cs = append(*cs, CodeSetResult{
		CodeSet: groupFound.Included,
		TestIds: []TestId{group.TestId},
	})
}
