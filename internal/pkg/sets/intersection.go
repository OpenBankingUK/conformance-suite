package sets

func Intersection(a, b []string) []string {
	set := make([]string, 0)
	for _, el := range a {
		if contains(b, el) {
			set = append(set, el)
		}
	}
	return set
}

func contains(a []string, e string) bool {
	for _, value := range a {
		if value == e {
			return true
		}
	}
	return false
}
