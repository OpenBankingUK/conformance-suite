package authentication

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

func CalcKid(modulus string) (string, error) {
	canonicalInput := fmt.Sprintf(`{"e":"AQAB","kty":"RSA","n":"%s"}`, modulus)

	sumer := sha1.New()
	_, err := io.WriteString(sumer, canonicalInput)
	if err != nil {
		return "", fmt.Errorf("authentication.CalcKid: io.WriteString(sumer, canonicalInput) failed: %w", err)
	}
	sum := sumer.Sum(nil)

	sumBase64 := base64.RawURLEncoding.EncodeToString(sum)
	sumBase64NoTrailingEquals := strings.TrimSuffix(sumBase64, "=")

	return sumBase64NoTrailingEquals, nil
}

// GetKID determines the value of the JWS Key ID
func GetKID(ctx ContextInterface, modulus []byte) (string, error) {
	modulusBase64 := base64.RawURLEncoding.EncodeToString(modulus)
	kid, err := CalcKid(modulusBase64)
	if err != nil {
		return "", fmt.Errorf("authentication.GetKID: CalcKid(modulusBase64) failed: %w", err)
	}
	nonOBDirectory, exists := ctx.Get("nonOBDirectoryTPP")
	if !exists {
		return "", errors.New("authentication.GetKID: unable get nonOBDirectory value from context")
	}
	nonOBDirectoryAsBool, ok := nonOBDirectory.(bool)
	if !ok {
		return "", errors.New("authentication.GetKID: unable to cast nonOBDirectory value to bool")
	}
	if nonOBDirectoryAsBool {
		kid, err = ctx.GetString("signingKid")
		if err != nil {
			return "", fmt.Errorf("authentication.GetKID: unable to retrieve signingKid from context: %w", err)
		}
	}

	return kid, nil
}
