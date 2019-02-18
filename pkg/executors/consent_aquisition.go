package executors

import (
	"errors"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/sirupsen/logrus"
)

var (
	consentChannelTimeout = 30
)

// InitiationConsentAcquisition - get required tokens
func InitiationConsentAcquisition(consentRequirements []model.SpecConsentRequirements, definition RunDefinition, ctx *model.Context) (TokenConsentIDs, error) {
	consentIDChannel := make(chan TokenConsentIDItem, 100)
	tokenParameters := getConsentTokensAndPermissions(consentRequirements)

	for tokenName, permissionList := range tokenParameters {
		runner := NewConsentAcquisitionRunner(definition, NewBufferedDaemonController())
		tokenAcquisitionType := definition.DiscoModel.DiscoveryModel.TokenAcquisition
		permissionString := buildPermissionString(permissionList)
		consentInfo := TokenConsentIDItem{TokenName: tokenName, Permissions: permissionString}
		runner.RunConsentAcquisition(consentInfo, ctx, tokenAcquisitionType, consentIDChannel)
	}

	consentItems, err := waitForConsentIDs(consentIDChannel, tokenParameters)
	for _, v := range consentItems {
		logrus.Debugf("Setting Token: %s, ConsentId: %s", v.TokenName, v.ConsentID)
		ctx.PutString(v.TokenName, v.ConsentID)
	}
	return consentItems, err
}

func waitForConsentIDs(consentIDChannel chan TokenConsentIDItem, tokenParameters map[string][]string) (TokenConsentIDs, error) {
	consentItems := TokenConsentIDs{}
	consentIDsRequired := len(tokenParameters)
	consentIDsReceived := 0
	logrus.Debugf("waiting for consentids items ...")
	for {
		select {
		case item := <-consentIDChannel:
			logrus.Debugf("received consent channel item item %#v", item)
			consentIDsReceived++
			consentItems = append(consentItems, item)
			if consentIDsReceived == consentIDsRequired {
				logrus.Infof("Got %d required tokens - progressing..", consentIDsReceived)
				for _, v := range consentItems {
					logrus.Infof("token: %s, consentid: %s", v.TokenName, v.ConsentID)
				}
				return consentItems, nil
			}
		case <-time.After(time.Duration(consentChannelTimeout) * time.Second):
			logrus.Warnf("consent channel timeout after %d seconds", consentChannelTimeout)
			return consentItems, errors.New("ConsentChannel Timeout")
		}
	}
}

func getConsentTokensAndPermissions(consentRequirements []model.SpecConsentRequirements) map[string][]string {
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
		logrus.Debugf("Getting ConsentToken: %s: %s", k, buildPermissionString(v))
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
	logrus.Debugf("Looking to exchange code %s, token: %s", tokenName, code)
	r := NewExchangeComponentRunner(definition, NewBufferedDaemonController())

	r.runningLock.Lock()
	defer r.runningLock.Unlock()
	if r.running {
		return "", errors.New("exchange Code for Access token test cases runner already running")
	}
	r.running = true
	go r.RunExchangeCodeComponent(tokenName, code, scope, ctx)
	return "", nil
}

// RunExchangeCodeComponent -
func (r *TestCaseRunner) RunExchangeCodeComponent(tokenName, code, scope string, ctx *model.Context) (accesstoken string, err error) {
	r.executor.SetCertificates(r.definition.SigningCert, r.definition.TransportCert)
	ruleCtx := r.makeRuleCtx(ctx)

	basicAuth, err := ctx.GetString("basic_authentication")
	tokenEndpoint, err := ctx.GetString("token_endpoint")
	redirectURL, err := ctx.GetString("redirectURL")

	ruleCtx.PutString("exchange_code", code)
	ruleCtx.PutString("exchange_basic_auth", basicAuth)
	ruleCtx.PutString("exchange_token_endpoint", tokenEndpoint)
	ruleCtx.PutString("exchange_redirect_url", redirectURL)
	ruleCtx.PutString("exchange_scope", scope)
	ruleCtx.PutString("exchange_access_token", tokenName)

	var comp model.Component
	comp, err = model.LoadComponent("PSUConsentProviderComponent.json")
	if err != nil {
		r.AppErr("Load PSU Component Failed: " + err.Error())
		r.setNotRunning()
		return
	}

	for _, testcase := range comp.GetTests() {
		testResult := r.executeTest(testcase, ruleCtx, r.logger)
		r.daemonController.Results() <- testResult
		if testResult.Pass {
			accessToken, err := ruleCtx.GetString(tokenName)
			logrus.Debugf("received access token: %s for named %s", accessToken, tokenName)
			return accessToken, err
		}
	}
	r.running = false
	return "", nil
}
