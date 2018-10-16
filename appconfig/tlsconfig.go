package appconfig

import "crypto/tls"

// sub package to handle basic tls configuration

// NewTLSConfig - configures and returns a tls config
// based on the transport certificates present in the app config
// so it assume that you've already loaded these
func (a *AppConfig) NewTLSConfig() (*tls.Config, error) {

	cert, err := tls.X509KeyPair([]byte(a.CertTransport), []byte(a.KeyTransport)) // create x509 key pair from certs read into app config
	if err != nil {
		return &tls.Config{}, err
	}

	// Select tls cipherspecs etc ...  TODO: Tighten this config for example tls1.2 only ...
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert}, // pull in our OB transport cert
		InsecureSkipVerify: true,                    // Skip cert verification
		MinVersion:         tls.VersionSSL30,        // TLS 1.0/SSL 3.0 very old ... needs to be tls 1.2
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256, // not available by default however used by OB
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_RC4_128_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		},
	}
	tlsConfig.BuildNameToCertificate()

	return tlsConfig, nil
}
