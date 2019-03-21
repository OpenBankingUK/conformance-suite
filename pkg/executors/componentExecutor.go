package executors

import (
	"fmt"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
)

// ExecuteComponent -
func ExecuteComponent(name string, ctx *model.Context) (model.Context, error) {
	var comp model.Component
	var err error
	if name == "headlessTokenProvider" {
		comp, err = model.LoadComponent("headlessTokenProviderComponent.json")
		if err != nil {
			return model.Context{}, fmt.Errorf("error loading headlessTokenProvider component:" + err.Error())
		}
	} else {
		return model.Context{}, fmt.Errorf("unknown component name:" + name)
	}

	err = comp.ValidateParameters(ctx)
	if err != nil {
		return model.Context{}, fmt.Errorf("error validating headlessTokenProvider component parameters:" + err.Error())
	}

	tests := comp.GetTests()
	// run sequentially - don't care about async ... its a startup task, not a run task.
	for k, test := range tests {
		_, _ = k, test
		//Add permissions/named tokens to context to have the right stuff result.
		//Execute the tests passing context between
		//Maybe need run defintion in here somewhere with certs and stuff ...
	}

	return model.Context{}, nil
}
