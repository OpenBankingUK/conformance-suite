package authentication

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"strings"
	"testing"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These tests exercise JWSVerify offline with locally generated keys, so they
// need neither a live JWKS endpoint nor the OB test private key.

func newRSAKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return key
}

func detachedSignature(t *testing.T, key *rsa.PrivateKey, body string, b64 bool) string {
	t.Helper()
	detached, err := buildSignature(b64, "test-kid", "test-issuer", "openbanking.org.uk", body, SigningMethodPS256, key)
	require.NoError(t, err)
	return detached
}

// withBody returns the full "header.body.signature" message in the form
// ValidateSignature hands to JWSVerify.
func withBody(t *testing.T, detached, body string, b64 bool) string {
	t.Helper()
	message, err := insertBodyIntoJWT(detached, body, b64)
	require.NoError(t, err)
	return message
}

func TestJWSVerifyAcceptsGenuineSignature(t *testing.T) {
	key := newRSAKey(t)
	for _, b64 := range []bool{true, false} {
		message := withBody(t, detachedSignature(t, key, rawBody, b64), rawBody, b64)

		payload, err := JWSVerify(message, jwa.PS256, &key.PublicKey, b64)

		require.NoError(t, err, "b64=%t", b64)
		assert.Equal(t, rawBody, string(payload), "b64=%t", b64)
	}
}

func TestJWSVerifyRejectsInvalidSignatures(t *testing.T) {
	key := newRSAKey(t)
	otherKey := newRSAKey(t)
	ecKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	for _, b64 := range []bool{true, false} {
		detached := detachedSignature(t, key, rawBody, b64)
		message := withBody(t, detached, rawBody, b64)
		tampered := withBody(t, detached, rawBody+" ", b64)
		undecodable := message[:strings.LastIndex(message, ".")+1] + "!!"

		tests := []struct {
			name    string
			message string
			alg     jwa.SignatureAlgorithm
			key     interface{}
			wantErr string
		}{
			{"tampered body", tampered, jwa.PS256, &key.PublicKey, "failed to verify message"},
			{"different RSA key", message, jwa.PS256, &otherKey.PublicKey, "failed to verify message"},
			{"RS256 verifier for a PS256 signature", message, jwa.RS256, &key.PublicKey, "failed to verify message"},
			{"non-RSA key", message, jwa.PS256, &ecKey.PublicKey, "failed to verify message"},
			{"undecodable signature", undecodable, jwa.PS256, &key.PublicKey, "failed to decode signature"},
		}
		for _, tt := range tests {
			t.Run(fmt.Sprintf("%s/b64=%t", tt.name, b64), func(t *testing.T) {
				payload, err := JWSVerify(tt.message, tt.alg, tt.key, b64)

				assert.EqualError(t, err, tt.wantErr)
				assert.Nil(t, payload)
			})
		}
	}
}
