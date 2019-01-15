package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	mustReadViperEnvConfig()
	rootCmd := createRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const bitBucketRepository = "https://api.bitbucket.org/2.0/repositories/openbankingteam/conformance-suite/refs/tags"

func createRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "fcs",
		Short: "Functional Conformance Suite CLI",
		Long:  `To use with pipelines and reproducible test runs`,
	}

	versionCmdWrapper := newVersionCommand(bitBucketRepository)
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of FCS CLI",
		Run:   versionCmdWrapper.run,
	}
	root.AddCommand(versionCmd)

	return root
}
