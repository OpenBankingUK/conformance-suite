package executors

import (
	"fmt"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/manifest"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func getCbpiiConsents(
	definition RunDefinition,
	requiredTokens []manifest.RequiredTokens,
	ctx *model.Context,
) (TokenConsentIDs, error) {
	executor := &Executor{}
	err := executor.SetCertificates(definition.SigningCert, definition.TransportCert)
	if err != nil {
		logrus.Error(fmt.Sprintf("error running cbpii consent acquisition: %s", err))
		return nil, err
	}

	logrus.Debugf("we have %d cbpii consent required tokens", len(requiredTokens))
	for _, rt := range requiredTokens {
		logrus.Tracef("%#v", rt)
	}

	requiredTokens, err = runCbpiiConsents(requiredTokens, ctx, executor)
	if err != nil {
		logrus.Errorf("getCbpiiConsents error: %s", err)
	}

	consentItems := make([]TokenConsentIDItem, 0)
	for _, rt := range requiredTokens {
		tci := TokenConsentIDItem{TokenName: rt.Name, ConsentURL: rt.ConsentURL, ConsentID: rt.ConsentID}
		consentItems = append(consentItems, tci)
	}

	logrus.Debugf("we have %d consentIds: %#v", len(consentItems), consentItems)
	return consentItems, err
}

func runCbpiiConsents(rt []manifest.RequiredTokens, ctx *model.Context, executor *Executor) ([]manifest.RequiredTokens, error) {
	localCtx := model.Context{}
	localCtx.PutContext(ctx)
	localCtx.PutString("scope", "fundsconfirmations")
	consentJobs := manifest.GetConsentJobs()

	tc, err := readClientCredentialGrant()
	if err != nil {
		return nil, errors.New("cbpii PSU consent load clientCredentials testcase failed")
	}

	// Check for MTLS vs client basic authentication
	authMethod, err := ctx.GetString("token_endpoint_auth_method")
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
		return nil, errors.Wrap(err, "Cbpii PSU consent execute clientCredential grant testcase failed")
	}

	ccgBearerToken, err := localCtx.GetString("client_access_token")
	ctx.PutString("cbpii_ccg_token", ccgBearerToken)
	logrus.Debugf("runCbpiiConsents: just retrieved cbpii_ccg_token %v", ccgBearerToken)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot get Token for consent client credentials grant")
	}

	logrus.Tracef("runCbpiiConsents %d requiredTokens %#v", len(rt), rt)

	for k, v := range rt {
		localCtx.PutString("token_name", v.Name)

		test, exists := consentJobs.Get(v.ConsentProvider)
		if !exists {
			return nil, fmt.Errorf("Testcase %s does not exist in consentJob list", v.ConsentProvider)
		}
		test.InjectBearerToken(ccgBearerToken)
		test.Input.Headers["Content-Type"] = "application/json"

		err = executePaymentTest(&test, &localCtx, executor)
		if err != nil {
			return nil, errors.Wrap(err, "Cbpii PSU consent test case failed")
		}
		v.ConsentID, err = localCtx.GetString(v.ConsentParam)
		if err != nil {
			return nil, errors.Wrap(err, "Cbpii PSU consent test case failed - cannot find consentID in context")
		}
		localCtx.PutString("consent_id", v.ConsentID)
		localCtx.PutString("token_name", v.Name)

		exchange, err := readPsuExchange()
		if err != nil {
			return nil, errors.New("Cbpii PSU consent load psu_exchange testcase failed")
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
			return nil, errors.Wrap(err, "Cbpii PSU consent exchange code failed")
		}
		v.ConsentURL, err = localCtx.GetString("consent_url")
		if err != nil {
			return nil, errors.Wrap(err, "Cbpii PSU exchange test case failed - cannot find `consent_url` in context")
		}
		localCtx.Delete("consent_url")
		ctx.PutContext(&localCtx)
		rt[k] = v
	}

	logrus.Debug("Exit runCbpiiConsents Consents")
	logrus.Tracef("%#v", rt)
	return rt, nil
}
