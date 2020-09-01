package manifest

import (
	"testing"

	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/assert"
)

// test to check that all the tests for the file are presented when

func TestAPIVersionsSimple(t *testing.T) {
	v201, err := semver.Make("2.0.1")
	assert.Nil(t, err)
	v300, err := semver.Make("3.0.0")
	assert.Nil(t, err)
	v312, err := semver.Make("3.1.2")
	assert.Nil(t, err)
	v314, err := semver.Make("3.1.4")
	assert.Nil(t, err)
	v317, err := semver.Make("3.1.7")
	assert.Nil(t, err)
	v400, err := semver.Make("4.0.0")
	assert.Nil(t, err)

	assert.True(t, v312.LT(v314))
	assert.True(t, v312.LT(v317))

	singleRange, err := semver.ParseRange("3.1.4")
	assert.True(t, singleRange(v314))

	multiRange, err := semver.ParseRange(">=3.1.4 <=3.1.8")
	assert.True(t, multiRange(v317))
	assert.False(t, multiRange(v312))
	assert.False(t, multiRange(v400))

	anotherRange, err := semver.ParseRange(">=3.1.0")
	assert.True(t, anotherRange(v400))
	assert.False(t, anotherRange(v300))
	assert.False(t, anotherRange(v201))

}
