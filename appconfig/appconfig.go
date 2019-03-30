// Package appconfig - Package to read certs and parameters via the config/config.json file
/*

Typically this files looks like the following or simpler:-

    // config.json
	{
		"softwareStatementId": "1p8bGrHhJRrphjFk0qwNAU",
  		"clientScopes": "AuthoritiesReadAccess ASPSPReadAccess TPPReadAccess",
  		"keyId": "BFnipP2g4ZaaFySsIaigOUoCP2E",
  		"tokenUrl": "https://matls-sso.openbanking.me.uk/as/token.oauth2",
  		"tppTestUrl1":"https://matls-api.openbanking.me.uk/scim/v2/OBAccountPaymentServiceProviders",
  		"tppTestUrl":"https://tls-api.openbanking.me.uk/scim/v2/OBAccountPaymentServiceProviders",
  		"aud": "https://matls-sso.openbanking.me.uk/as/token.oauth2"
	}

	The following files will also be picked up from the config directory
	certSigning.pem - public signing cert
	certTransport.pem - public transport cert
	privateKeySigning.key - private signing key
	privateKeyTransport.key - private transport key

*/
package appconfig

import (
	"github.com/sirupsen/logrus"
)

// AccessToken - Generic Access token
type AccessToken struct {
	AccessToken string `json:"access_token" form:"access_token" query:"access_token"`
	ExpiresIn   int    `json:"expires_in" form:"expires_in" query:"expires_in"`
	TokenType   string `json:"token_type" form:"token_type" query:"token_type"`
}

// AppConfig - application config
// partly read from config.json
// captures AccessTokens, Signing and Transport certs
//
// To get these in single-line form use these commands:
// $ cat certTransport.pem | awk '{print}' ORS='\\n' | pbcopy
// $ cat certSigning.pem | awk '{print}' ORS='\\n' | pbcopy
// $ cat privateKeySigning.key | awk '{print}' ORS='\\n' | pbcopy
// $ cat privateKeyTransport.key | awk '{print}' ORS='\\n' | pbcopy
type AppConfig struct {
	SoftwareStatementID   string      `json:"softwareStatementId" form:"softwareStatementId" query:"softwareStatementId" validate:"required"` // OB Directory software statementid
	KeyID                 string      `json:"keyId" form:"keyId" query:"keyId" validate:"required"`                                           // Signing cert key id
	TargetHost            string      `json:"targetHost" form:"targetHost" query:"targetHost" validate:"required"`                            // Host to proxy against
	Verbose               bool        `json:"verbose" form:"verbose" query:"verbose" validate:"required"`                                     // verbose output
	Spec                  string      `json:"specLocation" form:"specLocation" query:"specLocation" validate:"required"`                      // Spec location
	Bind                  string      `json:"bindAddress" form:"bindAddress" query:"bindAddress" validate:"required"`                         // bind address
	CertTransport         string      `json:"certTransport" form:"certTransport" query:"certTransport" validate:"required"`
	CertSigning           string      `json:"certSigning" form:"certSigning" query:"certSigning" validate:"required"`
	KeySigning            string      `json:"keySigning" form:"keySigning" query:"keySigning" validate:"required"`
	KeyTransport          string      `json:"keyTransport" form:"keyTransport" query:"keyTransport" validate:"required"`
	ClientCredentialToken AccessToken `json:"client_credential_token" form:"client_credential_token" query:"client_credential_token"`
	AccountRequestToken   AccessToken `json:"account_request_token" form:"account_request_token" query:"account_request_token"`
	PaymentRequestToken   AccessToken `json:"payment_request_token" form:"payment_request_token" query:"payment_request_token"`
}

// PrintAppConfig - dumps application config to console
func (a *AppConfig) PrintAppConfig() {
	logrus.StandardLogger().WithFields(logrus.Fields{
		"SoftwareStatementId": a.SoftwareStatementID,
		"KeyID":               a.KeyID,
		"TargetHost":          a.TargetHost,
		"Verbose":             a.Verbose,
		"SpecLocation":        a.Spec,
		"BindAddress":         a.Bind,
	}).Info("AppConfig")
}
