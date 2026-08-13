package main

import (
	"fmt"

	"github.com/OpenBankingUK/conformance-suite/pkg/client"
	suiteVersion "github.com/OpenBankingUK/conformance-suite/pkg/version"
	"github.com/spf13/cobra"
)

const cliVersion = suiteVersion.FullVersion

func versionCmd(service client.Service) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of FCS CLI and Server",
		RunE:  versionCmdRun(service),
	}
}

func versionCmdRun(service client.Service) func(_ *cobra.Command, _ []string) error {
	return func(_ *cobra.Command, _ []string) error {
		version, err := service.Version()
		if err != nil {
			return err
		}
		fmt.Printf("CLI version %s, Server version %s\n", cliVersion, version.Version)
		return nil
	}
}
