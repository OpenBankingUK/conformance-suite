package main

import (
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestConfigWelcome(t *testing.T) {
	expected := "I am about to—or I am going to—die. Either expression is correct."
	require.NoError(t, os.Setenv("FCS_WELCOME", expected))

	config, err := readViperEnvConfig()

	require.NoError(t, err)
	assert.Equal(t, expected, config.Welcome)
}
