package executors

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/manifest"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	resty "gopkg.in/resty.v1"
)

// errors
var (
	errTokenEndpointMethodUnsupported = errors.New("token_endpoint_auth_method unsupported")
)

var (
	consentChannelTimeout = 30
)

// GetPsuConsent -
func GetPsuConsent(definition RunDefinition, ctx *model.Context, runTests *generation.TestCasesRun, permissions map[string][]manifest.RequiredTokens) (TokenConsentIDs, map[string]string, error) {
	consentRequirements := runTests.SpecConsentRequirements
	var consentIdsToReturn TokenConsentIDs
	logrus.Debugf("running with %#v\n", permissions)

	for specType := range permissions {
		logrus.Tracef("Getting PSU Consent for api type: %s\n", specType)
		tests, err := getSpecForSpecType(specType, runTests)
		if err != nil {
			return nil, nil, err
		}

		switch specType {
		case "accounts":
			consentIds, _, err := getAccountConsents(consentRequirements, definition, permissions["accounts"], ctx)
			consentIdsToReturn = append(consentIdsToReturn, consentIds...)
			if err != nil {
				logrus.Error("GetPSUConsent - accounts error: " + err.Error())
				return nil, nil, err
			}

		case "payments": //TODO: Handle multipe spec returns
			consentIds, err := getPaymentConsents(tests, definition, permissions["payments"], ctx)
			consentIdsToReturn = append(consentIdsToReturn, consentIds...)
			if err != nil {
				logrus.Error("GetPSUConsent - payments error: " + err.Error())
				return nil, nil, err
			}
		default:
			logrus.Fatalf("Support for spec type (%s) not implemented yet", specType)
		}
	}

	logrus.Warnf("No Consent Acquistion Performed\n")
	return consentIdsToReturn, nil, nil
}

func getSpecForSpecType(stype string, runTests *generation.TestCasesRun) ([]model.TestCase, error) {
	for _, spec := range runTests.TestCases {
		specType, err := manifest.GetSpecType(spec.Specification.SchemaVersion)
		if err != nil {
			logrus.Warnf("cannot get spec type from SchemaVersion %s\n", spec.Specification.SchemaVersion)
			return nil, errors.New("Cannot get spec type from SchemaVersion " + spec.Specification.SchemaVersion)
		}
		if stype == specType {
			return spec.TestCases, nil
		}
	}

	return nil, errors.New("Cannot find test cases for spec type " + stype)
}

// getAccountConsents - get required tokens
func getAccountConsents(consentRequirements []model.SpecConsentRequirements, definition RunDefinition, permissions []manifest.RequiredTokens, ctx *model.Context) (TokenConsentIDs, map[string]string, error) {
	consentIDChannel := make(chan TokenConsentIDItem, 100)
	logger := logrus.StandardLogger().WithField("module", "getAccountConsents")
	logger.Tracef("getAccountConsents")

	tokenParameters := map[string]string{}

	//requiredTokens, err := manifest.GetRequiredTokensFromTests(spec.TestCases, "accounts")
	requiredTokens := permissions
	logrus.Tracef("we require %d tokens for `accounts`", len(requiredTokens))
	logrus.Tracef("required tokens %#v\n", requiredTokens)
	for _, rt := range requiredTokens {
		tokenParameters[rt.Name] = buildPermissionString(rt.Perms)
	}
	logrus.Tracef("required tokens %#v\n", tokenParameters)
	//for tokenName, permissionList := range tokenParameters {
	for _, rt := range requiredTokens {
		permissionList := rt.Perms
		tokenName := rt.Name
		runner := NewConsentAcquisitionRunner(logrus.StandardLogger().WithField("module", "InitiationConsentAcquisition"), definition, NewBufferedDaemonController())
		tokenAcquisitionType := definition.DiscoModel.DiscoveryModel.TokenAcquisition
		permissionString := buildPermissionString(permissionList)
		consentInfo := TokenConsentIDItem{TokenName: tokenName, Permissions: permissionString}
		err := runner.RunConsentAcquisition(consentInfo, ctx, tokenAcquisitionType, consentIDChannel)
		if err != nil {
			logger.WithError(err).Debug("InitiationConsentAcquisition")
		}
	}

	consentItems, err := waitForConsentIDs(consentIDChannel, len(tokenParameters))
	for _, v := range consentItems {
		logger.Debugf("Setting Token: %s, ConsentId: %s", v.TokenName, v.ConsentID)
		ctx.PutString(v.TokenName, v.ConsentID)
	}
	logrus.Debugf("we have %d consentIds: %#v", len(consentItems), consentItems)
	return consentItems, tokenParameters, err
}

func getTokenParametersFromRequiredTokens(tokens []manifest.RequiredTokens) map[string][]string {
	tokenParameters := map[string][]string{}
	for _, reqToken := range tokens {
		tokenParameters[reqToken.Name] = reqToken.Perms
	}
	return tokenParameters
}

func waitForConsentIDs(consentIDChannel chan TokenConsentIDItem, consentIDsRequired int) (TokenConsentIDs, error) {
	logger := logrus.StandardLogger().WithFields(logrus.Fields{
		"function":           "waitForConsentIDs",
		"consentIDsRequired": consentIDsRequired,
	})

	consentItems := TokenConsentIDs{}
	consentIDsReceived := 0
	logger.Debug("Waiting for items on consentIDChannel ...")
	for {
		select {
		case item := <-consentIDChannel:
			logrus.Debugf("received consent channel item item %#v", item)
			consentIDsReceived++
			consentItems = append(consentItems, item)

			logger.WithFields(logrus.Fields{
				"consentIDsReceived": consentIDsReceived,
				"consentIDsRequired": consentIDsRequired,
				"item":               fmt.Sprintf("%#v", item),
			}).Info("Progressing ...")

			errs := errors.New("")
			if consentIDsReceived == consentIDsRequired {
				logrus.Infof("Got %d required tokens - progressing..", consentIDsReceived)
				for _, v := range consentItems {
					if len(v.Error) > 0 {
						errs = errors.WithMessage(errs, v.Error)
					} else {
						logrus.Infof("token: %s, consentid: %s", v.TokenName, v.ConsentID)
					}
				}

				if len(errs.Error()) > 0 {
					return consentItems, errs
				}
				return consentItems, nil
			}
		case <-time.After(time.Duration(consentChannelTimeout) * time.Second):
			logrus.Warnf("consent channel timeout after %d seconds", consentChannelTimeout)
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

	for tokenParameterKey, tokenParameterValue := range tokenParameters {
		logger.WithFields(logrus.Fields{
			"tokenParameterValue":   tokenParameterValue,
			"tokenParameterKey":     tokenParameterKey,
			"buildPermissionString": buildPermissionString(tokenParameterValue),
		}).Debugf("Getting ConsentToken")
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
	logger := logrus.StandardLogger().WithFields(logrus.Fields{
		"module":    "ExchangeCodeForAccessToken",
		"tokenName": tokenName,
		"code":      code,
		"scope":     scope,
	})

	grantToken, err := exchangeCodeForToken(code, scope, ctx, logger)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("exchangeCodeForToken failed")
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

func exchangeCodeForToken(code string, scope string, ctx *model.Context, logger *logrus.Entry) (*grantToken, error) {
	logger = logger.WithFields(logrus.Fields{
		"function": "exchangeCodeForToken",
		"code":     code,
		"scope":    scope,
	})
	ctx.DumpContext()

	if scope == "" {
		scope = "accounts" // lets default to a scope of accounts
	}

	basicAuth, err := ctx.GetString("basic_authentication")
	if err != nil {
		return nil, errors.New("cannot get basic authentication for code")
	}
	tokenEndpoint, err := ctx.GetString("token_endpoint")
	if err != nil {
		return nil, errors.New("cannot get token_endpoint for code exchange")
	}
	redirectURI, err := ctx.GetString("redirect_url")
	if err != nil {
		return nil, errors.New("Cannot get redirect_url for code exchange")
	}
	clientID, err := ctx.GetString("client_id")
	if err != nil {
		return nil, errors.New("cannot get client_id for exchange code")
	}
	alg, err := ctx.GetString("requestObjectSigningAlg")
	if err != nil {
		return nil, errors.New("cannot get requestObjectSigningAlg for exchange code")
	}
	privKey, err := ctx.GetString("signingPrivate")
	if err != nil {
		return nil, errors.New("input, couldn't find `signingPrivate` in context")
	}
	pubKey, err := ctx.GetString("signingPublic")
	if err != nil {
		return nil, errors.New("input, couldn't find `signingPublic` in context")
	}
	cert, err := authentication.NewCertificate(pubKey, privKey)
	if err != nil {
		return nil, errors.Wrap(err, "input, couldn't create `certificate` from pub/priv keys")
	}

	// Check for MTLS vs client basic authentication
	authMethod, err := ctx.GetString("token_endpoint_auth_method")
	if err != nil {
		authMethod = authentication.ClientSecretBasic
	}

	var resp *resty.Response
	var errResponse error
	switch authMethod {
	case authentication.ClientSecretBasic:
		resp, err = resty.R().
			SetHeader("content-type", "application/x-www-form-urlencoded").
			SetHeader("accept", "*/*").
			SetHeader("authorization", "Basic "+basicAuth).
			SetFormData(map[string]string{
				authentication.GrantType: authentication.GrantTypeAuthorizationCode,
				"scope":                  scope, // accounts or payments currently
				"code":                   code,
				"redirect_uri":           redirectURI,
			}).
			Post(tokenEndpoint)
	case authentication.TlsClientAuth:
		resp, err = resty.R().
			SetHeader("content-type", "application/x-www-form-urlencoded").
			SetHeader("accept", "*/*").
			SetFormData(map[string]string{
				authentication.GrantType: authentication.GrantTypeAuthorizationCode,
				"scope":                  scope, // accounts or payments currently
				"code":                   code,
				"redirect_uri":           redirectURI,
				"client_id":              clientID,
			}).
			Post(tokenEndpoint)
	case authentication.PrivateKeyJwt:
		now := time.Now()
		iat := now.Unix()
		exp := now.Add(30 * time.Minute).Unix()
		jti := uuid.New().String()
		// https://openid.net/specs/openid-connect-core-1_0.html#ClientAuthentication
		// iss
		// REQUIRED. Issuer. This MUST contain the client_id of the OAuth Client.
		// sub
		// REQUIRED. Subject. This MUST contain the client_id of the OAuth Client.
		// aud
		// REQUIRED. Audience. The aud (audience) Claim. Value that identifies the Authorization Server as an intended audience. The Authorization Server MUST verify that it is an intended audience for the token. The Audience SHOULD be the URL of the Authorization Server's Token Endpoint.
		claims := jwt.MapClaims{
			"iss": clientID,
			"sub": clientID,
			"aud": tokenEndpoint,
			"iat": iat,
			"exp": exp,
			"jti": jti,
		}

		var signingMethod jwt.SigningMethod
		switch strings.ToUpper(alg) {
		case "PS256":
			// Workaround
			// https://github.com/dgrijalva/jwt-go/issues/285
			fixedSigningMethodPS256 := &jwt.SigningMethodRSAPSS{
				SigningMethodRSA: jwt.SigningMethodPS256.SigningMethodRSA,
				Options: &rsa.PSSOptions{
					SaltLength: rsa.PSSSaltLengthEqualsHash,
				},
			}
			signingMethod = fixedSigningMethodPS256
		case "RS256":
			signingMethod = jwt.SigningMethodRS256
		case "NONE":
			fallthrough
		default:
			return nil, errors.Errorf("unsupported algorithm: %q", alg)
		}

		token := jwt.NewWithClaims(signingMethod, claims) // create new token

		modulus := cert.PublicKey().N.Bytes()
		modulusBase64 := base64.RawURLEncoding.EncodeToString(modulus)
		kid, err := authentication.CalcKid(modulusBase64)
		if err != nil {
			return nil, errors.Wrap(err, "could not calculate kid")
		}
		token.Header["kid"] = kid

		clientAssertion, err := token.SignedString(cert.PrivateKey()) // sign the token - get as encoded string
		if err != nil {
			return nil, errors.Wrap(err, "could not generate client_assertion")
		}

		resp, errResponse = resty.R().
			SetHeader("content-type", "application/x-www-form-urlencoded").
			SetHeader("accept", "*/*").
			SetFormData(map[string]string{
				authentication.GrantType:           authentication.GrantTypeAuthorizationCode,
				"scope":                            scope,
				"code":                             code,
				"redirect_uri":                     redirectURI,
				authentication.ClientAssertionType: authentication.ClientAssertionTypeValue,
				authentication.ClientAssertion:     clientAssertion,
			}).
			Post(tokenEndpoint)
	default:
		logger.WithFields(logrus.Fields{
			"authMethod": authMethod,
		}).Error(errTokenEndpointMethodUnsupported)
		return nil, errTokenEndpointMethodUnsupported
	}

	logger.WithFields(logrus.Fields{
		"tokenEndpoint": tokenEndpoint,
	}).Debug("Attempting POST")

	if errResponse != nil {
		logger.WithFields(logrus.Fields{
			"tokenEndpoint": tokenEndpoint,
			"err":           err,
		}).Debug("Error accessing exchange code")
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("bad status code %d from exchange token %q", resp.StatusCode(), tokenEndpoint)
	}

	grantToken := &grantToken{}
	err = json.Unmarshal(resp.Body(), grantToken)
	logger.WithFields(logrus.Fields{
		"grantToken": grantToken,
	}).Debugf("OK")
	return grantToken, err
}
