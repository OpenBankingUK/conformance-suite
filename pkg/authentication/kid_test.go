package authentication

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalcKid(t *testing.T) {
	modulus := "tGzvc5H2KLufptikvbL1crtdSaV901mJY4dAxjWK2V-W6hhgNIgdQgusn3k8AW6KKFckDLIs0hYKmIJTVN0MGaruG4USN4sRlRT2kkizJaXU9ZtHZ5yiwP9BMEiaKgY6IGWy4vVxR9ii83HhAXbTo-gI9HaK73i2kLIYUYwiAUG32Oo5Z226dISMBiGxDU7EeLCJ8uhdKPTi05z5fPE0Lw3eszLwaJN8qQ1BIFON_QXCVS7BDMdmWh2XEEljD_h5d6W1SPXikWod2XWK9PbxbKzGkpIJHV_Ty74c48eQE3_0rkUEZ9iCHtuFxgN0SEy1Hj5-5TDMVXkVQO_rGyYv4w"

	kid, err := CalcKid(modulus)

	require.NoError(t, err)
	expected := "QuFYBRJnWdI6_NHFgamuXNr5R20"
	assert.Equal(t, expected, kid)
}
