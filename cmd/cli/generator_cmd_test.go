package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func ExampleGeneratorCommand_runNoFilename() {
	generator := &MockGenerator{}
	generatorCmdWrapper := newGeneratorCmdWrapperWithOptions(generator)
	root := newRootCommand(generatorCmdWrapper.run)

	_, err := executeCommand(root, "generate")
	if err != nil {
		fmt.Println(err.Error())
	}

	// Output:
	// You need to provide a discovery filename.
}

func TestGeneratorCommand(t *testing.T) {
	generator := &MockGenerator{}
	generator.On("Generate", mock.Anything, mock.Anything).Return(nil)
	generatorCmdWrapper := newGeneratorCmdWrapperWithOptions(generator)
	root := newRootCommand(generatorCmdWrapper.run)

	_, err := executeCommand(root, "generate", "--filename", "generator_cmd_test.go")

	require.NoError(t, err)
}

func TestGeneratorWritesToFile(t *testing.T) {
	generator := &MockGenerator{}
	generator.On(
		"Generate", mock.Anything, mock.Anything).Return(nil)
	generatorCmdWrapper := newGeneratorCmdWrapperWithOptions(generator)
	root := newRootCommand(generatorCmdWrapper.run)
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
