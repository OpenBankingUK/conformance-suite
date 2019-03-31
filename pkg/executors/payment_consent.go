package executors

import (
	"errors"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/manifest"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/sirupsen/logrus"
)

func getPaymentConsents(spec generation.SpecificationTestCases, definition RunDefinition, ctx *model.Context) (TokenConsentIDs, error) {
	psuConsentIDChannel := make(chan TokenConsentIDItem, 100)
	executor := &Executor{}
	err := executor.SetCertificates(definition.SigningCert, definition.TransportCert)
	if err != nil {
		logrus.Error("error running payment consent acquisition async: " + err.Error())
		return nil, err
	}
	requiredTokens, err := manifest.GetRequiredTokensFromTests(spec.TestCases, "payments")
	if err != nil {
		return nil, err
	}
	tests := spec.TestCases
	for _, rt := range requiredTokens {
		logrus.Tracef("%#v\n", rt)
	}

	err = runPaymentConsents(tests, requiredTokens, ctx, psuConsentIDChannel, executor)
	if err != nil {
		logrus.Errorf("getPaymentConsents error: " + err.Error())
	}

	tokenParameters := make(map[string]string, 0)
	for _, rt := range requiredTokens {
		tokenParameters[rt.Name] = buildPermissionString(rt.Perms)
	}

	logrus.Trace("<<==========================")
	dumpJSON(tokenParameters)
	logrus.Trace("<<==========================")
	consentItems, err := waitForConsentIDs(psuConsentIDChannel, len(tokenParameters))
	for _, v := range consentItems {
		logrus.Debugf("Setting Token: %s, ConsentId: %s", v.TokenName, v.ConsentID)
		ctx.PutString(v.TokenName, v.ConsentID)
	}
	logrus.Debugf("we have %d consentIds: %#v\n", len(consentItems), consentItems)
	return consentItems, err
}

func runPaymentConsents(tcs []model.TestCase, rt []manifest.RequiredTokens, ctx *model.Context,
	consentIDChannel chan<- TokenConsentIDItem, executor *Executor) error {

	// Execute client credential grant
	// load
	// process replacement fields
	// execute ccg
	// check we have access token
	// the for each in required tokens
	// execute POST consent - store consentid
	// return

	// $scope $token_endpoint $basic_authentication
	// client_access_token
	localCtx := model.Context{}
	localCtx.PutContext(ctx)
	localCtx.PutString("scope", "payments")

	tc, err := readClientCredentialGrant()
	if err != nil {
		return errors.New("Payment PSU consent load clientCredentials testcase failed")
	}

	tc.ProcessReplacementFields(&localCtx, true)
	err = executePaymentTest(&tc, &localCtx, executor)
	if err != nil {
		return errors.New("Payment PSU consent execute clientCredential grant testcase failed :" + err.Error())
	}

	bearerToken, err := localCtx.GetString("client_access_token")
	if err != nil {
		return errors.New("Cannot get Token for consent client credentials grant: " + err.Error())
	}

	for k, v := range rt {
		localCtx.PutString("token_name", v.Name)
		logrus.Warnln("Loop through requesting consent authorisation")
		test, err := findTest(tcs, v.ConsentProvider)
		if err != nil {
			return err
		}
		test.InjectBearerToken(bearerToken)
		test.Input.Headers["Content-Type"] = "application/json"
		logrus.Tracef("%v\n", test)
		err = executePaymentTest(test, &localCtx, executor)
		if err != nil {
			return errors.New("Payment PSU consent test case failed " + err.Error())
		}
		v.ConsentID, err = localCtx.GetString(v.ConsentParam)
		if err != nil {
			return errors.New("Payment PSU consent test case failed - cannot find consentID in context " + err.Error())
		}
		rt[k] = v
	}

	clientGrantToken, err := localCtx.GetString("client_access_token")
	if err == nil {
		logrus.Tracef("setting client credential grant token to %s", clientGrantToken)
		ctx.PutString("client_access_token", clientGrantToken)
	}
	logrus.Debug("Exit runPayment Consents")
	logrus.Tracef("%#v\n", rt)
	return nil
}

func findTest(tcs []model.TestCase, testID string) (*model.TestCase, error) {
	for k, test := range tcs {
		if test.ID == testID {
			return &tcs[k], nil
		}
	}
	return nil, errors.New("Test " + testID + " not found in findTest")
}

func executePaymentTest(tc *model.TestCase, ctx *model.Context, executor *Executor) error {
	req, err := tc.Prepare(ctx)
	if err != nil {
		logrus.Errorf("preparing to execute test %s: %s", tc.ID, err.Error())
		return err
	}
	resp, _, err := executor.ExecuteTestCase(req, tc, ctx)
	if err != nil {
		return err
	}
	result, errs := tc.Validate(resp, ctx)
	if errs != nil {
		return err
	}
	if result == false {
		return errors.New("testcase validation failed:" + err.Error())
	}
	return nil
}

func readClientCredentialGrant() (model.TestCase, error) {
	sc, err := model.LoadTestCaseFromJSONFile("components/clientcredentialgrant.json")
	if err != nil {
		sc, err = model.LoadTestCaseFromJSONFile("../../components/clientcredentialgrant.json")
	}
	return sc, err
}

func (r *TestCaseRunner) executePaymentConsent(tc model.TestCase, ruleCtx *model.Context, log *logrus.Entry) (bool, []string) {
	testresult := r.executeTest(tc, ruleCtx, log)
	return testresult.Pass, testresult.Fail

}
