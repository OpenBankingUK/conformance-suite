package authentication

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
)

const (
	publicCertValid = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCFENGw33yGihy92pDjZQhl0C3
6rPJj+CvfSC8+q28hxA161QFNUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6
Z4UMR7EOcpfdUE9Hf3m/hs+FUR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJw
oYi+1hqp1fIekaxsyQIDAQAB
-----END PUBLIC KEY-----`
	privateCertValid = `-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDCFENGw33yGihy92pDjZQhl0C36rPJj+CvfSC8+q28hxA161QF
NUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6Z4UMR7EOcpfdUE9Hf3m/hs+F
UR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJwoYi+1hqp1fIekaxsyQIDAQAB
AoGBAJR8ZkCUvx5kzv+utdl7T5MnordT1TvoXXJGXK7ZZ+UuvMNUCdN2QPc4sBiA
QWvLw1cSKt5DsKZ8UETpYPy8pPYnnDEz2dDYiaew9+xEpubyeW2oH4Zx71wqBtOK
kqwrXa/pzdpiucRRjk6vE6YY7EBBs/g7uanVpGibOVAEsqH1AkEA7DkjVH28WDUg
f1nqvfn2Kj6CT7nIcE3jGJsZZ7zlZmBmHFDONMLUrXR/Zm3pR5m0tCmBqa5RK95u
412jt1dPIwJBANJT3v8pnkth48bQo/fKel6uEYyboRtA5/uHuHkZ6FQF7OUkGogc
mSJluOdc5t6hI1VsLn0QZEjQZMEOWr+wKSMCQQCC4kXJEsHAve77oP6HtG/IiEn7
kpyUXRNvFsDE0czpJJBvL/aRFUJxuRK91jhjC68sA7NsKMGg5OXb5I5Jj36xAkEA
gIT7aFOYBFwGgQAQkWNKLvySgKbAZRTeLBacpHMuQdl1DfdntvAyqpAZ0lY0RKmW
G6aFKaqQfOXKCyWoUiVknQJAXrlgySFci/2ueKlIE1QqIiLSZ8V8OlpFLRnb1pzI
7U1yQXnTAEFYM560yJlzUpOb1V4cScGd365tiSMvxLOvTA==
-----END RSA PRIVATE KEY-----`
)

func TestCertificateValidateValidKeys(t *testing.T) {
	require := test.NewRequire(t)

	publicCert := publicCertValid
	privateCert := privateCertValid
	cert, err := NewCertificate(publicCert, privateCert)

	require.NotNil(cert)
	require.NoError(err)
}

func TestCertificateValidateValidKeysRandomlyGenerated(t *testing.T) {
	require := test.NewRequire(t)

	reader := rand.Reader
	bitSize := 512 // smallest no. of bits that is allowed by `rsa.GenerateKey`
	// bitSize := 2048
	// bitSize := 4096

	// Generate random keys
	privateKey, err := rsa.GenerateKey(reader, bitSize)
	require.NoError(err)
	require.NotNil(privateKey)

	require.NoError(privateKey.Validate())
	publicKey := &privateKey.PublicKey

	// Do not call MarshalPKCS1PublicKey
	// https://github.com/dgrijalva/jwt-go/blob/master/rsa_utils.go#L86
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	require.NoError(err)
	// Can call either MarshalPKCS1PrivateKey or MarshalPKCS8PrivateKey
	// https://github.com/dgrijalva/jwt-go/blob/master/rsa_utils.go#L27
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	// Get the keys in PEM format
	publicKeyPemData := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		},
	)
	privateKeyPemData := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		},
	)

	// Make new certificate
	publicCert := string(publicKeyPemData)
	privateCert := string(privateKeyPemData)
	cert, err := NewCertificate(publicCert, privateCert)

	require.NoError(err)
	require.NotNil(cert)
}

func TestCertificateValidatePrivateKeyInvalid(t *testing.T) {
	require := test.NewRequire(t)

	publicCert := publicCertValid
	privateCert := `-----BEGIN RSA PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCABIDFa/cp23OH
PZwpnBme3mvVun8ErtTpMtjCHBKbFyiVKI84e1sZt6BosIiXVbhJ3wsG+tmcJAVK
+rVRqSHUjPh9lHSI+4QQNvWIQC+v6zPyLKVRo47VK6dKTIZrsfA/A/+3hk/GrNx8
gbx9yENZ6WLg7GQ0mUwdV+kDwudll5sXjL7PHMIkVuxDOFv1cXXYmPuCtnFKvn/X
5XnD6fV8IOgLbwXTNzUxLdfnBn90PFYbRRUJVwIXzqSVJb3EXsEL6WAl2KdQ2BPS
UEv7jEw2Ja+fztUZUjqxWcrckq3rKUeHn3ykhTlq0+Iqg5sbgqz2zYt+6Zq4w42E
XrCqGlmVAgMBAAECggEAYLcGMiBnEpBgr4O0PxtXn9aZ0VacL4WGBMgNSli7FcBh
QI7r5NgM81jvLyhviSWRnP2M7zEExhnQhdzyr0b/7/ywnu9RO0wJcdaTmOQlItqm
3Acuvoa6mgHo2REHXMWJo5H510T5cDeYO9gn9z8c4wiXUyZEbhiCkIih2d2dw/mv
SBwfwOQN/3JXNBlW5c012usG6MLKvbOAYhqLfzq89ZqKJnrRW58Y6m+qb91fx5tt
DKrqVEXPBLcC0faXK+iMINGIbXqv8l1hd69f/SzSveI49yMgJBPCS5mtckcJqCm7
oxMafu3zGhaR7TYUM6CaqP3Sk2nLv6Sy9vNsFt4eAQKBgQD1TM/2nMV2HhVOl6Gm
hqCSGCxb6M+bv92tlkvuijpU0Bx5JjHTQ/dk8+0L0h/x+jPMo1+yrYyXNmE4kzyY
2s4jcDcBch5d8bSKksPv29sKYkMhEv0O31fcGVFhOyHvBm92EtA+l2vDVn0WSo7r
SAjsD3QSAdEr/pOqrUjInJfvMQKBgQCFmgjkVuMHgg5fjM+xQnPtLpJBRbe+gBdj
09jCj5D0mN5WpqvjO2cV40i2vYqaWz3/BgWcmlT4Crf2MtYxmR0rWsSQQsnk34EU
vej5kmkT8Pq4HRMskH2f/kNu55yHjY4TvMn13Gl9UWIb+g5oYdMbpb6B9jol9gm8
Op1wiopfpQKBgAbTvHYApv5CmBU34yffV1i5k4J7WEvday4JoNNixXzWzfQRPBHF
Mn18zHwnvPvfGtH3OhKfAeqzeME6V9VpQZN67Az+QBodQAkbTJjAZbhEQ9oHzUM8
tBVMHxe1rZwZccC3hVQ4oqctIQ4dxRyHRLhNNc3KfyfaTgHSENSEhzYBAoGAS4j8
MAVD1JHeiH03S9PzcQzcmdTN/wGyt7klm1LKNNBdHIadNhr2vHRFPzRIsd6WXaJM
9+51zctZmPPDEEWuLT3jVmC8fw8yjsSUfM4fZKvhRMkDdzW2IQgDnieK40TQKC6b
zMqyRa0GmCS3kqKEVeROonHRDHdfp7FIJEHf3BUCgYAjzIH744n31iTa/WhLXoIx
tR8rvTgdW6tn0YKGKy2RvgkABdDdjFc1oGBOSgNRopn1mL7cVtScMOenPTqKCEgO
yzVVOgE0ON1YmF/t3phOcJYIC9dy1CMHOjnG3deWxu2z3NAidF2TcAf8G/LtKOY2
t6uWCaHlDg2r/WDkbpxsgA==
-----END RSA PRIVATE KEY-----`
	cert, err := NewCertificate(publicCert, privateCert)

	require.Nil(cert)
	require.EqualError(err, `error verifying: crypto/rsa: verification error`)
}
func TestCertificateValidatePublicKeyInvalid(t *testing.T) {
	require := test.NewRequire(t)

	publicCert := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArBPlIPFnrLAOPY7Ltkxh
+Y1Jn7G0hjY+OVSb1I4Qi2zaypXMAeDKtcVoyCaEfGk1U6oMDVLj09iehUGtNBt8
jtV5IVoTu4DBlc9gKLdBt61mUb3l9ClZS6JWCiac3puInA+6VyBqtjcrwNITrHvw
O/qWRTmovG/Mw8kh49rYPDZ1jh3JMCAfbck3dvg/ULQVYMtjvTDDz9rKoPywz6Mn
68EgXy86jXfB0ekhhwuTnEALlO1bDQ0hrdzg8tum08OJJQH89VmbV2jB8lCZleS3
u11NvpdCkkjWK5+lyvPaVrg76pIYZ6gmpD9l6MbIK9XNmMHemsmRMi18HLuh+bTE
XwIDAQAB
-----END PUBLIC KEY-----`
	privateCert := privateCertValid
	cert, err := NewCertificate(publicCert, privateCert)

	require.Nil(cert)
	require.EqualError(err, `error verifying: crypto/rsa: verification error`)
}

func TestCertificateValidateInvalidPublicAndPrivateKeyEmpty(t *testing.T) {
	require := test.NewRequire(t)

	publicCert := ``
	privateCert := ``
	cert, err := NewCertificate(publicCert, privateCert)

	require.Nil(cert)
	require.EqualError(err, `error with public key: Invalid Key: Key must be PEM encoded PKCS1 or PKCS8 private key`)
}

func TestCertificateValidateInvalidPublicKeyEmpty(t *testing.T) {
	require := test.NewRequire(t)

	publicCert := ``
	privateCert := privateCertValid
	cert, err := NewCertificate(publicCert, privateCert)

	require.Nil(cert)
	require.EqualError(err, `error with public key: Invalid Key: Key must be PEM encoded PKCS1 or PKCS8 private key`)
}

func TestCertificateValidateInvalidPrivateKeyEmpty(t *testing.T) {
	require := test.NewRequire(t)

	publicCert := publicCertValid
	privateCert := ``
	cert, err := NewCertificate(publicCert, privateCert)

	require.Nil(cert)
	require.EqualError(err, `error with private key: Invalid Key: Key must be PEM encoded PKCS1 or PKCS8 private key`)
}

func TestCertificateValidateInvalidPublicKeyRSA(t *testing.T) {
	require := test.NewRequire(t)

	publicCert := `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE8xOUetsCa8EfOlDEBAfREhJqspDo
yEh6Szz2in47Tv5n52m9dLYyPCbqZkOB5nTSqtscpkQD/HpykCggvx09iQ==
-----END PUBLIC KEY-----`
	privateCert := privateCertValid
	cert, err := NewCertificate(publicCert, privateCert)

	require.Nil(cert)
	require.EqualError(err, `error with public key: Key is not a valid RSA public key`)
}

func TestCertificateValidateInvalidPrivateKeyRSA(t *testing.T) {
	require := test.NewRequire(t)

	publicCert := publicCertValid
	privateCert := `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQg6UIOwuA6Ww4E4ucb
xGhLoE5LKkxDFYz2X0dHwmV7q8ihRANCAARn6FZBTy1MQsxOGq4sDrZj5UELX1m2
VZ4E4JhiGX3fT07ucdN8nDXufQ8WplCX0nWFPh3P3Z7snLq2E3m8yqIQ
-----END PRIVATE KEY-----`
	cert, err := NewCertificate(publicCert, privateCert)

	require.Nil(cert)
	require.EqualError(err, `error with private key: Key is not a valid RSA private key`)
}
