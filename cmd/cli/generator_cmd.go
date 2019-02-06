package main

import (
	"fmt"
	"os"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/server"
	"github.com/spf13/cobra"
)

func generatorCmd(runFunc cobraCmdRunFunc) *cobra.Command {
	generatorCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate test cases from a discovery model",
		Long:  "Generated test cases will output to standard output.",
		Run:   runFunc,
	}
	generatorCmd.Flags().StringP("filename", "f", "", "Discovery filename")
	generatorCmd.Flags().StringP("output", "o", "", "Output filename, defaults to stdout")

	return generatorCmd
}

// GeneratorCommand executes a discovery model to get test cases
type GeneratorCommand struct {
	Generator
}

func newGeneratorCmdWrapper() GeneratorCommand {
	checker := model.NewConditionalityChecker()
	validatorEngine := discovery.NewFuncValidator(checker)
	testGenerator := generation.NewGenerator()
	journey := server.NewJourney(testGenerator, validatorEngine)
	generator := newGenerator(journey)
	return newGeneratorCmdWrapperWithOptions(generator)
}

func newGeneratorCmdWrapperWithOptions(generator Generator) GeneratorCommand {
	return GeneratorCommand{generator}
}

func (g GeneratorCommand) run(cmd *cobra.Command, _ []string) {
	// check if input (discovery model) filename if provided
	filenameFlag, err := cmd.Flags().GetString("filename")
	if err != nil || filenameFlag == "" {
		fmt.Println("You need to provide a discovery filename.")
		return
	}

	// set where to write results (testcases) defaults to stdout, flag output to choose a file
	output := os.Stdout
	outputFlag, err := cmd.Flags().GetString("output")
	if err == nil && outputFlag != "" {
		file, err := os.Create(outputFlag)
		if err != nil {
			exitError(err, "Error creating output file")
		}
		output = file
		defer func() {
			err := file.Close()
			if err != nil {
				exitError(err, "Error closing output file")
			}
		}()
	}

	input, err := os.Open(filenameFlag)
	if err != nil {
		exitError(err, "Error running generation command, opening input file")
	}
	defer func() {
		err := input.Close()
		if err != nil {
			exitError(err, "Error closing output file")
		}
	}()

	err = g.Generate(input, output)
	if err != nil {
		exitError(err, "Error running generation command")
	}
}

func exitError(err error, message string) {
	fmt.Fprint(os.Stderr, message+"\n")
	fmt.Fprint(os.Stderr, err.Error()+"\n")
	os.Exit(1)
}
