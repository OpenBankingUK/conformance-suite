package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

func runCmd(service Service) *cobra.Command {
	generatorCmd := &cobra.Command{
		Use:   "run",
		Short: "Run test cases from a discovery model",
		Long:  "Run test cases will output to standard output.",
		Run:   run(service),
	}
	generatorCmd.Flags().StringP("filename", "f", "", "Discovery filename")
	generatorCmd.Flags().StringP("config", "c", "", "Config filename")
	generatorCmd.Flags().StringP("output", "o", "", "Output filename, defaults to stdout")
	return generatorCmd
}

// run runs the functional conformance workflow to generate test case run report
func run(service Service) func(cmd *cobra.Command, _ []string) {
	return func(cmd *cobra.Command, _ []string) {
		// check if input (discovery model) filename if provided
		filenameFlag, err := cmd.Flags().GetString("filename")
		if err != nil || filenameFlag == "" {
			fmt.Println("You need to provide a discovery filename.")
			return
		}

		// check if input (discovery model) filename if provided
		configFlag, err := cmd.Flags().GetString("config")
		if err != nil || filenameFlag == "" {
			fmt.Println("You need to provide a discovery filename.")
			return
		}

		err = service.SetDiscoveryModel(filenameFlag)
		if err != nil {
			fmt.Printf("Error setting discovery model: %s\n", err.Error())
			return
		}

		err = service.SetConfig(configFlag)
		if err != nil {
			fmt.Printf("Error setting config: %s\n", err.Error())
			return
		}

		err = service.TestCases()
		if err != nil {
			fmt.Printf("Error generating test cases: %s\n", err.Error())
			return
		}

		resultsChan := make(chan TestCase)
		endedChan := make(chan struct{})

		err = service.RunTests(resultsChan, endedChan)
		if err != nil {
			fmt.Printf("Error generating test cases: %s\n", err.Error())
			return
		}

		resultPrinter(resultsChan, endedChan)
	}
}

// resultPrinter print incoming results from channel to console
func resultPrinter(resultChan chan TestCase, endedChan chan struct{}) {
	var passMsg = map[bool]string{true: "PASS", false: "FAIL"}
	const timeoutRunningTests = 5 * time.Minute

	deadline := time.NewTicker(timeoutRunningTests)
	defer deadline.Stop()
	for {
		select {
		case result := <-resultChan:
			fmt.Printf("=== %s: %s\n", passMsg[result.Pass], result.Id)
			if !result.Pass {
				fmt.Printf("\t %s\n", result.Fail)
			}

		case <-endedChan:
			return

		case <-deadline.C:
			fmt.Println("=== FAIL: Timeout waiting for tests for finish.")
			return
		}
	}
}
