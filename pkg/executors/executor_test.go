package executors

import (
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"

	"github.com/stretchr/testify/require"
	resty "gopkg.in/resty.v1"
)

const (
	signingPublic = `-----BEGIN CERTIFICATE-----
MIIFLTCCBBWgAwIBAgIEWcVzXDANBgkqhkiG9w0BAQsFADBTMQswCQYDVQQGEwJH
QjEUMBIGA1UEChMLT3BlbkJhbmtpbmcxLjAsBgNVBAMTJU9wZW5CYW5raW5nIFBy
ZS1Qcm9kdWN0aW9uIElzc3VpbmcgQ0EwHhcNMTkwMzA0MTE1MzMxWhcNMjAwNDA0
MTIyMzMxWjBhMQswCQYDVQQGEwJHQjEUMBIGA1UEChMLT3BlbkJhbmtpbmcxGzAZ
BgNVBAsTEjAwMTU4MDAwMDEwNDFSZUFBSTEfMB0GA1UEAxMWcVkxZ1pEU2ZZTlRN
Y2U3ZGlmSXQ4UzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKGmfAiI
5EHmeDFtSHCpstndKvhXtUQl/SLCxU9uD5FTvwOyECb2nZ4mNdR6TlxnxoGgmwlj
TTSebTTMAF26pECxWf4SeQek6d0hyN2qe+TzjoaFrhhJhgsg1/plh+aXyVPrRC3A
uyi/dXy13ME+cPN1Otchm1vTTD7vFhT0GGQx4gTA2irYVGJCWtTMoZ6VYDyxa4Ls
0wp+2muks/ybX6PoZLUCkqLplKEck9MpfyS0NQc00spnfb3Tj+gAoYCUhwK/SQ8P
OSJn9TQ5xY7mRiwh1e8OUMdM+ZkT/k6dwm+5dgiFQa/IbkufNyd+JQluCn48msXT
EobpKjDp+tCRzmMCAwEAAaOCAfkwggH1MA4GA1UdDwEB/wQEAwIGwDAVBgNVHSUE
DjAMBgorBgEEAYI3CgMMMIHgBgNVHSAEgdgwgdUwgdIGCysGAQQBqHWBBgFkMIHC
MCoGCCsGAQUFBwIBFh5odHRwOi8vb2IudHJ1c3Rpcy5jb20vcG9saWNpZXMwgZMG
CCsGAQUFBwICMIGGDIGDVXNlIG9mIHRoaXMgQ2VydGlmaWNhdGUgY29uc3RpdHV0
ZXMgYWNjZXB0YW5jZSBvZiB0aGUgT3BlbkJhbmtpbmcgUm9vdCBDQSBDZXJ0aWZp
Y2F0aW9uIFBvbGljaWVzIGFuZCBDZXJ0aWZpY2F0ZSBQcmFjdGljZSBTdGF0ZW1l
bnQwbQYIKwYBBQUHAQEEYTBfMCYGCCsGAQUFBzABhhpodHRwOi8vb2IudHJ1c3Rp
cy5jb20vb2NzcDA1BggrBgEFBQcwAoYpaHR0cDovL29iLnRydXN0aXMuY29tL29i
X3BwX2lzc3VpbmdjYS5jcnQwOgYDVR0fBDMwMTAvoC2gK4YpaHR0cDovL29iLnRy
dXN0aXMuY29tL29iX3BwX2lzc3VpbmdjYS5jcmwwHwYDVR0jBBgwFoAUUHORxiFy
03f0/gASBoFceXluP1AwHQYDVR0OBBYEFGC6qjNxiOGjKsx1kqLIM9ZO0tWHMA0G
CSqGSIb3DQEBCwUAA4IBAQA2Kxag0hhJCtAP1oFICBnWLr6ThoJ6/htEYqEHGsRw
wGrI8Rh6KDDATtzwMGyZt+uMJvYApzkc2oBJFNYUGppa9DWEVPGKJAyeSbWNmFML
Peojk1kHqJkCnuzAzlZyV3KIea5RiWs26/ju3l/0QutHTzqb44fH5GC3c+x/atsE
UkhWuxXVUhw0T7lPOKyhhFeex+py5PNpQOVPj5xVDoUFEXGLliPX3az7w8hlAbFX
fx1oxsb22raWsERwS8DDux08kkH7dVCLPgvwJQz0vtlX6nAeoUVW7LO7PxnZEITZ
bt7hj9v+uKN8pBlvAym8wCY58zb/v/Di0k2r6ZRz0T5+
-----END CERTIFICATE-----
`
	signingPrivate = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQChpnwIiORB5ngx
bUhwqbLZ3Sr4V7VEJf0iwsVPbg+RU78DshAm9p2eJjXUek5cZ8aBoJsJY000nm00
zABduqRAsVn+EnkHpOndIcjdqnvk846Gha4YSYYLINf6ZYfml8lT60QtwLsov3V8
tdzBPnDzdTrXIZtb00w+7xYU9BhkMeIEwNoq2FRiQlrUzKGelWA8sWuC7NMKftpr
pLP8m1+j6GS1ApKi6ZShHJPTKX8ktDUHNNLKZ32904/oAKGAlIcCv0kPDzkiZ/U0
OcWO5kYsIdXvDlDHTPmZE/5OncJvuXYIhUGvyG5LnzcnfiUJbgp+PJrF0xKG6Sow
6frQkc5jAgMBAAECggEBAKGPHw/4oQksKpxbuLbBJDuSTEwAfO4reZ+wQjtsqKp6
pMIwyOvoNwfE8K/3vTGllkQgHFN5p8QbQtItwX/r9hWiK2s/Uy9Mp1+XUIYaydC9
i4jvOlyTvyCIJtPffb/9m/3/eRixM106XVXS/Vs16PWqCLDSqc9QkzejBNLUUzxu
93mImF62GUitMR35FRS2keXnRcX50uQZ4N6LmO/MV9HrPDNNdIAeGOs9cuEKD+/i
1jKl7OKTh6Jk/H6UkmnrcE1cW8yxLiCRNBgtMDUMg5qG/rnWYNCfoWOyK3oJxEC7
NE4yMK6uEvw1QPCdwLDnixMnMn9o+2PogXhs+EMeAiECgYEA1WUXn1w/IL3yzJ7B
4LWX3y2y9lRM5AA+KX4gKO3v/bzG0O96KF0MPHkjrXtl/FR4al8BeQ9lnR0sF8sa
RPFX3bYSBglXlyJy+ziMPPe0rsnfawlZ6LeYOSYHskJPXcl5TeIu63MP+Ipi0mRD
JGv7FDw9mjUBYR1hobjaHciE+xUCgYEAweyn/LqP6lc6PWuEEuq6GPfiCYI8GCIb
RrInuoXy4FamNE2R481hzb8ry2F778WbhgkQ2VhzmmjbAHiCLXof1wZOtJCiMCoz
liBBZP/zmgOWgXgL4bXPqjiYhJgMzSJyKiU3rwxW5l7KT5JhcLaOQh33qeIFoep1
xJeV2y+0IZcCgYEAtkAUoMIEGE6iIygjpWryPmWlRsRQtxmN/Zn+lXZBVY/4rVEa
H4b4gF1lnzCYtZzfCtoBRAdmXX0gv2FzGhaVWIG7evRXnniJgw2UmC1mXzGCYsQl
yZ+jnotgX1pKtmrv8xiNwgEPTtHB/LYssdqXIX0hj6ZdezfAvoJFptIu4NECgYAk
g58N40L96Pa6Yeg4d6Ia2XHiQHd4Q9PG9/yrDlWxEB+zcXeq4R0tVHW2keB4QUkL
b+GQSytZQ60Y5Zf9YCVmo3VmYmVnlEqqVeB6WAdSVKKeNjBmi4lSj92H+elPJtFA
Rkm52CT0s5x8Zx+ZzYXzxRjBECHnXvJV1gUNhGnyeQKBgA1Y64fM0LMMW9HPCILB
hUqchcOjocFfReVkrremr2Nz54/DSPMT79AwEzjYrElQyhtsPLmNJtT6FmzfQgMc
dG/g5RM0K9bni0KrTP2fw9R59oRxdnu667L3U+b8c7ogy9sH3N9bl4F0Gam+1B0/
29Z8cZKxCjXKlrjW/OP5Zyro
-----END PRIVATE KEY-----
`
)

const (
	transportPublic = `-----BEGIN CERTIFICATE-----
MIIFODCCBCCgAwIBAgIEWcVzXzANBgkqhkiG9w0BAQsFADBTMQswCQYDVQQGEwJH
QjEUMBIGA1UEChMLT3BlbkJhbmtpbmcxLjAsBgNVBAMTJU9wZW5CYW5raW5nIFBy
ZS1Qcm9kdWN0aW9uIElzc3VpbmcgQ0EwHhcNMTkwMzA0MTE1NDQ0WhcNMjAwNDA0
MTIyNDQ0WjBhMQswCQYDVQQGEwJHQjEUMBIGA1UEChMLT3BlbkJhbmtpbmcxGzAZ
BgNVBAsTEjAwMTU4MDAwMDEwNDFSZUFBSTEfMB0GA1UEAxMWcVkxZ1pEU2ZZTlRN
Y2U3ZGlmSXQ4UzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMeEjHvN
mRMiTeF6U4tr0oMQFQeQjaSFYEs4w6FEhXlIgfeZ7vZuNvlduYwVXCNFulrikJJT
NlU95abdBo1Gd+MkAe+6fuJ1fW1wr6IvB9eTGOG26EQYXdOk3U5c4PNopyCe8C71
itxAwU4frlTIqbZzavld9Eb6Fx10qm+Tl36E8zFIDYmB9ld7IFIem6n1opCkn5mq
OHke04weStqpC3b7rXh4rNW50B4c+KWB3v0a9KKyl+nLeh5Vcjm6KfRWt61nVVjd
Vmx8cunCsjzeOzPvwDrTUae5m3AIblOSa8wqrTpGHcp91gebvHoMhqWCgYBl3j9z
jdJwGNv6u7/g3CECAwEAAaOCAgQwggIAMA4GA1UdDwEB/wQEAwIHgDAgBgNVHSUB
Af8EFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwgeAGA1UdIASB2DCB1TCB0gYLKwYB
BAGodYEGAWQwgcIwKgYIKwYBBQUHAgEWHmh0dHA6Ly9vYi50cnVzdGlzLmNvbS9w
b2xpY2llczCBkwYIKwYBBQUHAgIwgYYMgYNVc2Ugb2YgdGhpcyBDZXJ0aWZpY2F0
ZSBjb25zdGl0dXRlcyBhY2NlcHRhbmNlIG9mIHRoZSBPcGVuQmFua2luZyBSb290
IENBIENlcnRpZmljYXRpb24gUG9saWNpZXMgYW5kIENlcnRpZmljYXRlIFByYWN0
aWNlIFN0YXRlbWVudDBtBggrBgEFBQcBAQRhMF8wJgYIKwYBBQUHMAGGGmh0dHA6
Ly9vYi50cnVzdGlzLmNvbS9vY3NwMDUGCCsGAQUFBzAChilodHRwOi8vb2IudHJ1
c3Rpcy5jb20vb2JfcHBfaXNzdWluZ2NhLmNydDA6BgNVHR8EMzAxMC+gLaArhilo
dHRwOi8vb2IudHJ1c3Rpcy5jb20vb2JfcHBfaXNzdWluZ2NhLmNybDAfBgNVHSME
GDAWgBRQc5HGIXLTd/T+ABIGgVx5eW4/UDAdBgNVHQ4EFgQU/Ip5tlBmO9C7cme9
zv0kNYAv00EwDQYJKoZIhvcNAQELBQADggEBADRRvzTSdp1IfrveCj/OiumKfP78
guQUYyffd48+P5CgIleaPSKNYV7NQ7tCP8BXrLw6kFa45uQgyE55kCXmWgsdsds8
yHDIqUHJAWixlPY8hIfEmVoYLK6Ncd6O6vMfv2Y8WxWxoeJMdFcmkRNruZSJ/dnM
wF6s6LkkmGWF2XtPN8OXwvVD69Ey6wDHIRmDlGbXVcoqPUJPjwm06u21wpeeALbH
1Brd2pDaH4P7IBqwP8Tpp4cAp6ZVRVFmINTcvml2oGLPmtOYJcuGEu5O1nEUjNUP
rwWLuKuKBCfZ4QE74JX+1yDMXHdoqtW+Dz+GeNhwZI5prrGaIzQz+Wy6Vrs=
-----END CERTIFICATE-----
`
	transportPrivate = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDHhIx7zZkTIk3h
elOLa9KDEBUHkI2khWBLOMOhRIV5SIH3me72bjb5XbmMFVwjRbpa4pCSUzZVPeWm
3QaNRnfjJAHvun7idX1tcK+iLwfXkxjhtuhEGF3TpN1OXODzaKcgnvAu9YrcQMFO
H65UyKm2c2r5XfRG+hcddKpvk5d+hPMxSA2JgfZXeyBSHpup9aKQpJ+Zqjh5HtOM
HkraqQt2+614eKzVudAeHPilgd79GvSispfpy3oeVXI5uin0VretZ1VY3VZsfHLp
wrI83jsz78A601GnuZtwCG5TkmvMKq06Rh3KfdYHm7x6DIalgoGAZd4/c43ScBjb
+ru/4NwhAgMBAAECggEAFiKvj2C9FfFdYKG4uSQqQ9455w/zlwgxKcdPdQnsIQuZ
V8YdS/voX3w1hMQt57/psAGo9oMC6Swn2X52JqBl1q59BILVZvyQAN9arQy4uwMX
5JrtY/isGDoXT4Vgc8DtoeHgVeVqFYudprQ/HCrzIUnm2WnCG0nN3Le/3Qcr7J5S
+xW+BalqeMhbHsS7VZ+isxrgeMsa8OLVSQ+omyesXAsi5DAtUDOe3Tw7VBi0UDth
uule8T7kN3dqReuMWQfRpLhFPEQZOwCSdpFnmi0DRBs470sMml7rUZV7OB7Vp3vD
chvoZyXO7r6OQPMPAyNZx/iC1bIh7ZrIGfz13ArLXQKBgQD+CmOflB/EHXm+GJPW
j5gt5LMxDMXT3CPR8go0AvqMJzZpDRa/XcXKjKZvjyoFC1kYRjI6DG2oC4jxa8fx
UuqUDM0i24g9V9yyIwj7e64349d7vNAjqYL23lpi6uZO/ubi3yX3bA0uIruz25l8
x11U/nA50ReO/OvXIQYuYWSUHwKBgQDJDoCuMxaiZLgeVrMmeg/8aLXEwIXDJeDp
K8XMLBD0/OhoUhC4L1scCUYRpiIJ3XzslQHeP3Z8Qg+5nnTVtbtVka1RIx3os2Kv
cuYkDmC5yyqBCE+57nT4+yLNa8TNLzwnGWoRxSKqnu8jXAkQDrhP/IJyiJc6qY01
xH6cDQaHvwKBgBjhxZ40sOPRi0IOQDSsvdgI5XAxcxLsJeoDTfKINCgUEyU47fhy
Y9QR8J9Oo2v5D5HsFjFPVFI4RwJ2bw/48hbsJg969x4jA+/CtLeFBqxcuZdaB/zm
NnidkLbNkR89ojmoZ5yTTbsuFbppEOCC2mZfwXg4PZl4tlTM3EEgsuw7AoGAD7Bc
BjviVkW5wFRPon7/5FhfZr0HMxUvmcJaqvX9VMCvegR9XYIEgAmROCtYmKB58RQn
kyosmsGk7H0a7NpDhgfaGGy/Frt4xewXXVTp41WhOXRmlEGxSwR90L3KG6DF9t8a
0cwqSlogmwfBhUlAxK0VmM5jzqYQaNOudYrmqY0CgYEAoIi8Bj1OczJS95eOKdqD
7sOwpu1P61TvC3SG/UbW9Ftpo4UUWVf4zNYtXugG1LsC+hCQsonWhuO+mdF+qO2j
M7EMy0B1icq9PGlvl+EmtTUQeW/eTMvlKa4/tukphZqXzimkbF9Rc1brT6zgeCKQ
HJ2zHQe3vcccMRFHglcf1eo=
-----END PRIVATE KEY-----
`
)

func disableTestExecutor_SetCertificates(t *testing.T) {

	t.Run("InvalidTransportCertificate", func(t *testing.T) {
		require := require.New(t)

		executor := NewExecutor()
		require.NotNil(executor)

		certificateSigning, err := authentication.NewCertificate(signingPublic, signingPrivate)
		require.NotNil(certificateSigning)
		require.NoError(err)

		certificateTransport, err := authentication.NewCertificate(signingPublic, signingPrivate)
		require.NotNil(certificateTransport)
		require.NoError(err)

		require.NoError(executor.SetCertificates(certificateSigning, certificateTransport))

		// https://ob19-rs1.o3bank.co.uk:4501/open-banking/v3.1/aisp
		res, err := resty.R().Get("https://ob19-rs1.o3bank.co.uk:4501/open-banking/v3.1/aisp")
		require.NotNil(res)
		require.NoError(err)

		require.Equal(
			"<html>\r\n<head><title>400 The SSL certificate error</title></head>\r\n<body bgcolor=\"white\">\r\n<center><h1>400 Bad Request</h1></center>\r\n<center>The SSL certificate error</center>\r\n<hr><center>nginx/1.14.1</center>\r\n</body>\r\n</html>\r\n",
			string(res.Body()),
		)
	})

	t.Run("ValidTransportCertificate", func(t *testing.T) {
		require := require.New(t)

		executor := NewExecutor()
		require.NotNil(executor)

		certificateSigning, err := authentication.NewCertificate(signingPublic, signingPrivate)
		require.NotNil(certificateSigning)
		require.NoError(err)

		certificateTransport, err := authentication.NewCertificate(transportPublic, transportPrivate)
		require.NotNil(certificateTransport)
		require.NoError(err)

		require.NoError(executor.SetCertificates(certificateSigning, certificateTransport))

		// https://ob19-rs1.o3bank.co.uk:4501/open-banking/v3.1/aisp
		res, err := resty.R().Get("https://ob19-rs1.o3bank.co.uk:4501/open-banking/v3.1/aisp")
		require.NotNil(res)
		require.NoError(err)

		require.Equal(
			"<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n<meta charset=\"utf-8\">\n<title>Error</title>\n</head>\n<body>\n<pre>Cannot GET /open-banking/v3.1/aisp</pre>\n</body>\n</html>\n",
			string(res.Body()),
		)
	})
}
