package main

import (
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version/mocks"
	"github.com/stretchr/testify/mock"
)

func ExampleVersionCommand() {
	humanVersion := "0.1.2-RC1"
	warningMsg := "Version v0.1.2 of the Conformance Suite is out-of-date, please update to v0.1.3"
	formatted := "0.1.2"

	v := &mocks.Version{}
	v.On("GetHumanVersion").Return(humanVersion)
	v.On("UpdateWarningVersion", mock.AnythingOfType("string")).Return(warningMsg, true, nil)
	v.On("VersionFormatter", mock.AnythingOfType("string")).Return(formatted, nil)

	versionCommand := newVersionCommandWithOptions(v)
	err := versionCommand.run(nil, nil)
	if err != nil {
		return
	}
	// Output:
	// FCS - Functional Conformance Suite
	// Version v0.1.2 of the Conformance Suite is out-of-date, please update to v0.1.3
}
