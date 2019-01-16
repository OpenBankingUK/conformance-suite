package main

import (
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version"
	"fmt"
	"github.com/spf13/cobra"
)

const bitBucketRepository = "https://api.bitbucket.org/2.0/repositories/openbankingteam/conformance-suite/refs/tags"

func versionCmd() *cobra.Command {
	versionCmdWrapper := newVersionCommand(bitBucketRepository)
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of FCS CLI",
		Run:   versionCmdWrapper.run,
	}
	return versionCmd
}

type versionCommand struct {
	versionChecker version.Version
}

func newVersionCommand(bitBucketRepository string) versionCommand {
	return versionCommand{
		version.New(bitBucketRepository),
	}
}

func newVersionCommandWithOptions(versionChecker version.Version) versionCommand {
	return versionCommand{
		versionChecker,
	}
}

func (v versionCommand) run(_ *cobra.Command, _ []string) {
	fmt.Println(banner)
	softwareVersion := v.versionChecker.GetHumanVersion()
	uiMessage, _, _ := v.versionChecker.UpdateWarningVersion(softwareVersion)
	fmt.Println(uiMessage)
}
