package authentication

import (
	"encoding/base64"
	"errors"
	"fmt"
)

// CalculateClientSecretBasicToken tests the generation of `client secret basic` value as a product of
// `client_id` and `client_secret` as per https://tools.ietf.org/html/rfc7617
func CalculateClientSecretBasicToken(clientID, clientSecret string) (string, error) {
	if clientID == "" {
		return "", errors.New("clientID cannot be empty")
	}
	if clientSecret == "" {
		return "", errors.New("clientSecret cannot be empty")
	}
	subject := fmt.Sprintf("%s:%s", clientID, clientSecret)
	return base64.URLEncoding.EncodeToString([]byte(subject)), nil
}
