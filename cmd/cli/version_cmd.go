package main

import (
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version"
	"fmt"
	"github.com/spf13/cobra"
)

const bitBucketRepository = "https://api.bitbucket.org/2.0/repositories/openbankingteam/conformance-suite/refs/tags"

// VersionCommand VersionCommand
type VersionCommand struct {
	versionChecker version.Checker
}

func versionCmd() *cobra.Command {
	versionCmdWrapper := newVersionCommand(bitBucketRepository)
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of FCS CLI",
		Run:   versionCmdWrapper.run,
	}
	return versionCmd
}

func newVersionCommand(bitBucketRepository string) VersionCommand {
	return VersionCommand{
		version.NewBitBucket(bitBucketRepository),
	}
}

func newVersionCommandWithOptions(versionChecker version.Checker) VersionCommand {
	return VersionCommand{
		versionChecker,
	}
}

func (v VersionCommand) run(_ *cobra.Command, _ []string) {
	fmt.Println(banner)
	softwareVersion := v.versionChecker.GetHumanVersion()
	uiMessage, _, _ := v.versionChecker.UpdateWarningVersion(softwareVersion)
	fmt.Println(uiMessage)
}
