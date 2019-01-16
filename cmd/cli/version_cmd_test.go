package main

import (
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version"
)

func ExampleversionCommand() {
	versionCommand := newVersionCommandWithOptions(version.Version{})
	versionCommand.run(nil, nil)
	// Output:
	// FCS - Functional Conformance Suite
	// Version check is unavailable at this time.
}
