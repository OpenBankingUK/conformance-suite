package signature

import (
	"crypto/rsa"
	"encoding/json"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// struct to hold all parameters required for signature creation
type SignatureParameters struct {
	cert           authentication.Certificate
	issuer         string
	nonOBDirectory bool
	kid            string
	trustAnchor    string
	alg            jwt.SigningMethod
	apiVersion     string
	body           string
}

type SignatureBodyClaim struct {
	Body string
}

// Body needs to be valid JSON
// So run it through the standard JSON marshaller to check
func (s SignatureBodyClaim) Valid() error {
	_, err := json.Marshal(s.Body) // check valid
	return err
}

func (s SignatureBodyClaim) MarshalJSON() ([]byte, error) {
	data := json.RawMessage(s.Body)
	return data.MarshalJSON()
}

func CreateSignature(sig SignatureParameters) (jwt.Token, error) {

	return jwt.Token{}, nil
}

var SigningMethodPS256 = &jwt.SigningMethodRSAPSS{
	SigningMethodRSA: jwt.SigningMethodPS256.SigningMethodRSA,
	Options: &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthEqualsHash,
	},
}

func SigningCertFromContext(ctx model.Context) (authentication.Certificate, error) {
	privKey, err := ctx.GetString("signingPrivate")
	if err != nil {
		return nil, errors.New("authentication.SigningCertFromContext: couldn't find `SigningPrivate` in context")
	}
	pubKey, err := ctx.GetString("signingPublic")
	if err != nil {
		return nil, errors.New("authentication.SigningCertFromContext: couldn't find `SigningPublic` in context")
	}
	cert, err := authentication.NewCertificate(pubKey, privKey)
	if err != nil {
		return nil, errors.Wrap(err, "authentication.SigningCertFromContext: couldn't create `certificate` from pub/priv keys")
	}
	return cert, nil
}
