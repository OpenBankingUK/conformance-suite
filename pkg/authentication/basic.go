package authentication

import (
	"encoding/base64"
	"errors"
	"fmt"
)

// calculateClientSecretBasicToken tests the generation of `client secret basic` value as a product of
// `client_id` and `client_secret` as per https://tools.ietf.org/html/rfc6749#section-4.4
func calculateClientSecretBasicToken(clientID, clientSecret string) (string, error) {
	if clientID == "" {
		return "", errors.New("clientID cannot be empty")
	}
	if clientSecret == "" {
		return "", errors.New("clientSecret cannot be empty")
	}
	subject := fmt.Sprintf("%s:%s", clientID, clientSecret)
	return base64.URLEncoding.EncodeToString([]byte(subject)), nil
}
