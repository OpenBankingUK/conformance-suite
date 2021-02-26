package authentication

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

var jwkCache = make(map[string]JWK)

var hsbcTanList = []string{
	"https://ob.hsbc.co.uk/jwks/public.jwks",
	"https://ob.firstdirect.com/jwks/public.jwks",
	"https://ob.mandsbank.com/jwks/public.jwks",
	"https://ob.business.hsbc.co.uk/jwks/public.jwks",
	"https://ob.hsbckinetic.co.uk/jwks/public.jwks",
	"https://ob.hsbcnet.com/jwks/public.jwks",
	"ob.hsbc.co.uk",
	"ob.firstdirect.com",
	"ob.mandsbank.com",
	"ob.business.hsbc.co.uk",
	"ob.hsbckinetic.co.uk",
	"ob.hsbcnet.com",
}

func isHSBCTrustAnchor(tan string) bool {
	for _, v := range hsbcTanList {
		if tan == v {
			return true
		}
	}
	return false
}

// getCertForKid
// Given a Kid - return the public cert from the JWKS keystore of the TrustAnchor
func getCertForKid(kid, jwks_uri string) (*x509.Certificate, error) {
	jwk, err := getJwkFromJwks(kid, jwks_uri)
	if err != nil {
		return nil, err
	}

	if len(jwk.X5c) == 0 {
		return nil, errors.New(fmt.Sprintf("No X5c certificate chain found for kid %s", kid))
	}

	certs, err := parseCertificateChain(jwk.X5c)
	if err != nil {
		return nil, err
	}

	cert := certs[0] // assumes a single certificate in chain which is the style used by the OB directory

	return cert, nil
}

// getJwkFromJwks
// Retieve the jwk representing a single public key from the jwks keystore
func getJwkFromJwks(kid, jwks string) (JWK, error) {
	if jwk, ok := jwkCache[kid]; !ok {
		logrus.Traceln("Retrieving JWKS url: " + jwks)
		jwks, err := getJwks(jwks)
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
	logrus.Traceln("no matching key found")
	return JWK{}, nil
}

// getJwks
// download the JWKS key store from the given url
func getJwks(url string) (JWKS, error) {
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

// parseCertificateChain
// takes a JWKS x5c claim containing a set of certs as strings - x509 cert objects in an array
func parseCertificateChain(chain []string) ([]*x509.Certificate, error) {
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
