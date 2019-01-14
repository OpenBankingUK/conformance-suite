package main

import (
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	rootCmd := createRootCommand()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "fcs",
		Short: "Functional Conformance Suite CLI",
		Long:  `To use with pipelines and reproducible test runs`,
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of FCS CLI",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("FCS CLI version: %s\n", version.Version)
		},
	}

	root.AddCommand(versionCmd)

	return root
}
