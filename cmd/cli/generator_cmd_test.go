package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
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

func TestGeneratorWritesToFile(t *testing.T) {
	generator := &MockGenerator{}
	generator.On("Generate", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	generatorCmdWrapper := newGeneratorCmdWrapperWithOptions(generator)
	root := newRootCommand(generatorCmdWrapper.run)
	output, err := tempFileName("fcs", ".json")
	require.NoError(t, err)

	_, err = executeCommand(
		root,
		"generate",
		"--filename",
		"testdata/discovery-model.json",
		"--output",
		output,
		"--config",
		"testdata/config.json",
	)

	require.NoError(t, err)
	require.FileExists(t, output)
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

func tempFileName(prefix, suffix string) (string, error) {
	randBytes := make([]byte, 16)
	_, err := rand.Read(randBytes)
	if err != nil {
		return "", err
	}
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix), nil
}
