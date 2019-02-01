package permissions

type group struct {
	included CodeSet
	excluded CodeSet
}

type groupSet []*group

func (gs groupSet) firstCompatible(aGroup *group) (*group, bool) {
	for _, group := range gs {
		if aGroup.isCompatible(group) && group.isCompatible(aGroup) {
			return group, true
		}
	}
	return nil, false
}

func (g *group) isSatisfiedByAnyOf(groups []*group) bool {
	for _, group := range groups {
		if g.isSatisfiedBy(group) {
			return true
		}
	}
	return false
}

func (g *group) isSatisfiedBy(group *group) bool {
	if !group.included.HasAll(g.included) {
		return false
	}

	if group.included.HasAny(g.excluded) {
		return false
	}

	return true
}

func (g *group) isCompatible(group *group) bool {
	if group.included.HasAny(g.excluded) {
		return false
	}

	if g.excluded.HasAny(group.included) {
		return false
	}

	return true
}

func (g *group) Union(g2 *group) {
	g.included = g.included.Union(g2.included)
	g.excluded = g.excluded.Union(g2.excluded)
}
