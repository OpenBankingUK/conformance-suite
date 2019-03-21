package executors

import (
	"encoding/json"
	"errors"
	"fmt"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/manifest"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/sirupsen/logrus"
)

// AcquireHeadlessTokens from manifest generated test cases
func AcquireHeadlessTokens(tests []model.TestCase, ctx *model.Context, definition RunDefinition) error {
	bodyDataStart := "{\"Data\": { \"Permissions\": ["
	//TODO: sort out consent transaction timestamps
	bodyDataEnd := "], \"TransactionFromDateTime\": \"2016-01-01T10:40:00+02:00\", \"TransactionToDateTime\": \"2025-12-31T10:40:00+02:00\" },  \"Risk\": {} }"
	component, err := getHeadlessTokenComponent()
	if err != nil {
		return err
	}

	executor := NewExecutor()
	err = executor.SetCertificates(definition.SigningCert, definition.TransportCert)
	if err != nil {
		return err
	}

	testcasePermissions, err := manifest.GetTestCasePermissions(tests)
	if err != nil {
		return err
	}
	requiredTokens, err := manifest.GetRequiredTokens(testcasePermissions)

	for _, tokenGatherer := range requiredTokens {
		localCtx := model.Context{}
		localCtx.PutContext(ctx)
		localCtx.Put("SigningCert", definition.SigningCert) // For RS256 Claim signing
		permString := buildPermissionString(tokenGatherer.Perms)
		bodyData := bodyDataStart + permString + bodyDataEnd
		tokenName := tokenGatherer.Name
		localCtx.PutString("permission_payload", bodyData)
		localCtx.PutString("result_token", tokenName)

		//TODO: Implement component call + error return
		returnCtx, err := executeComponent(component, &localCtx, executor)
		if err != nil {
			return err
		}
		fmt.Println("Return Context:")
		for k, v := range returnCtx {
			fmt.Printf("%s %s", k, v)
		}
	}

	return nil
}

func getHeadlessTokenComponent() (*model.Component, error) {
	comp, err := model.LoadComponent("headlessTokenProviderComponent.json")
	if err != nil {
		return &comp, fmt.Errorf("error loading headlessTokenProvider component:" + err.Error())
	}
	return &comp, nil

}

// ExecuteComponent -
func executeComponent(comp *model.Component, ctx *model.Context, executor TestCaseExecutor) (model.Context, error) {

	err := comp.ValidateParameters(ctx)
	if err != nil {
		return model.Context{}, fmt.Errorf("error validating headlessTokenProvider component parameters:" + err.Error())
	}

	tests := comp.GetTests()
	executeCtx := &model.Context{}
	executeCtx.PutContext(ctx)

	// run sequentially - don't care about async ... its a startup task, not a run task.
	for k, test := range tests {
		test.ProcessReplacementFields(executeCtx, false)
		_, _ = k, test
		dumpJSON(test)
		fmt.Println("Executing ------->>")

		req, err := test.Prepare(executeCtx)
		if err != nil {
			return model.Context{}, err
		}
		resp, _, err := executor.ExecuteTestCase(req, &test, executeCtx)
		if err != nil {
			return model.Context{}, fmt.Errorf("Test case %s failed with error %s", test.ID, err.Error())
		}

		result, err := test.Validate(resp, ctx)
		if err != nil {
			return model.Context{}, fmt.Errorf("Test case %s Validation faiilure error %s", test.ID, err.Error())
		}

		if !result {
			logrus.Errorf("Component testcase %s failed to Validate", test.ID)
			return model.Context{}, errors.New("testcase failed to validate testid:" + test.ID)
		}

		fmt.Println("Executed  <<-------")

		//Add permissions/named tokens to context to have the right stuff result.
		//Execute the tests passing context between
		//Maybe need run defintion in here somewhere with certs and stuff ...
	}

	return model.Context{}, nil
}

// Utility to Dump Json
func dumpJSON(i interface{}) {
	var model []byte
	model, _ = json.MarshalIndent(i, "", "    ")
	fmt.Println(string(model))
}
