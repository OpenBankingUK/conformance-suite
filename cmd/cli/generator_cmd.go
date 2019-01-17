package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func generatorCmd() *cobra.Command {
	generatorCmdWrapper := newGeneratorCmdWrapper()
	generatorCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate test cases from a discovery model",
		Run:   generatorCmdWrapper.run,
	}
	generatorCmd.Flags().String("filename", "dist/ob-v3.0-generic.json", "Discovery filename")
	return generatorCmd
}

// GeneratorCommand executes a discovery model to get test cases
type GeneratorCommand struct {
}

func newGeneratorCmdWrapper() GeneratorCommand {
	return GeneratorCommand{}
}

func (v GeneratorCommand) run(cmd *cobra.Command, _ []string) {
	fmt.Println(banner)
	filenameFlag := cmd.Flag("filename")
	if filenameFlag == nil {
		fmt.Println("you need to provide a discovery filename")
		return
	}
	fmt.Printf("it will generate from %s\n", filenameFlag.Value)
}
