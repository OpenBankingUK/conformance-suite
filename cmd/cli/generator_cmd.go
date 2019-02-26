package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
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
	generatorCmd.Flags().StringP("config", "c", "", "Config filename")
	generatorCmd.Flags().StringP("output", "o", "", "Output filename, defaults to stdout")

	return generatorCmd
}

// GeneratorCommand executes a discovery model to get test cases
type GeneratorCommand struct {
	Generator
}

func newGeneratorCmdWrapper(logger *logrus.Entry) GeneratorCommand {
	checker := model.NewConditionalityChecker()
	validatorEngine := discovery.NewFuncValidator(checker)
	testGenerator := generation.NewGenerator()
	// we overwrite the output of logs so it wont collide with testcases output (stdOut)
	logger.Logger.Out = os.Stderr
	journey := server.NewJourney(logger, testGenerator, validatorEngine)
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

	config, err := makeJourneyConfigFromFlag(cmd)
	if err != nil {
		exitError(err, "Error making config")
		return
	}

	// set where to write results (testcases) defaults to stdout, flag output to choose a file
	output := os.Stdout
	outputFlag, err := cmd.Flags().GetString("output")
	if err == nil && outputFlag != "" {
		file, err2 := os.Create(outputFlag)
		if err2 != nil {
			exitError(err2, "Error creating output file")
		}
		output = file
		defer func() {
			err3 := file.Close()
			if err3 != nil {
				exitError(err3, "Error closing output file")
			}
		}()
	}

	input, err := os.Open(filenameFlag)
	if err != nil {
		exitError(err, "Error running generation command, opening input file")
	}
	defer func() {
		err2 := input.Close()
		if err2 != nil {
			exitError(err2, "Error closing output file")
		}
	}()

	err = g.Generate(config, input, output)
	if err != nil {
		exitError(err, "Error running generation command")
	}
}

func makeJourneyConfigFromFlag(cmd *cobra.Command) (server.JourneyConfig, error) {
	// check if config filename if provided
	configFlag, err := cmd.Flags().GetString("config")
	if err != nil || configFlag == "" {
		return server.JourneyConfig{}, errors.New("you need to provide a config filename.")
	}

	configContent, err := ioutil.ReadFile(configFlag)
	if err != nil {
		return server.JourneyConfig{}, err
	}

	globalConfig := &server.GlobalConfiguration{}
	err = json.Unmarshal(configContent, globalConfig)
	if err != nil {
		return server.JourneyConfig{}, err
	}

	return server.MakeJourneyConfig(globalConfig)
}

func exitError(err error, message string) {
	_, err2 := fmt.Fprint(os.Stderr, fmt.Sprintf("%s\n%v\n", message, err))
	if err2 != nil {
		panic(err2)
	}
	os.Exit(1)
}
