package authentication

import (
	"log"
	"testing"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/stretchr/testify/assert"
)

var detachedJWT = `eyJ0eXAiOiJKT1NFIiwiY3R5IjoiYXBwbGljYXRpb24vanNvbiIsImh0dHA6Ly9vcGVuYmFua2luZy5vcmcudWsvaWF0IjoxNTg4NTg3NjgyLjQ1NiwiaHR0cDovL29wZW5iYW5raW5nLm9yZy51ay9pc3MiOiIwMDE1ODAwMDAxMDQxUkhBQVkiLCJodHRwOi8vb3BlbmJhbmtpbmcub3JnLnVrL3RhbiI6Im9wZW5iYW5raW5nLm9yZy51ayIsImNyaXQiOlsiaHR0cDovL29wZW5iYW5raW5nLm9yZy51ay9pYXQiLCJodHRwOi8vb3BlbmJhbmtpbmcub3JnLnVrL2lzcyIsImh0dHA6Ly9vcGVuYmFua2luZy5vcmcudWsvdGFuIl0sImFsZyI6IlBTMjU2Iiwia2lkIjoiREtlUE9MQU9pWEx3WWhNZkxTOGFTNllVLWQwIn0..1zMW5n7jXFGaOhVvL-Qz6ELVRzbfDzZahdXR3ioWA_H2MOib1Z346ZRaSczqjF2AY5qJfUX6AVpDopjCEDqmlCvSYsBOSFk0gwaNqnQVK4AN-yWK5OqC-gmo7W8RSTTF6s41yuXTdvZAPw7cdqmGKTHRvg2QpPkdHP8wXXurWqOgnUSgI6Czn_VKeIsc5W7rNpYF9onxY1HMDpXoYyXF_znYyWR3dNCueQaTHkIdt6b0MCBXINcgsY7pXsyHn-hZVGAW877sJjRC4GUfbZWKvkR2URLUOYKlzLYSGitsjtoHocESCG2uoovknTMLSIertSqbnm3VDVPRtBbJ0RSCuQ`

var rawBody = `{"Data":{"FundsAvailableResult":{"FundsAvailableDateTime":"2020-05-04T10:21:22.456Z","FundsAvailable":true}},"Links":{"Self":"http://ob19-rs1.o3bank.co.uk/open-banking/v3.1/pisp/domestic-payment-consents/sdp-1-241c9cc1-5dbc-46ca-a0df-9d512799c869/funds-confirmation"},"Meta":{}}`

func TestOzonePublicKey2EncodePayload(t *testing.T) {

	kid, err := getKidFromToken(detachedJWT)
	assert.Nil(t, err)
	jwk, err := getJwkFromJwks(kid, "https://keystore.openbankingtest.org.uk/0015800001041RHAAY/0015800001041RHAAY.jwks")
	assert.Nil(t, err)

	certs, err := ParseCertificateChain(jwk.X5c)
	assert.Nil(t, err)
	cert := certs[0]

	signedJWT, err := insertBodyIntoJWT(detachedJWT, rawBody, true)
	assert.Nil(t, err)

	verified, err := jws.Verify([]byte(signedJWT), jwa.PS256, cert.PublicKey)
	if err != nil {
		log.Printf("failed to verify message: %s", err)
		return
	}
	log.Printf("signed message verified! -> %s", verified)
}
