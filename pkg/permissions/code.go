package permissions

// Code is a string representing a OB access permission
type Code string

// CodeSet is a set of OB code permissions
type CodeSet []Code

// NoCodeSet represents empty or no permissions set
func NoCodeSet() CodeSet {
	return CodeSet{}
}

// Has check if a set Has a code
func (c CodeSet) Has(searchCode Code) bool {
	for _, code := range c {
		if code == searchCode {
			return true
		}
	}
	return false
}

// HasAll check is a set has all codes in other set
func (c CodeSet) HasAll(otherSet CodeSet) bool {
	for _, code := range otherSet {
		if !c.Has(code) {
			return false
		}
	}
	return true
}

// Equals check if 2 sets have the SAME codes
func (c CodeSet) Equals(otherSet CodeSet) bool {
	if len(otherSet) != len(c) {
		return false
	}

	if !c.HasAll(otherSet) {
		return false
	}

	return true
}

// HasAny checks if has any of the codes of other set
func (c CodeSet) HasAny(otherSet CodeSet) bool {
	for _, code := range otherSet {
		if c.Has(code) {
			return true
		}
	}
	return false
}

// add returns a new set with all Code from 2 sets
func (c CodeSet) Union(otherSet CodeSet) CodeSet {
	union := CodeSet{}
	for _, code := range c {
		if !union.Has(code) {
			union = append(union, code)
		}
	}
	for _, code := range otherSet {
		if !union.Has(code) {
			union = append(union, code)
		}
	}
	return union
}
