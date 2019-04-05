package sets

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntersection(t *testing.T) {
	setA := []string{"A", "B", "C"}
	setB := []string{"A", "C", "D"}

	result := Intersection(setA, setB)

	expected := []string{"A", "C"}
	assert.Equal(t, expected, result)
}

func TestInsensitiveIntersection(t *testing.T) {
	setA := []string{"A", "B", "C"}
	setB := []string{"A", "c", "D"}

	result := InsensitiveIntersection(setA, setB)

	expected := []string{"A", "C"}
	assert.Equal(t, expected, result)
}
