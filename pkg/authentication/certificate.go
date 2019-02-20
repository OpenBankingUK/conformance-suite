package authentication

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Certificate - create new Certificate.
type Certificate interface {
	PublicKey() *rsa.PublicKey
	PrivateKey() *rsa.PrivateKey
	TLSCert() tls.Certificate
}

// certificate implements Certificate
type certificate struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	tlsCert    tls.Certificate
}

// NewCertificate - create new Certificate.
//
// Parameters:
// * publicKeyPem=PEM encoded public key.
// * privateKeyPem=PEM encoded private key.
//
// Returns Certificate, or nil with error set if something is invalid.
func NewCertificate(publicKeyPem, privateKeyPem string) (Certificate, error) {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyPem))
	if err != nil {
		return nil, errors.Wrap(err, "error with public key")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPem))
	if err != nil {
		return nil, errors.Wrap(err, "error with private key")
	}

	tlsCert, err := tls.X509KeyPair([]byte(publicKeyPem), []byte(privateKeyPem))
	if err != nil {
		logrus.Warnln("tls.X509KeyPair, err=", err)
	}

	if err := validateKeys(publicKey, privateKey); err != nil {
		return nil, err
	}

	return &certificate{
		publicKey:  publicKey,
		privateKey: privateKey,
		tlsCert:    tlsCert,
	}, nil
}

func (c certificate) PublicKey() *rsa.PublicKey {
	return c.publicKey
}

func (c certificate) PrivateKey() *rsa.PrivateKey {
	return c.privateKey
}

func (c certificate) TLSCert() tls.Certificate {
	return c.tlsCert
}

func validateKeys(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) error {
	// validate public and private key pair
	// see:
	// * https://stackoverflow.com/questions/20655702/signing-and-decoding-with-rsa-sha-in-go
	// * http://play.golang.org/p/bzpD7Pa9mr
	plaintext := []byte(`date: Thu, 05 Jan 2012 21:31:40 GMT`)

	hashed := sha256.Sum256(plaintext)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return errors.Wrap(err, "error signing")
	}

	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature); err != nil {
		return errors.Wrap(err, "error verifying")
	}

	return nil
}
