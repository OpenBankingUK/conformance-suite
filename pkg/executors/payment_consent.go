package executors

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"github.com/pkg/errors"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/manifest"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/sirupsen/logrus"
)

func getPaymentConsents(definition RunDefinition, requiredTokens []manifest.RequiredTokens, ctx *model.Context) (TokenConsentIDs, error) {
	executor := &Executor{}
	err := executor.SetCertificates(definition.SigningCert, definition.TransportCert)
	if err != nil {
		logrus.Error("error running payment consent acquisition async: " + err.Error())
		return nil, err
	}

	logrus.Debugf("we have %d required tokens", len(requiredTokens))
	for _, rt := range requiredTokens {
		logrus.Tracef("%#v", rt)
	}

	requiredTokens, err = runPaymentConsents(requiredTokens, ctx, executor)
	if err != nil {
		logrus.Errorf("getPaymentConsents error: " + err.Error())
	}

	consentItems := make([]TokenConsentIDItem, 0)
	for _, rt := range requiredTokens {
		tci := TokenConsentIDItem{TokenName: rt.Name, ConsentURL: rt.ConsentURL, ConsentID: rt.ConsentID}
		consentItems = append(consentItems, tci)
	}

	logrus.Debugf("we have %d consentIds: %#v", len(consentItems), consentItems)
	return consentItems, err
}

func runPaymentConsents(rt []manifest.RequiredTokens, ctx *model.Context, executor *Executor) ([]manifest.RequiredTokens, error) {
	localCtx := model.Context{}
	localCtx.PutContext(ctx)
	localCtx.PutString("scope", "payments")
	consentJobs := manifest.GetConsentJobs()

	tc, err := readClientCredentialGrant()
	if err != nil {
		return nil, errors.New("payment PSU consent load clientCredentials testcase failed")
	}

	// Check for MTLS vs client basic authentication
	authMethod, err := ctx.GetString("token_endpoint_auth_method")
	logrus.Tracef("runpaymentconsent auth %s", authMethod)

	if err != nil {
		authMethod = "client_secret_basic"
	}

	switch authMethod {
	case authentication.ClientSecretBasic:
		tc.Input.SetHeader("authorization", "Basic $basic_authentication")
	case authentication.PrivateKeyJwt:
		clientID, err := ctx.GetString("client_id")
		if err != nil {
			return nil, errors.Wrap(err, "cannot find client_id for private_key_jwt form field")
		}
		tokenEndpoint, err := ctx.GetString("token_endpoint")
		if err != nil {
			return nil, errors.Wrap(err, "cannot find token_endpoint for private_key_jwt form field")
		}
		if tc.Input.Claims == nil {
			tc.Input.Claims = map[string]string{}
		}
		tc.Input.Claims["iss"] = clientID
		tc.Input.Claims["sub"] = clientID
		tc.Input.Claims["aud"] = tokenEndpoint
		clientAssertion, err := tc.Input.GenerateRequestToken(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "cannot generate request token for private_key_jwt form field")
		}
		tc.Input.SetFormField(authentication.ClientAssertionType, authentication.ClientAssertionTypeValue)
		tc.Input.SetFormField(authentication.ClientAssertion, clientAssertion)
	case authentication.TlsClientAuth:
		clientid, err := ctx.GetString("client_id")
		if err != nil {
			logrus.Warn("cannot locate client_id for tls_client_auth form field")
		}
		tc.Input.SetFormField("client_id", clientid)
	}

	tc.ProcessReplacementFields(&localCtx, true)
	err = executePaymentTest(&tc, &localCtx, executor)
	if err != nil {
		return nil, errors.New("Payment PSU consent execute clientCredential grant testcase failed :" + err.Error())
	}

	bearerToken, err := localCtx.GetString("client_access_token")
	ctx.PutString("payment_ccg_token", bearerToken) // store payment ccg token for later
	if err != nil {
		return nil, errors.New("Cannot get Token for consent client credentials grant: " + err.Error())
	}

	logrus.Tracef("runPaymentConsents %d requiredTokens %#v", len(rt), rt)

	for k, v := range rt {
		localCtx.PutString("token_name", v.Name)

		test, exists := consentJobs.Get(v.ConsentProvider)
		if !exists {
			return nil, errors.New("Testcase " + v.ConsentProvider + " does not exist in consentJob list")
		}
		test.InjectBearerToken(bearerToken) //client credential grant token
		test.Input.Headers["Content-Type"] = "application/json"

		err = executePaymentTest(&test, &localCtx, executor)
		if err != nil {
			return nil, errors.New("Payment PSU consent test case failed " + err.Error())
		}
		v.ConsentID, err = localCtx.GetString(v.ConsentParam)
		if err != nil {
			return nil, errors.New("Payment PSU consent test case failed - cannot find consentID in context " + err.Error())
		}
		localCtx.PutString("consent_id", v.ConsentID)
		localCtx.PutString("token_name", v.Name)

		exchange, err := readPsuExchange()
		if err != nil {
			return nil, errors.New("Payment PSU consent load psu_exchange testcase failed")
		}
		if authMethod == "tls_client_auth" {
			clientid, err := ctx.GetString("client_id")
			if err != nil {
				logrus.Warn("cannot locate client_id for tls_client_auth form field")
			}
			exchange.Input.SetFormField("client_id", clientid)

		} else {
			exchange.Input.SetHeader("authorization", "Basic $basic_authentication")
		}

		localCtx.DumpContext("before exchange", "token_name", "consent_id")
		err = executePaymentTest(&exchange, &localCtx, executor)
		if err != nil {
			return nil, errors.New("Payment PSU consent exchange code failed " + err.Error())
		}
		v.ConsentURL, err = localCtx.GetString("consent_url")
		if err != nil {
			return nil, errors.New("Payment PSU exchange test case failed - cannot find `consent_url` in context " + err.Error())
		}
		localCtx.Delete("consent_url")
		ctx.PutContext(&localCtx)
		rt[k] = v
	}

	clientGrantToken, err := localCtx.GetString("client_access_token")
	if err == nil {
		logrus.Tracef("setting payments client credential grant token to %s", clientGrantToken)
		ctx.PutString("client_access_token", clientGrantToken)
	}
	logrus.Debug("Exit runPayment Consents")
	logrus.Tracef("%#v", rt)
	return rt, nil
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
	if !result {
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

func readPsuExchange() (model.TestCase, error) {
	sc, err := model.LoadTestCaseFromJSONFile("components/psu_exchange.json")
	if err != nil {
		sc, err = model.LoadTestCaseFromJSONFile("../../components/psu_exchange.json")
	}
	return sc, err
}

func (r *TestCaseRunner) executePaymentConsent(tc model.TestCase, ruleCtx *model.Context, log *logrus.Entry) (bool, []string) {
	testresult := r.executeTest(tc, ruleCtx, log)
	return testresult.Pass, testresult.Fail

}
