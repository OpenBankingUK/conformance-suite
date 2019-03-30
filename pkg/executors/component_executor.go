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
func AcquireHeadlessTokens(tests []model.TestCase, ctx *model.Context, definition RunDefinition) ([]manifest.RequiredTokens, error) {
	logrus.Debug("=================================================================================================================")
	defer logrus.Debug("=================================================================================================================")
	logrus.Debug("AcquireHeadlessTokens")
	bodyDataStart := "{\"Data\": { \"Permissions\": ["
	//TODO: sort out consent transaction timestamps
	bodyDataEnd := "], \"TransactionFromDateTime\": \"2016-01-01T10:40:00+02:00\", \"TransactionToDateTime\": \"2025-12-31T10:40:00+02:00\" },  \"Risk\": {} }"

	executor := NewExecutor()
	err := executor.SetCertificates(definition.SigningCert, definition.TransportCert)
	if err != nil {
		return nil, err
	}

	requiredTokens, err := manifest.GetRequiredTokensFromTests(tests)
	logrus.Debugf("required tokens %#v\n", requiredTokens)

	for k, tokenGatherer := range requiredTokens {

		localCtx := model.Context{}
		localCtx.PutContext(ctx)
		localCtx.Put("SigningCert", definition.SigningCert) // For RS256 Claim signing
		permString := buildPermissionString(tokenGatherer.Perms)
		if len(permString) == 0 {
			continue
		}
		bodyData := bodyDataStart + permString + bodyDataEnd
		tokenName := tokenGatherer.Name
		localCtx.PutString("permission_payload", bodyData)
		localCtx.PutString("result_token", tokenName)

		returnCtx, err := executeComponent(&localCtx, executor)
		if err != nil {
			return nil, err
		}
		returnCtx.DumpContext("Return Context", tokenName, "client_access_token")
		clientGrantToken, _ := returnCtx.GetString("client_access_token")
		ctx.PutString("client_access_token", clientGrantToken)
		token, err := returnCtx.GetString(tokenName)
		if err != nil {
			return nil, err
		}
		tokenGatherer.Token = token
		requiredTokens[k] = tokenGatherer
	}

	return requiredTokens, nil
}

func getHeadlessTokenComponent() (*model.Component, error) {
	comp, err := model.LoadComponent("headlessTokenProviderComponent.json")
	if err != nil {
		return &comp, fmt.Errorf("error loading headlessTokenProvider component:" + err.Error())
	}
	return &comp, nil

}

// ExecuteComponent -
func executeComponent(ctx *model.Context, executor TestCaseExecutor) (*model.Context, error) {
	comp, err := getHeadlessTokenComponent()
	if err != nil {
		return nil, err
	}

	logrus.Debug("executeComponent - entry")
	err = comp.ValidateParameters(ctx)
	if err != nil {
		msg := fmt.Sprintf("error validating headlesstTokenProvider component %s", err.Error())
		logrus.Debug(msg)
		return &model.Context{}, fmt.Errorf(msg)
	}

	tests := comp.GetTests()
	executeCtx := &model.Context{}
	executeCtx.PutContext(ctx)
	logrus.Debugf("We have %d tests to run ", len(tests))
	// run sequentially - don't care about async ... its a startup task, not a run task.
	for k, test := range tests {
		test.ProcessReplacementFields(executeCtx, false)
		_, _ = k, test
		dumpJSON(test)
		logrus.Debug("Executing ------->>")

		req, err := test.Prepare(executeCtx)
		if err != nil {
			return &model.Context{}, err
		}
		resp, _, err := executor.ExecuteTestCase(req, &test, executeCtx)
		if err != nil {
			return &model.Context{}, fmt.Errorf("Test case %s failed with error %s", test.ID, err.Error())
		}

		result, errs := test.Validate(resp, executeCtx)
		if errs != nil {
			return &model.Context{}, fmt.Errorf("Test case %s Validation faiilure error %s", test.ID, errs[0].Error())
		}

		if !result {
			logrus.Errorf("Component testcase %s failed to Validate", test.ID)
			return &model.Context{}, errors.New("testcase failed to validate testid:" + test.ID)
		}

		logrus.Debug("Executed  <<-------")
		executeCtx.DumpContext("execution loop")

		//Add permissions/named tokens to context to have the right stuff result.
		//Execute the tests passing context between
		//Maybe need run defintion in here somewhere with certs and stuff ...
	}

	return executeCtx, nil
}

// Utility to Dump Json
func dumpJSON(i interface{}) {
	var model []byte
	model, _ = json.MarshalIndent(i, "", "    ")
	fmt.Println(string(model))
}
