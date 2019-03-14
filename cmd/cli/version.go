package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

const cliVersion = "0.0.1"

func versionCmd(service Service) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of FCS CLI and Server",
		RunE:  versionCmdRun(service),
	}
}

func versionCmdRun(service Service) func(_ *cobra.Command, _ []string) error {
	return func(_ *cobra.Command, _ []string) error {
		version, err := service.Version()
		if err != nil {
			return err
		}
		fmt.Printf("CLI version %s, Server version %s\n", cliVersion, version.Version)
		return nil
	}
}
