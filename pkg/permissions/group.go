package permissions

// Group represents a test and it's context permission set config
type Group struct {
	TestId   TestId
	Included CodeSet
	Excluded CodeSet
}

func NewGroup(testId string, included, excluded CodeSet) Group {
	return Group{
		TestId:   TestId(testId),
		Included: included,
		Excluded: excluded,
	}
}

type groupSet []*Group

func (gs groupSet) firstCompatible(aGroup *Group) (*Group, bool) {
	for _, group := range gs {
		if aGroup.isCompatible(group) && group.isCompatible(aGroup) {
			return group, true
		}
	}
	return nil, false
}

func (g *Group) isSatisfiedByAnyOf(groups []*Group) bool {
	for _, group := range groups {
		if g.isSatisfiedBy(group) {
			return true
		}
	}
	return false
}

func (g *Group) isSatisfiedBy(group *Group) bool {
	if !group.Included.HasAll(g.Included) {
		return false
	}

	if group.Included.HasAny(g.Excluded) {
		return false
	}

	return true
}

func (g *Group) isCompatible(group *Group) bool {
	return !group.Included.HasAny(g.Excluded)
}

func (g *Group) add(g2 *Group) {
	g.Included = g.Included.Union(g2.Included)
	g.Excluded = g.Excluded.Union(g2.Excluded)
}
