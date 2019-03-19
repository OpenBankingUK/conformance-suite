package report

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type reportClaims struct {
	jwt.StandardClaims
	ReportDigest    string `json:"reportDigest,omitempty"`
	DiscoveryDigest string `json:"discoveryDigest,omitempty"`
	ManifestDigest  string `json:"manifestDigest,omitempty"`
}

// Validate reportClaims
func (sc reportClaims) Valid() error {
	if err := sc.StandardClaims.Valid(); err != nil {
		return err
	}
	return nil
}

// required standardClaims fields:
// - Issuer
// - Subject
// - Id
// - NotBefore
// - ExpiresAt
// - ReportDigest
// - DiscoveryDigest
// - ManifestDigest
func sign(claims reportClaims, meta map[string]string, privateKey *rsa.PrivateKey) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)

	for k, v := range meta {
		t.Header[k] = v
	}

	signed, err := t.SignedString(privateKey)
	if err != nil {
		return "", errors.Wrap(err, "t.signedString()")
	}

	return signed, nil
}

func verify(rawJwt string, publicKey *rsa.PublicKey, claims reportClaims) error {
	keyFunc := func(*jwt.Token) (interface{}, error) {
		return publicKey, nil
	}

	parsedClaims := reportClaims{}
	t, err := jwt.ParseWithClaims(rawJwt, &parsedClaims, keyFunc)

	if err != nil {
		return errors.Wrap(err, "jwt.Parse()")
	}
	if !t.Valid {
		return errors.New("Token not invalid - unspecified")
	}
	if err := t.Claims.Valid(); err != nil {
		return errors.New("Token claims invalid - unspecified")
	}
	if parsedClaims.ReportDigest != claims.ReportDigest {
		return errors.New("report digest mismatch")
	}
	if parsedClaims.DiscoveryDigest != claims.DiscoveryDigest {
		return errors.New("discovery digest mismatch")
	}
	if parsedClaims.ManifestDigest != claims.ManifestDigest {
		return errors.New("manifest digest mismatch")
	}

	return nil
}
