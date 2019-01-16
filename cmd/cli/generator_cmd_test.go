package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func ExampleGeneratorCommand_runNoFilename() {
	generatorCommand := newGeneratorCmdWrapper()
	generatorCommand.run(&cobra.Command{}, nil)
	// Output:
	// FCS - Functional Conformance Suite
	// you need to provide a discovery filename
}

func ExampleGeneratorCommand_run() {
	generatorCommand := newGeneratorCmdWrapper()
	cobraCmd := &cobra.Command{}
	cobraCmd.Flags().String("filename", "dist/ob-v3.0-generic.json", "Discovery filename")
	err := cobraCmd.ParseFlags([]string{"--filename", "file"})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	generatorCommand.run(cobraCmd, nil)
	// Output:
	// FCS - Functional Conformance Suite
	// it will generate from file
}
