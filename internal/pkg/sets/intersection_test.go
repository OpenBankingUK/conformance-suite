package sets

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestIntersection(t *testing.T) {
	setA := []string{"A", "B", "C"}
	setB := []string{"A", "C", "D"}

	result := Intersection(setA, setB)

	expected := []string{"A", "C"}
	assert.Equal(t, expected, result)
}
