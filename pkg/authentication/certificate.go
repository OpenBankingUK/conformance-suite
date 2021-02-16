package authentication

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

// Certificate - create new Certificate.
type Certificate interface {
	PublicKey() *rsa.PublicKey
	PrivateKey() *rsa.PrivateKey
	TLSCert() tls.Certificate
	DN() (string, string, string, error)
	SignatureIssuer(bool) (string, error)
}

// certificate implements Certificate
type certificate struct {
	publicKey     *rsa.PublicKey
	privateKey    *rsa.PrivateKey
	tlsCert       tls.Certificate
	publicCertPem []byte
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
		return nil, fmt.Errorf("error with public key: %w", err)
	}
	publicPem := []byte(publicKeyPem)

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPem))
	if err != nil {
		return nil, fmt.Errorf("error with private key: %w", err)
	}

	tlsCert, err := tls.X509KeyPair([]byte(publicKeyPem), []byte(privateKeyPem))
	if err != nil {
		logrus.StandardLogger().Warnln("tls.X509KeyPair, err=", err)
	}

	if err := validateKeys(publicKey, privateKey); err != nil {
		return nil, err
	}

	return &certificate{
		publicKey:     publicKey,
		privateKey:    privateKey,
		tlsCert:       tlsCert,
		publicCertPem: publicPem,
	}, nil
}

// creates a certificate from only the public key, in the case of the aspsp public cert to validate signatures
func NewPublicCertificate(publicKeyPem string) (Certificate, error) {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyPem))
	if err != nil {
		return nil, fmt.Errorf("error with public key: %w", err)
	}
	publicPem := []byte(publicKeyPem)

	return &certificate{
		publicKey:     publicKey,
		publicCertPem: publicPem,
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
		return fmt.Errorf("error signing: %w", err)
	}

	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature); err != nil {
		return fmt.Errorf("error verifying: %w", err)
	}

	return nil
}

func (c certificate) DN() (string, string, string, error) {
	co, o, ou, cn, err := c.nameComponents()
	if err != nil {
		return "", "", "", errors.New("error getting certificate DN " + err.Error())
	}
	dn := fmt.Sprintf("C=%s, O=%s, OU=%s, CN=%s", co, o, ou, cn)
	return dn, ou, cn, err
}

// SignatureIssuer - produces a string for use as the "http://openbanking.org.uk/iss" header field for
// x-jws-signature generation for the payments api v3.1
// for v3.0 the DN is used as this issuer
// addresses: https://openbanking.atlassian.net/browse/OBSD-8663
// bug fix: https://openbanking.atlassian.net/browse/REFAPP-784
// The handling of the issuer will likely change for in future iterations to handle thinks like EIDAS certs
// where the organisation id and software statement id are not present in the certificate DN
func (c certificate) SignatureIssuer(tpp bool) (string, error) {
	_, _, ou, cn, err := c.nameComponents()
	if err != nil {
		return "", errors.New("error getting certificate DN for SignatureIssuer: " + err.Error())
	}

	if tpp {
		if ou == "" {
			logrus.Warn("certificate ou is empty - if you're using EIDAS certificates you need to configure issuer")
		}
		return ou + "/" + cn, nil
	}

	return ou, nil

}

func (c certificate) nameComponents() (string, string, string, string, error) {
	cpb, _ := pem.Decode(c.publicCertPem)
	crt, err := x509.ParseCertificate(cpb.Bytes)
	if err != nil {
		logrus.Errorf("cannot parse cert %s", err.Error())
		return "", "", "", "", err
	}

	subject := crt.Subject

	var co string
	if len(subject.Country) > 0 {
		co = subject.Country[0]
	}

	var o string
	if len(subject.Organization) > 0 {
		o = subject.Organization[0]
	}

	var ou string
	if len(subject.OrganizationalUnit) > 0 {
		ou = subject.OrganizationalUnit[0]
	}

	cn := subject.CommonName

	return co, o, ou, cn, nil
}
