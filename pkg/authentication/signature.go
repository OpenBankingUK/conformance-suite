package authentication

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/sirupsen/logrus"
)

type JWKS struct {
	Keys []JWK
}

type JWK struct {
	Alg string   `json:"alg,omitempty"`
	Kty string   `json:"kty,omitempty"`
	X5c []string `json:"x5c,omitempty"`
	N   string   `json:"n,omitempty"`
	E   string   `json:"e,omitempty"`
	Kid string   `json:"kid,omitempty"`
	X5t string   `json:"x5t,omitempty"`
	X5u string   `json:"x5u,omitempty"`
	Use string   `json:"use,omitempty"`
}

// ValidateSignature
// take the signature JWT
// extract the kid
// used the kid to lookup the public key in the JWKS
//
func ValidateSignature(jwtToken, body, jwksUri string, b64 bool) (bool, error) {
	kid, err := getKidFromToken(jwtToken)
	if err != nil {
		return false, err
	}
	jwk, err := getJwkFromJwks(kid, jwksUri)
	if err != nil {
		return false, err
	}

	certs, err := ParseCertificateChain(jwk.X5c)
	if err != nil {
		return false, err
	}

	cert := certs[0]

	signature, err := insertBodyIntoJWT(jwtToken, body, b64) // b64claim
	if err != nil {
		logrus.Errorf("failed to insert body into signature message: %v", err)
		return false, err
	}
	logrus.Trace("Signature with payload: " + signature)

	verified, err := jws.Verify([]byte(signature), jwa.PS256, cert.PublicKey)
	if err != nil {
		logrus.Errorf("failed to verify message: %v", err)
		return false, err
	}

	logrus.Tracef("signed message verified! -> %s", verified)

	return true, nil
}

// insertBodyB64False
// jwt contains "header..signature"
// insert body into jwt resulting in "header.body.signature"
// b64 parameter controls body encoding
// b64=true  = R/W Api 3.1.4 and after
// b64=false = R/W Api 3.1.3 and prior
func insertBodyIntoJWT(token, body string, b64 bool) (string, error) {
	segments := strings.Split(token, ".")
	if len(segments) != 3 {
		return "", errors.New("Signature Token does not have 3 segments: " + token)
	}
	if b64 {
		segments[1] = base64.RawURLEncoding.EncodeToString([]byte(body))
	} else {
		segments[1] = body
	}
	return strings.Join(segments, "."), nil
}

func getKidFromToken(token string) (string, error) {
	var tokenHeader map[string]interface{}
	segments := strings.Split(token, ".")

	decodedPayload, _ := base64.RawURLEncoding.DecodeString(segments[0])

	json.Unmarshal(decodedPayload, &tokenHeader)

	kid, ok := tokenHeader["kid"].(string)
	if !ok {
		return "", fmt.Errorf("GetKidFromToken: error getting kid string from header")
	}
	if len(kid) == 0 {
		return "", fmt.Errorf("GetKidFromToken: error kid header is zero length")
	}

	return kid, nil
}

var jwkCache = make(map[string]JWK)

// Get JWK from JWKS_URI - cache responses
func getJwkFromJwks(kid, url string) (JWK, error) {
	if jwk, ok := jwkCache[kid]; !ok {
		logrus.Traceln("Retrieving JWKS url: " + url)
		jwks, err := GetJwks(url)
		if err != nil {
			return JWK{}, fmt.Errorf("GetJwkFromJwks: errors: %v", err)
		}
		for _, k := range jwks.Keys {
			if k.Kid == kid {
				jwkCache[kid] = k
				return k, nil
			}
		}
	} else {
		logrus.Traceln("Using cached jwk")
		return jwk, nil
	}
	return JWK{}, nil
}

func ParseCertificateChain(chain []string) ([]*x509.Certificate, error) {
	certchain := make([]*x509.Certificate, len(chain))
	for i, cert := range chain {
		raw, err := base64.StdEncoding.DecodeString(cert)
		if err != nil {
			return nil, errors.New("ParseCertificateChain: decode cert: " + err.Error())
		}
		certchain[i], err = x509.ParseCertificate(raw)
		if err != nil {
			return nil, errors.New("ParseCertificateChain: parse certificate: " + err.Error())
		}
	}
	return certchain, nil
}

func GetJwks(url string) (JWKS, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		return JWKS{}, fmt.Errorf("GetJwkss error retrieving url: %s, %v", url, err)
	}
	defer resp.Body.Close()
	jwksbytes := json.NewDecoder(resp.Body)
	var jwks JWKS
	if err := jwksbytes.Decode(&jwks); err != nil {
		return JWKS{}, fmt.Errorf("GetJwks: decoding error %s : %v", url, err)
	}

	return jwks, nil
}
