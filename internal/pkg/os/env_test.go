package os

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const envKeyName = "FCS_TEST_ENV_VAR"

func TestGetEnvOrDefaultGetFromDefault(t *testing.T) {
	assert.Equal(t, "DEFAULT", GetEnvOrDefault(envKeyName, "DEFAULT"))
}

func TestGetEnvOrDefaultGetFromEnv(t *testing.T) {
	require.NoError(t, os.Setenv(envKeyName, "VALUE"))
	assert.Equal(t, "VALUE", GetEnvOrDefault(envKeyName, "DEFAULT_VALUE"))
	require.NoError(t, os.Unsetenv(envKeyName))
}

func TestGetEnvOrDefaultGetFromEnvButEmpty(t *testing.T) {
	require.NoError(t, os.Setenv(envKeyName, ""))
	assert.Equal(t, "", GetEnvOrDefault(envKeyName, "DEFAULT_VALUE"))
	require.NoError(t, os.Unsetenv(envKeyName))
}
