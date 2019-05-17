# `CERTIFICATES`

## Certificates

### Production

<https://openbanking.atlassian.net/wiki/spaces/DZ/pages/80544075/OB+Root+and+Issuing+Certificates+for+Production>

### Sandbox

<https://openbanking.atlassian.net/wiki/spaces/DZ/pages/252018873/OB+Root+and+Issuing+Certificates+for+Sandbox>

### Convert Certificates

```sh
openssl x509 -inform der -in 'OB_SandBox_PP_Issuing CA.cer' -outform PEM -out OpenBanking_SandBox_IssuingCA.cer
openssl x509 -inform der -in 'OB_SandBox_PP_Root CA.cer' -outform PEM -out OpenBanking_SandBox_RootCA.cer
```

## Notes

If the certificate does not look valid, run this command:

```sh
mv OpenBankingRootCA.cer OpenBankingRootCA_original.cer
openssl x509 -in OpenBankingRootCA_original.cer -outform PEM -out OpenBankingRootCA.cer
```

## Testing

### With Certificates

**Note:** We do not pass in `-k` to `curl`.

```sh
$ cat OpenBankingIssuingCA.cer OpenBankingRootCA.cer > OpenBankingCA.cer
$ curl -v --cacert OpenBankingCA.cer https://private.api.hsbc.com/open-banking/v3.1/aisp/accounts
*   Trying 34.251.108.110...
* TCP_NODELAY set
* Connected to private.api.hsbc.com (34.251.108.110) port 443 (#0)
* ALPN, offering h2
* ALPN, offering http/1.1
* Cipher selection: ALL:!EXPORT:!EXPORT40:!EXPORT56:!aNULL:!LOW:!RC4:@STRENGTH
* successfully set certificate verify locations:
*   CAfile: OpenBankingCA.cer
  CApath: none
* TLSv1.2 (OUT), TLS handshake, Client hello (1):
* TLSv1.2 (IN), TLS handshake, Server hello (2):
* TLSv1.2 (IN), TLS handshake, Certificate (11):
* TLSv1.2 (IN), TLS handshake, Server key exchange (12):
* TLSv1.2 (IN), TLS handshake, Request CERT (13):
* TLSv1.2 (IN), TLS handshake, Server finished (14):
* TLSv1.2 (OUT), TLS handshake, Certificate (11):
* TLSv1.2 (OUT), TLS handshake, Client key exchange (16):
* TLSv1.2 (OUT), TLS change cipher, Client hello (1):
* TLSv1.2 (OUT), TLS handshake, Finished (20):
* TLSv1.2 (IN), TLS change cipher, Client hello (1):
* TLSv1.2 (IN), TLS handshake, Finished (20):
* SSL connection using TLSv1.2 / ECDHE-RSA-AES128-GCM-SHA256
* ALPN, server accepted to use h2
* Server certificate:
*  subject: C=GB; O=OpenBanking; OU=00158000016i44JAAQ; CN=6XhEzeZEYHnU6r2O7CZcvr
*  start date: Jul 11 10:35:19 2018 GMT
*  expire date: Jul 11 11:05:19 2020 GMT
*  subjectAltName: host "private.api.hsbc.com" matched cert's "private.api.hsbc.com"
*  issuer: C=GB; O=OpenBanking; CN=OpenBanking Issuing CA
*  SSL certificate verify ok.
* Using HTTP2, server supports multi-use
* Connection state changed (HTTP/2 confirmed)
* Copying HTTP/2 data in stream buffer to connection buffer after upgrade: len=0
* Using Stream ID: 1 (easy handle 0x7f85b180c600)
> GET /open-banking/v3.1/aisp/accounts HTTP/2
> Host: private.api.hsbc.com
> User-Agent: curl/7.54.0
> Accept: */*
>
* Connection state changed (MAX_CONCURRENT_STREAMS updated)!
< HTTP/2 401
< server: nginx
< date: Thu, 09 May 2019 13:43:34 GMT
< content-type: application/octet-stream
< content-length: 15
< strict-transport-security: max-age=31536000; includeSubDomains
<
* Connection #0 to host private.api.hsbc.com left intact
{"status": 401}
```

### Without Certificates

**Note:** We do not pass in `-k` to `curl`.

```sh
$ cat OpenBankingIssuingCA.cer OpenBankingRootCA.cer > OpenBankingCA.cer
$ curl -v https://private.api.hsbc.com/open-banking/v3.1/aisp/accounts
*   Trying 34.255.230.248...
* TCP_NODELAY set
* Connected to private.api.hsbc.com (34.255.230.248) port 443 (#0)
* ALPN, offering h2
* ALPN, offering http/1.1
* Cipher selection: ALL:!EXPORT:!EXPORT40:!EXPORT56:!aNULL:!LOW:!RC4:@STRENGTH
* successfully set certificate verify locations:
*   CAfile: /etc/ssl/cert.pem
  CApath: none
* TLSv1.2 (OUT), TLS handshake, Client hello (1):
* TLSv1.2 (IN), TLS handshake, Server hello (2):
* TLSv1.2 (IN), TLS handshake, Certificate (11):
* TLSv1.2 (OUT), TLS alert, Server hello (2):
* SSL certificate problem: unable to get local issuer certificate
* stopped the pause stream!
* Closing connection 0
curl: (60) SSL certificate problem: unable to get local issuer certificate
More details here: https://curl.haxx.se/docs/sslcerts.html

curl performs SSL certificate verification by default, using a "bundle"
 of Certificate Authority (CA) public keys (CA certs). If the default
 bundle file isn't adequate, you can specify an alternate file
 using the --cacert option.
If this HTTPS server uses a certificate signed by a CA represented in
 the bundle, the certificate verification probably failed due to a
 problem with the certificate (it might be expired, or the name might
 not match the domain name in the URL).
If you'd like to turn off curl's verification of the certificate, use
 the -k (or --insecure) option.
HTTPS-proxy has similar options --proxy-cacert and --proxy-insecure.
```
