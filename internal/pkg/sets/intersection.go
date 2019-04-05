package sets

import "strings"

func InsensitiveIntersection(a, b []string) []string {
	return intersection(a, b, insensitiveEquals)
}

func Intersection(a, b []string) []string {
	return intersection(a, b, equals)
}

func intersection(a, b []string, equals func(a string, b string) bool) []string {
	set := make([]string, 0)
	for _, el := range a {
		if contains(b, el, equals) {
			set = append(set, el)
		}
	}
	return set
}

func contains(a []string, e string, equals func(a string, b string) bool) bool {
	for _, value := range a {
		if equals(value, e) {
			return true
		}
	}
	return false
}

func equals(a, b string) bool {
	return a == b
}

func insensitiveEquals(a, b string) bool {
	return strings.ToUpper(a) == strings.ToUpper(b)
}
