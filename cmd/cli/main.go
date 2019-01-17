package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const banner = "FCS - Functional Conformance Suite"

func main() {
	mustReadViperEnvConfig()
	generatorCmdWrapper := newGeneratorCmdWrapper()
	rootCmd := rootCommand(generatorCmdWrapper.run)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type cobraCmdRunFunc func(cmd *cobra.Command, args []string)

func rootCommand(runFunc cobraCmdRunFunc) *cobra.Command {
	root := &cobra.Command{
		Use:   "fcs",
		Short: "Functional Conformance Suite CLI",
		Long:  `To use with pipelines and reproducible test runs`,
	}
	root.AddCommand(versionCmd())
	root.AddCommand(generatorCmd(runFunc))
	return root
}
