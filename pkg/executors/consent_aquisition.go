package executors

import (
	"encoding/json"
	"fmt"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/manifest"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/sirupsen/logrus"
	"github.com/pkg/errors"
	resty "gopkg.in/resty.v1"
)

var consentChannelTimeout = 30

// InitiationConsentAcquisition - get required tokens
func InitiationConsentAcquisition(consentRequirements []model.SpecConsentRequirements, definition RunDefinition, ctx *model.Context, runTests *generation.TestCasesRun) (TokenConsentIDs, map[string]string, error) {
	tokenMap := make(map[string]string, 0)
	consentIDChannel := make(chan TokenConsentIDItem, 100)
	logger := logrus.StandardLogger().WithField("module", "InitiationConsentAcquisition")
	tokenParameters := getConsentTokensAndPermissions(consentRequirements, logger)

	tests := make([]model.TestCase, 0)
	for _, v := range runTests.TestCases {
		tests = append(tests, v.TestCases...)
	}

	requiredTokens, err := manifest.GetRequiredTokensFromTests(tests)
	//tokenParameters = getTokenParametersFromRequiredTokens(requiredTokens)
	_ = requiredTokens
	logrus.Debugf("required tokens %#v\n", requiredTokens)

	for tokenName, permissionList := range tokenParameters {
		runner := NewConsentAcquisitionRunner(logrus.StandardLogger().WithField("module", "InitiationConsentAcquisition"), definition, NewBufferedDaemonController())
		tokenAcquisitionType := definition.DiscoModel.DiscoveryModel.TokenAcquisition
		permissionString := buildPermissionString(permissionList)
		consentInfo := TokenConsentIDItem{TokenName: tokenName, Permissions: permissionString}
		errRun := runner.RunConsentAcquisition(consentInfo, ctx, tokenAcquisitionType, consentIDChannel)
		if errRun != nil {
			logger.WithError(errRun).Debug("InitiationConsentAcquisition")
		}
	}

	consentItems, err := waitForConsentIDs(consentIDChannel, tokenParameters, logger)
	for _, v := range consentItems {
		if len(v.Error) > 0 {
			continue
		}
		logger.Debugf("Setting Token: %s, ConsentId: %s", v.TokenName, v.ConsentID)
		ctx.PutString(v.TokenName, v.ConsentID)
	}

	return consentItems, tokenMap, err
}

func getTokenParametersFromRequiredTokens(tokens []manifest.RequiredTokens) map[string][]string {
	tokenParameters := make(map[string][]string, 0)
	for _, reqToken := range tokens {
		tokenParameters[reqToken.Name] = reqToken.Perms
	}
	return tokenParameters
}

func waitForConsentIDs(consentIDChannel chan TokenConsentIDItem, tokenParameters map[string][]string, logger *logrus.Entry) (TokenConsentIDs, error) {
	consentItems := TokenConsentIDs{}
	consentIDsRequired := len(tokenParameters)
	consentIDsReceived := 0
	logger.Debugf("waiting for consentids items ...")
	for {
		select {
		case item := <-consentIDChannel:
			logger.Debugf("received consent channel item item %#v", item)
			consentIDsReceived++
			consentItems = append(consentItems, item)
			errs := errors.New("")
			if consentIDsReceived == consentIDsRequired {
				logger.Infof("Got %d required tokens - progressing..", consentIDsReceived)
				for _, v := range consentItems {
					if len(v.Error) > 0 {
						errs = errors.WithMessage(errs, v.Error)
					} else {
						logger.Infof("token: %s, consentid: %s", v.TokenName, v.ConsentID)
					}
				}
				if len(errs.Error()) > 0 {
					return consentItems, errs
				}
				return consentItems, nil
			}
		case <-time.After(time.Duration(consentChannelTimeout) * time.Second):
			logger.Warnf("consent channel timeout after %d seconds", consentChannelTimeout)
			return consentItems, errors.New("ConsentChannel Timeout")
		}
	}
}

func getConsentTokensAndPermissions(consentRequirements []model.SpecConsentRequirements, logger *logrus.Entry) map[string][]string {
	tokenParameters := make(map[string][]string)
	for _, v := range consentRequirements {
		for _, namedPermission := range v.NamedPermissions {
			codeset := namedPermission.CodeSet
			for _, b := range codeset.CodeSet {
				mystring := string(b)
				set := tokenParameters[namedPermission.Name]
				set = append(set, mystring)
				tokenParameters[namedPermission.Name] = set
			}
		}
	}
	for k, v := range tokenParameters {
		logger.Debugf("Getting ConsentToken: %s: %s", k, buildPermissionString(v))
	}

	return tokenParameters
}

func buildPermissionString(permissionSlice []string) string {
	var permissions string
	first := true
	for _, perms := range permissionSlice {
		if !first {
			permissions += ","
		} else {
			first = !first
		}
		permissions += "\"" + perms + "\""
	}
	return permissions
}

// ExchangeParameters - Captures the parameters require to exchange a code for an access token
type ExchangeParameters struct {
	Code                string
	BasicAuthentication string
	TokenEndpoint       string
	RedirectURL         string
	Scope               string
	TokenName           string
}

// ExchangeCodeForAccessToken - runs a testcase to perform this operation
func ExchangeCodeForAccessToken(tokenName, code, scope string, definition RunDefinition, ctx *model.Context) (accesstoken string, err error) {
	logger := logrus.StandardLogger().WithField("module", "ExchangeCodeForAccessToken")
	logger.Debugf("Looking to exchange code %s, tokenName: %s", code, tokenName)
	grantToken, err := exchangeCodeForToken(code, scope, ctx, logger)
	if err != nil {
		logger.Errorf("error attempting to exchange token %s", err.Error())
	}
	return grantToken.AccessToken, err
}

type grantToken struct {
	AccessToken string `json:"access_token,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	Expires     int32  `json:"expires_in,omitempty"`
	Scope       string `json:"scope,omitempty"`
	IDToken     string `json:"id_token,omitempty"`
}

func exchangeCodeForToken(code, scope string, ctx *model.Context, logger *logrus.Entry) (grantToken, error) {
	basicAuth, err := ctx.GetString("basic_authentication")
	if err != nil {
		return grantToken{}, errors.New("cannot get basic authentication for code")
	}
	tokenEndpoint, err := ctx.GetString("token_endpoint")
	if err != nil {
		return grantToken{}, errors.New("cannot get token_endpoint for code exchange")
	}
	redirectURL, err := ctx.GetString("redirect_url")
	if err != nil {
		return grantToken{}, errors.New("Cannot get redirectURL for code exchange")
	}
	if scope == "" {
		scope = "accounts" // lets default to a scope of accounts
	}

	logger.Debugf("[%s] attempting POST %s", "exchangeCodeForToken", tokenEndpoint)
	resp, err := resty.R().
		SetHeader("content-type", "application/x-www-form-urlencoded").
		SetHeader("accept", "*/*").
		SetHeader("authorization", "Basic "+basicAuth).
		SetFormData(map[string]string{
			"grant_type":   "authorization_code",
			"scope":        scope, // accounts or payments currently
			"code":         code,
			"redirect_uri": redirectURL,
		}).
		Post(tokenEndpoint)

	if err != nil {
		logger.Debugf("error accessing exchange code url %s: %s ", tokenEndpoint, err.Error())
		return grantToken{}, err
	}
	if resp.StatusCode() != 200 {
		return grantToken{}, fmt.Errorf("bad status code %d from exchange token %s/token", resp.StatusCode(), tokenEndpoint)
	}
	var t grantToken
	err = json.Unmarshal(resp.Body(), &t)
	logger.Debugf("exchangeCodeForToken  token: %#v", t)
	return t, err
}
