package main

import (
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version"
	"fmt"
	"github.com/spf13/cobra"
)

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
	softwareVersion := v.versionChecker.GetHumanVersion()
	uiMessage, _, _ := v.versionChecker.UpdateWarningVersion(softwareVersion)
	fmt.Println(uiMessage)
}
