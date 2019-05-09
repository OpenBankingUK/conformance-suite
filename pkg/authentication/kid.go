package authentication

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

func CalcKid(modulus string) (string, error) {
	canonicalInput := fmt.Sprintf(`{"e":"AQAB","kty":"RSA","n":"%s"}`, modulus)

	sumer := sha1.New()
	_, err := io.WriteString(sumer, canonicalInput)
	if err != nil {
		return "", nil
	}
	sum := sumer.Sum(nil)

	sumBase64 := base64.RawURLEncoding.EncodeToString(sum)
	sumBase64NoTrailingEquals := strings.TrimSuffix(sumBase64, "=")

	return sumBase64NoTrailingEquals, nil
}
