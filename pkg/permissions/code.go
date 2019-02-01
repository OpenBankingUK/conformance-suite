package permissions

// Code is a string representing a OB access permission
type Code string

type CodeSet []Code

var NoCodeSet []CodeSet

func (c CodeSet) Has(searchCode Code) bool {
	for _, code := range c {
		if code == searchCode {
			return true
		}
	}
	return false
}

func (c CodeSet) HasAll(codes []Code) bool {
	for _, code := range codes {
		if !c.Has(code) {
			return false
		}
	}
	return true
}

func (c CodeSet) Equals(codes CodeSet) bool {
	if len(codes) != len(c) {
		return false
	}

	if !c.HasAll(codes) {
		return false
	}

	if !codes.HasAll(c) {
		return false
	}

	return true
}

func (c CodeSet) HasAny(codes []Code) bool {
	for _, code := range codes {
		if c.Has(code) {
			return true
		}
	}
	return false
}

func (c CodeSet) Union(c2 CodeSet) CodeSet {
	union := CodeSet{}
	for _, code := range c {
		if !union.Has(code) {
			union = append(union, code)
		}
	}
	for _, code := range c2 {
		if !union.Has(code) {
			union = append(union, code)
		}
	}
	return union
}
