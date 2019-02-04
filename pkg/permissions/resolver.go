package permissions

// Resolver find minimal codeSet required to satisfy a set group of permissions (endpoints)
func Resolver(groups []group) []CodeSet {
	if len(groups) == 0 {
		return NoCodeSet
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

	return mapToCodeSets(groupsFound)
}

func mapToCodeSets(groups []*group) []CodeSet {
	var codeSets []CodeSet
	for _, group := range groups {
		codeSets = append(codeSets, group.included)
	}
	return codeSets
}
