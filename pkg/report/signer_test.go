package report

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const exampleRawJWT = "eyJhbGciOiJQUzI1NiIsImhlYWRlci1mb28iOiJoZWFkZXItYmFyIiwidHlwIjoiSldUIn0.eyJleHAiOjQwNzc2MTQxMzIsImp0aSI6InVuaXF1ZS1qd3QtaWQiLCJpc3MiOiJodHRwczovL29wZW5iYW5raW5nLm9yZy51ay9mY3MvcmVwb3J0aW5nIiwibmJmIjo5MDAwMCwic3ViIjoib3BlbmJhbmtpbmcub3JnLnVrIiwicmVwb3J0RGlnZXN0IjoicmVwb3J0LWhhc2gtc3VtIiwiZGlzY292ZXJ5RGlnZXN0IjoiZGlzY292ZXJ5LWhhc2gtc3VtIiwibWFuaWZlc3REaWdlc3QiOiJtYW5pZmVzdC1oYXNoLXN1bSJ9.i0ry9tHVyDIkhIGVIYXmpEVGcLDRxKbio1uQdwQM0lj1h8nyPKwZnvnnB7Y0IkHc0nmELa2_nIVyfZxYAio1bk7Nj-M6bqQFv2Q-hE8deeJMwzLPOni4KtSf-a2tOXQQM29wQhAQ5fTIj3hlIsJQaRY5SnZUTejgLRaBVdtaWxo6bwOkzeqPOUEGlS67cleQZJvS6EcXA_tXVhyjgfUxr5oc9W7qFvbsmDYzipcTNdQ2q5ZPMg4K7KDh7QaoTnJi1U2JnVr12xMHoJFgmJDj79L-orTsJaBxWb0-CjTDr5i4CCBqkrSCiJAttkZNERZy11zGNL0x2ojL3ymvNH_xkg"

const samplePrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAsK8mIapI9HetPmwfptmHr2+oYW5YGKhzq4xEl6zAISChNVk9
DSMMLALfnlOaAK02yPqSiVSOYpyPjlEK/mRwETMSsRQvO/i+pO5aI9NeSfTo0HAc
NZ4nEwzsdrgrL3vrEagBxA2UgJM397CYhij/kJNK7Gec/jvlJAZxvr/k0SPV1iik
mYPCk0jHdbuMCEYMJ769pDvNwZr/JtV/uwG4H6BJUFtfKj3T0kjLQC4Y/uGTjKfm
l39iEuEop8Gi05WniEvbcpcI64XZpKasnuj7nYpjBW57oZzRW8viok9/+UEW0nE6
whKnyfhQ9TfJEicyoKFeLEOVbq9UdrhC2vwcaQIDAQABAoIBAQCQ/vwFDrEWZwx2
uNb033n5oGGHq72CZuOeOdukubFmvldt55Exsbxwdd88GJG+0meuYexV5V2AUcmB
2sJx6M0LYGWLiuwEhGs4AR9aXUD44pMZU5fi7KpWePmpqBRQwJo2ADGKyjY/mhGJ
JJTXLNgmtqn6/kEZZt/yQ5OfHe3TLv21nwhjfk9BV7MTEdA1Rr2BXyHD7h9MjgxK
9Q7vHtmhwwV67fQZFXwwN7kgrSSNFVNXDz2S+8IOcDondGTQFTau+S6x/0ziY6rg
ENjkzyoV9JO2WqnQyKb+rtpaFnDNdkuJHw6oVdeZOZwIL6CQ/RhmWDQPeDoc/NDK
aD6mzwaBAoGBAMDgCctycsilaxqZlHY5l6Qt1V0NH57+RNpzFzHidS8tEjbwHii7
FFQusTTcCFnFPpiQzKTTrjbPqEH5dcC99ip0BA9CGIzVWaWHFSVUPI0comhZhLVX
Y5+uX309f3LrSqNQmuRkwAOCeNlH5l9r0ncrDX9DkXp3x6uNtmITvOoRAoGBAOqC
jo/z7X9XJRhWwgW39AGIpvxrJgK9lIWsxfQma+NglyqvUA/Dzw9Ou870oJCVBZ39
CfZkLiJAUAk3F6lOdeEByzRy4A7NL94O8B3lQ6huOayksgfr8A4ScCWMyrTIgtCh
zATedKa69QGm/TAg6KJA8eP3K0snAPRt38cvKnTZAoGAJxkDQ0eG9x95L6I0Uybn
k3NrDfrMDynSAUpVSFp0kMSdLZ/NLUqHG21/pIx58OCoCLtJkJwMc7XykLUl5pVb
Yk20SPeIDHxvOLvCUJfb0mscjPSgjzYQztzFJJkjzcLelW6Qh33Y4p0/LCSEEZHE
zz1d9g9XXTEMu7z1XLpNkFECgYAB0kPDMHTOwWGDX+Ef5D7b6DDL0xU3fjtyElZz
P/0khfKGnVf012N7TfQ9dj7tAItLn9R8+mg1UeSNPcVMRlS6C6aFYMMGumc9xUXu
JYKyAzElex3628VAhroiQIaugsQpVKhd/VBQnzEZ8y8SOZ8062Y1jAzlB4eFXnkX
dfFReQKBgB8uf54s2HOI+Yx1YuQ5bYzOUp/J4PjyzPECjylsKe5J4pBntV1TvA9m
GRw/287a1mi9hDUlbOZSIVNJHAzxArCnJJnrW6C29NDFqWGIAgUS066KRZIyEuB3
SEtoekWeoBbByz3ehuKpuBK5St7Mz1MHyqb7YHQlTFj7oRQ9uHNK
-----END RSA PRIVATE KEY-----`

const samplePublicKey = `-----BEGIN RSA PUBLIC KEY-----
MIIBCgKCAQEAsK8mIapI9HetPmwfptmHr2+oYW5YGKhzq4xEl6zAISChNVk9DSMM
LALfnlOaAK02yPqSiVSOYpyPjlEK/mRwETMSsRQvO/i+pO5aI9NeSfTo0HAcNZ4n
EwzsdrgrL3vrEagBxA2UgJM397CYhij/kJNK7Gec/jvlJAZxvr/k0SPV1iikmYPC
k0jHdbuMCEYMJ769pDvNwZr/JtV/uwG4H6BJUFtfKj3T0kjLQC4Y/uGTjKfml39i
EuEop8Gi05WniEvbcpcI64XZpKasnuj7nYpjBW57oZzRW8viok9/+UEW0nE6whKn
yfhQ9TfJEicyoKFeLEOVbq9UdrhC2vwcaQIDAQAB
-----END RSA PUBLIC KEY-----`

func TestReportSigning(t *testing.T) {
	require := test.NewRequire(t)

	pb, rest := pem.Decode([]byte(samplePrivateKey))
	require.NotNil(pb, "pem decode private key, block not nil")
	require.Len(rest, 0, "rest should be zero length")
	privateKey, err := x509.ParsePKCS1PrivateKey(pb.Bytes)
	require.NoError(err, "parse private key")

	meta := map[string]string{
		"header-foo": "header-bar",
	}

	knownClaims := reportClaims{
		ReportDigest:    "report-hash-sum",
		DiscoveryDigest: "discovery-hash-sum",
		ManifestDigest:  "manifest-hash-sum",
		StandardClaims: jwt.StandardClaims{
			Issuer:    "https://openbanking.org.uk/fcs/reporting",
			Subject:   "openbanking.org.uk",
			Id:        "unique-jwt-id",
			ExpiresAt: 4077614132,
			NotBefore: 90000,
		},
	}

	signed, err := sign(knownClaims, meta, privateKey)
	require.NoError(err, "Sign error")
	calcClaims := reportClaims{}
	tk, err := jwt.ParseWithClaims(signed, &calcClaims, keyFunc)
	require.NoError(err, "jwt.Parse")

	signedString, err := tk.SigningString()
	require.NoError(err, "tk.SigningString()")
	expectedString := fmt.Sprintf("%s.%s", signedString, tk.Signature)

	require.Equal(signed, expectedString)
	require.Equal(knownClaims.Issuer, calcClaims.Issuer, "claim Issuer")
	require.Equal(knownClaims.Subject, calcClaims.Subject, "claim Subject")
	require.Equal(knownClaims.Id, calcClaims.Id, "claim Id")
	require.Equal(knownClaims.NotBefore, calcClaims.NotBefore, "claim NotBefore")
	require.Equal(knownClaims.ExpiresAt, calcClaims.ExpiresAt, "claim ExpiresAt")
	require.Equal(knownClaims.ReportDigest, calcClaims.ReportDigest, "claim ReportDigest")
	require.Equal(knownClaims.DiscoveryDigest, calcClaims.DiscoveryDigest, "claim DiscoveryDigest")
	require.Equal(knownClaims.ManifestDigest, calcClaims.ManifestDigest, "claim ManifestDigest")
}

func keyFunc(_ *jwt.Token) (interface{}, error) {
	pb, _ := pem.Decode([]byte(samplePublicKey))
	publicKey, err := x509.ParsePKCS1PublicKey(pb.Bytes)
	return publicKey, err
}

func TestReportVerificationSuccess(t *testing.T) {
	require := test.NewRequire(t)

	pb, rest := pem.Decode([]byte(samplePublicKey))
	require.NotNil(pb, "pem decode public key, block not nil")
	require.Len(rest, 0, "rest should be zero length")
	publicKey, err := x509.ParsePKCS1PublicKey(pb.Bytes)
	require.NoError(err, "parse public key")

	vps := reportClaims{
		ReportDigest:    "report-hash-sum",
		DiscoveryDigest: "discovery-hash-sum",
		ManifestDigest:  "manifest-hash-sum",
	}

	err = verifySignature(exampleRawJWT, publicKey, vps)
	require.NoError(err, "verify error")
}

func TestReportVerificationFailure(t *testing.T) {
	require := test.NewRequire(t)

	pb, rest := pem.Decode([]byte(samplePublicKey))
	require.NotNil(pb, "pem decode public key, block not nil")
	require.Len(rest, 0, "rest should be zero length")
	publicKey, err := x509.ParsePKCS1PublicKey(pb.Bytes)
	require.NoError(err, "parse public key")

	tt := []struct {
		Label           string
		ReportDigest    string
		DiscoveryDigest string
		ManifestDigest  string
		Error           error
	}{
		{
			Label:           "Report digest mismatch",
			Error:           errors.New("report digest mismatch"),
			ReportDigest:    "invalid-report-hash-sum",
			DiscoveryDigest: "discovery-hash-sum",
			ManifestDigest:  "manifest-hash-sum",
		},
		{
			Label:           "Discovery digest mismatch",
			Error:           errors.New("discovery digest mismatch"),
			ReportDigest:    "report-hash-sum",
			DiscoveryDigest: "invalid-discovery-hash-sum",
			ManifestDigest:  "manifest-hash-sum",
		},
		{
			Label:           "Manifest digest mismatch",
			Error:           errors.New("manifest digest mismatch"),
			ReportDigest:    "report-hash-sum",
			DiscoveryDigest: "discovery-hash-sum",
			ManifestDigest:  "invalid-manifest-hash-sum",
		},
	}

	for _, ti := range tt {
		vps := reportClaims{
			ReportDigest:    ti.ReportDigest,
			DiscoveryDigest: ti.DiscoveryDigest,
			ManifestDigest:  ti.ManifestDigest,
		}

		err = verifySignature(exampleRawJWT, publicKey, vps)
		require.EqualError(err, ti.Error.Error(), ti.Label)
	}
}

func TestDigestCalculation(t *testing.T) {
	input := "foo-bar"
	hexSHA256 := "7d89c4f517e3bd4b5e8e76687937005b602ea00c5cba3e25ef1fc6575a55103e"
	require := test.NewRequire(t)

	calcSHA256Output, err := calculateDigest([]byte(input))
	require.NoError(err, "calculateDigest")
	require.Equal(hexSHA256, calcSHA256Output)

}
func TestDigestVerification(t *testing.T) {
	input := "foo-bar"
	hexSHA256 := "7d89c4f517e3bd4b5e8e76687937005b602ea00c5cba3e25ef1fc6575a55103e"
	require := test.NewRequire(t)

	err := verifyDigest([]byte(input), hexSHA256)
	require.NoError(err)
}
