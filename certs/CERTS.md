# `CERTS.md`
Server certificates generated using [https://github.com/FiloSottile/mkcert](https://github.com/FiloSottile/mkcert).

## generate
### `mkcert`
Replace `/Users/mbana/work/openbankinguk/bitbucket/conformance-suite` with your local checkout path.

```sh
$ ( \
  go get -u github.com/FiloSottile/mkcert && \
  mkdir /Users/mbana/work/openbankinguk/bitbucket/conformance-suite/certs && \
  cd /Users/mbana/work/openbankinguk/bitbucket/conformance-suite/certs && \
  mkcert -cert-file="conformancesuite_cert.pem" -key-file="conformancesuite_key.pem" localhost 127.0.0.1 0.0.0.0 ::1 \
)
Using the local CA at "/Users/mbana/Library/Application Support/mkcert" ✨
Warning: the local CA is not installed in the system trust store! ⚠️
Warning: the local CA is not installed in the Firefox trust store! ⚠️
Run "mkcert -install" to avoid verification errors ‼️

Created a new certificate valid for the following names 📜
 - "localhost"
 - "127.0.0.1"
 - "0.0.0.0"
 - "::1"

The certificate is at "conformancesuite_cert.pem" and the key at "conformancesuite_key.pem" ✅

$ ls -1 /Users/mbana/work/openbankinguk/bitbucket/conformance-suite/certs
CERTS.md
conformancesuite_cert.pem
conformancesuite_key.pem
```

### alternative
[Step 1: Generate a self-signed X.509 TLS certificate](https://echo.labstack.com/cookbook/http2)
