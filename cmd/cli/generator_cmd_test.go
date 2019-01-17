package main

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/server/mocks"
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
)

func ExampleGeneratorCommand_runNoFilename() {
	journeyMock := &mocks.Journey{}
	generator := newGeneratorCmdWrapperWithOptions(journeyMock)
	root := rootCommand(generator.run)

	_, err := executeCommand(root, "generate")
	if err != nil {
		fmt.Println(err.Error())
	}

	// Output:
	// You need to provide a discovery filename.
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
