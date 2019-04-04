package model

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCalcKid(t *testing.T) {
	modulus := "tPS6lkBEFf9MuVEfK7BET0oLYe7r6QjQR1SzXqwm37TmcnB8koB66ExmeFizSl8eJuTTjsNCDliGqbGdoe8p_Xw4hRLAPqtEEbq1-sQAAwPUHwgyAABOhIlWBsI6KxYX20UCp5pR4EzqM5cEj_nIvCjw7lmXZaOasMis9utAMw3iKFitduNS5Mj0g523CAes6CnlKusYf--k2l4TpgFRiYFGdVb7T-07xAqlyo5ljLguu8Tz_iwLaqvFKjb_m8gO7dy8P3h8wCv_nbdntxzh17EzsXiMIyh3PNKmxJUmUoAuKkOkpzaRVB5NsjIguIGZrrv0k_hZirxGA_SobsxAvQ"
	
	kid, err := calcKid(modulus)

	require.NoError(t, err)
	expected := "X3idGb9VFwA3FK101sgNnaHmM2Y"
	assert.Equal(t, expected, kid)
}
