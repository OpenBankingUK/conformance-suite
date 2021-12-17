package main

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/client"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func runCmd(service client.Service) *cobra.Command {
	generatorCmd := &cobra.Command{
		Use:   "run",
		Short: "Run test cases from a discovery model",
		Long:  "Run test cases will output to standard output.",
		Run:   run(service),
	}
	generatorCmd.Flags().StringP("filename", "f", "", "Discovery filename")
	generatorCmd.Flags().StringP("config", "c", "", "Config filename")
	generatorCmd.Flags().StringP("export", "e", "", "Export config filename")
	return generatorCmd
}

// run runs the functional conformance workflow to generate test case run report
func run(service client.Service) func(cmd *cobra.Command, _ []string) {
	return func(cmd *cobra.Command, _ []string) {
		filenameFlag, err := cmd.Flags().GetString("filename")
		if err != nil || filenameFlag == "" {
			fmt.Println("You need to provide a discovery filename.")
			return
		}

		configFlag, err := cmd.Flags().GetString("config")
		if err != nil || filenameFlag == "" {
			fmt.Println("You need to provide a config filename.")
			return
		}

		exportFlag, err := cmd.Flags().GetString("export")
		if err != nil || filenameFlag == "" {
			fmt.Println("You need to provide a export config filename.")
			return
		}

		results, err := service.Run(filenameFlag, configFlag, exportFlag)
		if err != nil {
			fmt.Printf("Error running tests: %s\n", err.Error())
			return
		}

		client.ResultWriter(os.Stdout, results)
	}
}
