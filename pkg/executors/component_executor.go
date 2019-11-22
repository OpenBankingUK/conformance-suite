package executors

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	resty "gopkg.in/resty.v1"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/manifest"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/sirupsen/logrus"
)

// GetHeadlessConsent -
func GetHeadlessConsent(definition RunDefinition, ctx *model.Context, specRun *generation.SpecRun, permissions map[string][]manifest.RequiredTokens) ([]manifest.RequiredTokens, error) {
	logger := logrus.WithFields(logrus.Fields{
		"module": "GetHeadlessConsent",
	})

	allRequiredTokens := []manifest.RequiredTokens{}

	for specType := range permissions {
		logger.Tracef("Getting Headless Consent for api type: %s", specType)
		tests, err := getSpecForSpecType(specType, specRun)
		if err != nil {
			return nil, err
		}

		switch specType {
		case "accounts":
			requiredTokens, err := getAccountsHeadlessTokens(tests, ctx, definition, logger)
			if err != nil {
				return nil, err
			}
			allRequiredTokens = append(allRequiredTokens, requiredTokens...)
		case "payments":
			requiredTokens, err := getPaymentHeadlessTokens(tests, ctx, definition, permissions["payments"], logger)
			if err != nil {
				return nil, err
			}
			allRequiredTokens = append(allRequiredTokens, requiredTokens...)
		default:
			logger.Fatalf("Support for spec type (%s) not implemented yet", specType)
		}
	}

	return allRequiredTokens, nil
}

func getPaymentHeadlessTokens(paymentTests []model.TestCase, ctx *model.Context, definition RunDefinition, requiredTokens []manifest.RequiredTokens, logger *logrus.Entry) ([]manifest.RequiredTokens, error) {
	logger.Debug("getPaymentHeadlessTokens")

	executor := Executor{}
	err := executor.SetCertificates(definition.SigningCert, definition.TransportCert)
	if err != nil {
		return nil, err
	}

	logger.Debugf("we have %d required tokens", len(requiredTokens))

	requiredTokens, err = runPaymentConsents(requiredTokens, ctx, &executor)
	if err != nil {
		logger.Errorf("getPaymentConsents error: " + err.Error())
	}

	tokendata, err := CallPaymentHeadlessConsentUrls(&requiredTokens, ctx, logger)
	if err != nil {
		return nil, err
	}

	for k, v := range requiredTokens {
		for x, y := range tokendata {
			if v.Name == x {
				v.Token = y
				requiredTokens[k] = v
			}
		}
	}

	logger.Tracef("updated requiredTokens: %#v", requiredTokens)
	return requiredTokens, err

}

// CallPaymentHeadlessConsentUrls -
func CallPaymentHeadlessConsentUrls(rt *[]manifest.RequiredTokens, ctx *model.Context, logger *logrus.Entry) (map[string]string, error) {
	var matchingGroup []string
	exchangeCode := ""
	exhangeCodeRegex := "code=([^&]*)&"
	consentedTokens := map[string]string{}

	for _, tokendata := range *rt {
		endpoint := tokendata.ConsentURL
		var resp *resty.Response

		resp, err := resty.R().
			SetHeader("accept", "*/*").
			Get(endpoint)

		if err != nil {
			if resp != nil && resp.StatusCode() == http.StatusFound { // catch status code 302 redirects and pass back as good response
				header := resp.Header()
				logger.Debugf("redirection headers: %#v", header)
				location := header.Get("Location")
				if location != "" {
					r, err := regexp.Compile(exhangeCodeRegex)
					if err != nil {
						return nil, err
					}
					matchingGroup = r.FindStringSubmatch(location)
					if len(matchingGroup[0]) < 2 {
						return nil, fmt.Errorf("Header Regex Context Match Failed - regex (%s) failed to find anything on Header (%s) value (%s)", exhangeCodeRegex, "Location", location)
					}
					exchangeCode = matchingGroup[1]
					logger.Tracef("retrieved Exchange code: %s", exchangeCode)
				}

			} else {
				logger.WithFields(logrus.Fields{
					"endpoint": endpoint,
					"err":      err,
				}).Debug("Error Calling Payment ConsentURL to get code")
				return nil, err
			}
		}

		if len(exchangeCode) < 1 {
			return nil, fmt.Errorf("Exchange code is empty - cannot complete exchange")
		}

		params, err := ctx.GetStrings("basic_authentication", "token_endpoint", "redirect_url")
		if err != nil {
			logger.Errorf("Consent Failed to get %s from context", err.Error())
			return nil, err
		}

		resp, err = resty.R().
			SetHeader("content-type", "application/x-www-form-urlencoded").
			SetHeader("accept", "application/json").
			SetHeader("authorization", "Basic "+params["basic_authentication"]).
			SetFormData(map[string]string{
				"code":         exchangeCode,
				"redirect_uri": params["redirect_url"],
				"grant_type":   "authorization_code",
				"scope":        "payments",
			}).
			Post(params["token_endpoint"])

		if err != nil {
			logger.WithFields(logrus.Fields{
				"tokenEndpoint": params["token_endpoint"],
				"err":           err,
			}).Debug("Payment headless code exchange failed")
			return nil, err
		}

		token, err := getAccessTokenFromJSONResponse(string(resp.Body()), logger)
		if err != nil {
			return nil, err
		}
		// Store the token against the token name for returning
		consentedTokens[tokendata.Name] = token
	}

	logger.Tracef("ConsentedTokens: %#v", consentedTokens)
	return consentedTokens, nil
}

func getAccessTokenFromJSONResponse(body string, logger *logrus.Entry) (string, error) {
	token := gjson.Get(body, "access_token")
	accessToken := token.String()
	if len(accessToken) == 0 {
		logger.WithFields(logrus.Fields{
			"body": body,
		}).Error("Access Token not found in JSON response body")
		return "", errors.New("Access Token not found in JSON response body")
	}
	return accessToken, nil
}

func getAccountsHeadlessTokens(tests []model.TestCase, ctx *model.Context, definition RunDefinition, logger *logrus.Entry) ([]manifest.RequiredTokens, error) {
	logger.Debug("getAccountsHeadlessTokens")
	bodyDataStart := "{\"Data\": { \"Permissions\": ["
	//TODO: sort out consent transaction timestamps
	txnFrom, err := ctx.GetString("transactionFromDate")
	if err != nil {
		return nil, errors.Wrap(err, "`transaction from date` not in context")
	}
	txnTo, err := ctx.GetString("transactionToDate")
	if err != nil {
		return nil, errors.Wrap(err, "`transaction to date` not in context")
	}

	bodyDataEnd := fmt.Sprintf(`], "TransactionFromDateTime": "%s", "TransactionToDateTime": "%s" },  "Risk": {} }`, txnFrom, txnTo)
	executor := NewExecutor()
	err = executor.SetCertificates(definition.SigningCert, definition.TransportCert)
	if err != nil {
		return nil, err
	}
	schemaVersion := definition.DiscoModel.DiscoveryModel.DiscoveryItems[0].APISpecification.SchemaVersion //TODO: Fix for more that one specification
	specType, err := manifest.GetSpecType(schemaVersion)
	if err != nil {
		return nil, errors.New("Error trying to determine specification type from API schemaVersion: " + err.Error())
	}

	requiredTokens, err := manifest.GetRequiredTokensFromTests(tests, specType)

	for k, tokenGatherer := range requiredTokens {

		localCtx := model.Context{}
		localCtx.PutContext(ctx)
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
