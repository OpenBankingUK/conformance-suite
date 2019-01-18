package main

import (
	"bitbucket.org/openbankingteam/conformance-suite/cmd/cli/mocks"
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

func ExampleGeneratorCommand_runNoFilename() {
	generator := &mocks.Generator{}
	generatorCmdWrapper := newGeneratorCmdWrapperWithOptions(generator)
	root := rootCommand(generatorCmdWrapper.run)

	_, err := executeCommand(root, "generate")
	if err != nil {
		fmt.Println(err.Error())
	}

	// Output:
	// You need to provide a discovery filename.
}

func TestGeneratorCommand(t *testing.T) {
	generator := &mocks.Generator{}
	generator.On("Generate", mock.Anything, mock.Anything).Return(nil)
	generatorCmdWrapper := newGeneratorCmdWrapperWithOptions(generator)
	root := rootCommand(generatorCmdWrapper.run)

	_, err := executeCommand(root, "generate", "--filename", "config.go")

	require.NoError(t, err)
}

func TestGeneratorWritesToFile(t *testing.T) {
	generator := &mocks.Generator{}
	generator.On(
		"Generate", mock.Anything, mock.Anything).Return(nil)
	generatorCmdWrapper := newGeneratorCmdWrapperWithOptions(generator)
	root := rootCommand(generatorCmdWrapper.run)
	output := tempFileName("fcs", "json")

	_, err := executeCommand(
		root,
		"generate",
		"--filename",
		"testdata/discovery-model.json",
		"--output",
		output,
	)

	require.NoError(t, err)
	assert.FileExists(t, output)
	require.NoError(t, os.Remove(output))
}

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOutput(buf)
	root.SetArgs(args)
	c, err = root.ExecuteC()
	return c, buf.String(), err
}

func tempFileName(prefix, suffix string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix)
}
