package executors

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/sirupsen/logrus"
	resty "gopkg.in/resty.v0"
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
	logrus.Debugf("Looking to exchange code %s, tokenName: %s", code, tokenName)
	grantToken, err := exchangeCodeForToken(code, scope, ctx)
	if err != nil {
		logrus.Errorf("error attempting to exchange token %s", err.Error())
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

func exchangeCodeForToken(code, scope string, ctx *model.Context) (grantToken, error) {
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
		Post(tokenEndpoint + "/token")

	if err != nil {
		logrus.Debugf("error accessing exchange code url %s: %s ", tokenEndpoint, err.Error())
		return grantToken{}, err
	}
	if resp.StatusCode() != 200 {
		return grantToken{}, fmt.Errorf("bad status code %d from exchange token %s/token", resp.StatusCode(), tokenEndpoint)
	}
	var t grantToken
	err = json.Unmarshal(resp.Body(), &t)
	logrus.Debugf("exchangeCodeForToken  token: %#v", t)
	return t, err
}
