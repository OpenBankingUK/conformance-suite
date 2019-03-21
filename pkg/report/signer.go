package report

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

/*
	A report signature is manifested as a Json Web Token (JWT) saved in a simple plaintext file, e.g. signature.json
	There are a number of claims in the payload, some of which are standard and others custom.
	The custom claims added to the token are as follows:
	- reportDigest: SHA256 digest of the report file.
	- discoveryDigest: SHA256 digest of the discovery file.
	- manifestDigest: SHA256 digest of the manifest file.

	The intention is that during the signing process, the above digests are passed as claims. To verify the report
	the overall JWT is validated first, then each of the above digests is matched against the calculated digest of
	each of the files.

	Note: The digests are represented as hex strings.
*/

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

func verifySignature(rawJwt string, publicKey *rsa.PublicKey, claims reportClaims) error {
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

func verifyDigest(data []byte, expDigest string) error {
	digest, err := calculateDigest(data)
	if err != nil {
		return errors.Wrap(err, "calculate digest")
	}

	if digest != expDigest {
		return errors.New("digest mismatch")
	}
	return nil
}

func calculateDigest(data []byte) (string, error) {
	calc := sha256.Sum256(data)

	result := hex.EncodeToString(calc[0:])
	if result == "" {
		return "", errors.New("unknown hex encoding error")
	}

	return result, nil
}
