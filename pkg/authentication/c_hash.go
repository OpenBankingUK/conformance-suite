package authentication

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// CalculateCHash calculates the code hash (c_hash) value
// as described in section 3.3.2.11 (ID Token) https://openid.net/specs/openid-connect-core-1_0.html#HybridIDToken
// List of valid algorithms https://openid.net/specs/openid-financial-api-part-2.html#jws-algorithm-considerations
// At the time of writing, the list shows "PS256", "ES256"
// https://openbanking.atlassian.net/wiki/spaces/DZ/pages/83919096/Open+Banking+Security+Profile+-+Implementer+s+Draft+v1.1.2#OpenBankingSecurityProfile-Implementer'sDraftv1.1.2-Step2:FormtheJOSEHeader
func CalculateCHash(alg string, code string) (string, error) {
	var digest []byte

	switch alg {
	case "ES256", "PS256":
		d := sha256.Sum256([]byte(code))
		//left most 256 bits.. 256/8 = 32bytes
		// no need to validate length as sha256.Sum256 returns fixed length
		digest = d[0:32]
	default:
		return "", fmt.Errorf("authentication.CalculateCHash: %q algorithm not supported", alg)
	}

	left := digest[0 : len(digest)/2]
	return base64.RawURLEncoding.EncodeToString(left), nil
}
