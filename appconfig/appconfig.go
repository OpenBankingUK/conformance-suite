// Package appconfig - Package to read certs and parameters via the config/config.json file
/*

Typically this files looks like the following or simplar:-

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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

// AccessToken - Generic Access token
type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// AppConfig - application config
// partly read from config.json
// captures AccessTokens, Signing and Transport certs
type AppConfig struct {
	SoftwareStatementID   string `json:"softwareStatementId"` // OB Directory software statementid
	KeyID                 string `json:"keyId"`               // Signing cert key id
	TargetHost            string `json:"targetHost"`          // Host to proxy against
	Verbose               bool   `json:"verbose"`             // vebose output
	Spec                  string `json:"specLocation"`        // Spec location
	Bind                  string `json:"bindAddress"`         // bind adderss
	CertTransport         []byte // Certificates ..
	CertSigning           []byte
	KeySigning            []byte
	KeyTransport          []byte
	ClientCredentialToken AccessToken // Access Tokens
	AccountRequestToken   AccessToken
	PaymentRequestToken   AccessToken
}

// LoadAppConfiguration - from file - typically config/config.json
func LoadAppConfiguration(dir string) (*AppConfig, error) {
	var config *AppConfig
	file := dir + "/config.json"
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
		return &AppConfig{}, err
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		return &AppConfig{}, err
	}
	err = config.readCerts(dir)
	return config, err
}

// PrintAppConfig - dumps application config to console
func (a *AppConfig) PrintAppConfig() {
	logrus.WithFields(logrus.Fields{
		"SoftwareStatementId": a.SoftwareStatementID,
		"KeyID":               a.KeyID,
		"TargetHost":          a.TargetHost,
		"Verbose":             a.Verbose,
		"SpecLocation":        a.Spec,
		"BindAddress":         a.Bind,
	}).Info("AppConfig")
}

// Read certificates and keys from specified path
func (a *AppConfig) readCerts(configdir string) error {
	certTransport, err := ioutil.ReadFile(configdir + "/certTransport.pem")
	if err != nil {
		return fmt.Errorf("cannot read transport certificate")
	}
	keyTransport, err := ioutil.ReadFile(configdir + "/privateKeyTransport.key")
	if err != nil {
		return fmt.Errorf("cannot read transport key")
	}
	certSigning, err := ioutil.ReadFile(configdir + "/certSigning.pem")
	if err != nil {
		return fmt.Errorf("cannot read signing certificate")
	}
	keySigning, err := ioutil.ReadFile(configdir + "/privateKeySigning.key")
	if err != nil {
		return fmt.Errorf("cannot read signing key")
	}
	a.CertTransport = certTransport
	a.KeyTransport = keyTransport
	a.CertSigning = certSigning
	a.KeySigning = keySigning
	return nil
}
