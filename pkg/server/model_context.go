package server

import (
	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
)

const (
	CtxConstClientID                = "client_id"
	CtxConstClientSecret            = "client_secret"
	CtxConstTokenEndpoint           = "token_endpoint"
	CtxResponseType                 = "responseType"
	CtxConstTokenEndpointAuthMethod = "token_endpoint_auth_method"
	CtxConstFapiFinancialID         = "x-fapi-financial-id"
	CtxConstFapiCustomerIPAddress	= "x-fapi-customer-ip-address"
	CtxConstRedirectURL             = "redirect_url"
	CtxConstAuthorisationEndpoint   = "authorisation_endpoint"
	CtxConstBasicAuthentication     = "basic_authentication"
	CtxConstResourceBaseURL         = "resource_server"
	CtxConstIssuer                  = "issuer"
	CtxAPIVersion                   = "api-version"
	CtxConsentedAccountID           = "consentedAccountId"
	CtxStatementID                  = "statementId"
	CtxCreditorSchema               = "creditorScheme"
	CtxCreditorIdentification       = "creditorIdentification"
	CtxCreditorName                 = "creditorName"
	CtxInstructedAmountCurrency     = "instructedAmountCurrency"
	CtxInstructedAmountValue        = "instructedAmountValue"
	CtxCurrencyOfTransfer           = "currencyOfTransfer"
	CtxTransactionFromDate          = "transactionFromDate"
	CtxTransactionToDate            = "transactionToDate"
	CtxRequestObjectSigningAlg      = "requestObjectSigningAlg"
	CtxSigningPrivate               = "signingPrivate"
	CtxSigningPublic                = "signingPublic"
	CtxPhase                        = "phase"
	CtxNonOBDirectory               = "nonOBDirectory"
	CtxSigningKid                   = "signingKid"
	CtxSignatureTrustAnchor         = "signatureTrustAnchor"
)

func PutParametersToJourneyContext(config JourneyConfig, context model.Context) error {
	config.apiVersion = "v3.1"

	context.PutString(CtxConstClientID, config.clientID)
	context.PutString(CtxConstClientSecret, config.clientSecret)
	context.PutString(CtxConstTokenEndpoint, config.tokenEndpoint)
	context.PutString(CtxResponseType, config.ResponseType)
	context.PutString(CtxConstTokenEndpointAuthMethod, config.tokenEndpointAuthMethod)
	context.PutString(CtxConstFapiFinancialID, config.xXFAPIFinancialID)
	context.PutString(CtxConstFapiCustomerIPAddress, config.xXFAPICustomerIPAddress)
	context.PutString(CtxConstRedirectURL, config.redirectURL)
	context.PutString(CtxConstAuthorisationEndpoint, config.authorizationEndpoint)
	context.PutString(CtxConstResourceBaseURL, config.resourceBaseURL)
	context.PutString(CtxAPIVersion, config.apiVersion)
	context.PutString(CtxConsentedAccountID, config.resourceIDs.AccountIDs[0].AccountID)
	context.PutString(CtxStatementID, config.resourceIDs.StatementIDs[0].StatementID)
	context.PutString(CtxCreditorSchema, config.creditorAccount.SchemeName)
	context.PutString(CtxCreditorIdentification, config.creditorAccount.Identification)
	context.PutString(CtxCreditorName, config.creditorAccount.Name)
	context.PutString(CtxInstructedAmountCurrency, config.instructedAmount.Currency)
	context.PutString(CtxInstructedAmountValue, config.instructedAmount.Value)
	context.PutString(CtxCurrencyOfTransfer, config.currencyOfTransfer)
	context.PutString(CtxRequestObjectSigningAlg, config.requestObjectSigningAlgorithm)
	context.PutString(CtxSigningPrivate, config.signingPrivate)
	context.PutString(CtxSigningPublic, config.signingPublic)
	context.PutString(CtxTransactionFromDate, config.transactionFromDate)
	context.PutString(CtxTransactionToDate, config.transactionToDate)
	context.Put(CtxNonOBDirectory, config.useNonOBDirectory)
	context.PutString(CtxSigningKid, config.signingKid)
	context.PutString(CtxSignatureTrustAnchor, config.signatureTrustAnchor)

	basicauth, err := authentication.CalculateClientSecretBasicToken(config.clientID, config.clientSecret)
	if err != nil {
		return err
	}
	context.PutString(CtxConstBasicAuthentication, basicauth)
	context.PutString(CtxConstIssuer, config.issuer)
	context.PutString(CtxPhase, "unknown")

	return nil
}
