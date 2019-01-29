package server

// Note: Do not run the server tests in parallel.
// The server starts and stops proxy server at a particular port number.
// Starting and stopping proxy server at the same port cannot be done in parallel.

import (
	"bytes"
	"context"
	"flag"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version"
	versionmock "bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version/mocks"
	"github.com/stretchr/testify/mock"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	appConfigJSON        = appConfigJSONWithURL("https://rs.aspsp.ob.forgerock.financial:443")
	appConfigJSONWithURL = func(url string) string {
		return `{
    "softwareStatementId": "5b5a2008b093465496d238fc",
    "keyId": "d6c3f49c-7112-4c5c-9c9d-84926e992c74",
    "targetHost": "` + url + `",
    "verbose": true,
    "specLocation": "../../swagger/rw20test.json",
    "bindAddress": ":8989",
    "certTransport": "-----BEGIN CERTIFICATE-----\nmiIDkjCCAnqgAwIBAgIUfofLkR37LWwG11wRB70OFEDNwfcwDQYJKoZIhvcNAQELBQAwezELMAkG\nA1UEBhMCVUsxDTALBgNVBAgTBEF2b24xEDAOBgNVBAcTB0JyaXN0b2wxEjAQBgNVBAoTCUZvcmdl\nUm9jazEcMBoGA1UECxMTZm9yZ2Vyb2NrLmZpbmFuY2lhbDEZMBcGA1UEAxMQb2JyaS1leHRlcm5h\nbC1jYTAgFw0xNzA5MjExMTQ2MzZaGA8yMTE4MDgyODExNDYzNlowgYgxCzAJBgNVBAYTAlVLMQ0w\nCwYDVQQIEwRBdm9uMRAwDgYDVQQHEwdCcmlzdG9sMRIwEAYDVQQKEwlGb3JnZVJvY2sxITAfBgNV\nBAsTGDViNTA3MDY1YjA5MzQ2NTQ5NmQyMzhhODEhMB8GA1UEAxMYNWI1YTIwMDhiMDkzNDY1NDk2\nZDIzOGZjMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAi2XZZoHcZVC2zPING7xm8zr0\nT7AruqB+oQ/YOULW3mHI0oeflNpuQ45h8LzyqO+4HzO8xW1nSU7qke7y8LCFhvOltatyvIFDbq/t\nmF/Jg/KaIlFxe4KTFTl8crqfIirrOb+rz3qHxqbDNDPyFefNmmy0KhqcOEDe7TYSevAiJjG68yxl\nNS2/sT6/3wTAo8FcarTLHkSYNAuARghlDfhOxni7P0z7O8cOY5qhgRbyygFx8cxp0tGxHIIBjgxE\nO1FKgjFGn9TInfaHbKdGc+GCE4IG6FHwWsxDKEDVuPfUtLq3DydK6zu4u747+dP0ViGkZi19zki7\n93iOCL+QOIe96QIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQBYhgJ3BljZjTSlR66cRNk4xd6MeCz7\nfOhl8mucaXURGwI2y6/VH6+gVdkV/bJWhGp2dcO2DulXCtJefKkW0Y+cEs8YHzHnkyfneHPpNSL7\nhq6kQkpWJGKmge71NVFmODGqb8rGWYJMUtocTtcPq3o9EdS0nreEZmd+VPc2NQIm/0BACQ3IxxOW\n0RNu6CdodVm7xujdaiJJQyCQVsvSUXFAQY0ClWOQRAp7x9cQ2bN71rZxCpT9M/gb1UKlcR33qZ2g\nOZ3UhHaIi7CeMgWDNs9LuLV4565ERFHdG/xSkLLDf1UdhQfFFzyGBR0nZ7bbVVpqYTLEbbnoqUW6\nYQ7nVD63\n-----END CERTIFICATE-----\n",
    "certSigning": "-----BEGIN CERTIFICATE-----\nmiIDkjCCAnqgAwIBAgIUJgoHICdF1y4c1binOIG2IacLWC0wDQYJKoZIhvcNAQELBQAwezELMAkG\nA1UEBhMCVUsxDTALBgNVBAgTBEF2b24xEDAOBgNVBAcTB0JyaXN0b2wxEjAQBgNVBAoTCUZvcmdl\nUm9jazEcMBoGA1UECxMTZm9yZ2Vyb2NrLmZpbmFuY2lhbDEZMBcGA1UEAxMQb2JyaS1leHRlcm5h\nbC1jYTAgFw0xNzA4MjcxNDM4MTFaGA8yMTE4MDgwMzE0MzgxMVowgYgxCzAJBgNVBAYTAlVLMQ0w\nCwYDVQQIEwRBdm9uMRAwDgYDVQQHEwdCcmlzdG9sMRIwEAYDVQQKEwlGb3JnZVJvY2sxITAfBgNV\nBAsTGDViNTA3MDY1YjA5MzQ2NTQ5NmQyMzhhODEhMB8GA1UEAxMYNWI1YTIwMDhiMDkzNDY1NDk2\nZDIzOGZjMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAjwuGfH0I0g59o1kbd+kJgrfo\nQYwXaBnme5ozVEf4NC3/xO7Lk/f1wNYeNE78u712IW8HtEQPhUjhUz4bsck9p4nb5JLRIQPjvRRC\nOBPfPA+nLOCtUzpUIjmiZAac5Mxan0UqJfDvxsMXj3VatHKC1feknhIyQjyqKSbR5h0LoNjLDqnF\n9YdNIOoSkX9EdDuhPVp/JSdiNB8qBY+ARiPwPIkeauLPaBoAYypndzlLPZcNxZai+83xx1x3F9xt\nLZAyq89gO5be8mkv2aN7P0p2zt4vZHKfXSO4xHFIVRV2DA4ip/8M9rqG8HDbXiHnb016u0x2y8sb\nv/AThIccVD4z6QIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQAfJk5d6zMaTHgEtUidrUtkbofFxYC7\naCsnYJtf4+SIy28tQ6Et/yvIZKXsL8iPCdub0A4SXBto0xHRE4UcK+lpj/j7IktB4qPxWtrq99cL\nZGPPpYIa8HOThpBn9uoLcNxSXSpqhqWdn/cSxoo0+ynrXU2nziqMC2NKFgsTR5gc5wuLPfAIi5i5\nhb1VhYZXj7eujvZpxc+9lCWsMg7a1kSPmKodQ4ty+5MZJZ7TS6YcHIOmavu7nUhavmfXfKHKrA7E\n/n7b5X0AgFXL3QJa6s8jWQpYfvtpncmNKbjVbBwNX4bqg6z6DupaVE0JWMgTUBlmp4dF1bhMM53/\nFVWWLCSH\n-----END CERTIFICATE-----\n",
    "keySigning": "-----BEGIN PRIVATE KEY-----\nmiIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCPC4Z8fQjSDn2jWRt36QmCt+hB\njBdoGeZ7mjNUR/g0Lf/E7suT9/XA1h40Tvy7vXYhbwe0RA+FSOFTPhuxyT2nidvkktEhA+O9FEI4\nE988D6cs4K1TOlQiOaJkBpzkzFqfRSol8O/GwxePdVq0coLV96SeEjJCPKopJtHmHQug2MsOqcX1\nh00g6hKRf0R0O6E9Wn8lJ2I0HyoFj4BGI/A8iR5q4s9oGgBjKmd3OUs9lw3FlqL7zfHHXHcX3G0t\nkDKrz2A7lt7yaS/Zo3s/SnbO3i9kcp9dI7jEcUhVFXYMDiKn/wz2uobwcNteIedvTXq7THbLyxu/\n8BOEhxxUPjPpAgMBAAECggEAc6uLNbFZ55pGKEfO+Xjc8vJKAm8JImoHQZ3gsd98qp0jvRioUF/r\nPuMmC4BvyFSdaM3CuhdrQYk8g7auaGZlz8ufn8bFC2B80RHHtlcDZir2MUkBf1KkZASc9yuNxUom\nYbJpMcMR8XUi4SOxlEcg22rkl9n5ACzUIHC+vMhx9b8DfdwvtkK5zhFL1MnbT2lWEkdYGmdnR3tk\nwphHrnz2/Jf10LBkosmBWJxGN1zcjS+t4L7V8JfsxZz1idTjIzzkOk+DyAh0fr5t8a/3Zu8Tjw6E\nqiXtuqUBJuB5rAOyRjG08zrgz1PVMG/uKIF5A4XqBqiPB5KTAHddzG8Dpd1szQKBgQDAsEvtVEYV\nqS+VpnIPMlIA+UT2PpXvxL+4oEi6sKrCv7hVHAKLG6f7Sf3+mNPF3cegLGcSWVv190WBHC8uYI/b\nVH5PA//4ycD79pylxH3WgxRxuila1LxswgiVcRgidIxKziYDnqYbdexx6Y7myRuAXlyNoESZqBiw\nCB+uWJwYUwKBgQC+C4hBqk9id1xgueMFg9GcToLv9rM2+abPSfV+sMIvjWi6O72okn+rsTSRxsFI\nycAA6WKy6SWvylevmgS0S8MbDbzPO5QTThhLYexfBybBAi7i4c/ElycafHi0dA6SDrLYMvbgDPv6\npxp/RzJqhvwanotMtufqVB3KBa1mrzZLUwKBgQCvqAePcyPw2yrl4bZY5CadfJ/BW4yT52hfhr7G\ncgc5Qk1oSQCIj82y5uEFF4z29BbnjZLox01uDNzvtiHMxXpfF8eNgLf4tPOYvlhPRbDxvM0GYA8T\nHpwnCTuKAG9f+Z9rEkLVSetjXT0PGzuKaAsKGvuEoHXpHbRjxQQci+rAwQKBgBViGsy4qwH7SCuh\n/sdKE7Wwp870RSn0YS6Ftdexb8gF8zixLB/hi/f3kmCsqmbUPIRdvjs/PHxRGhiqDclzlNpga1Qt\n8fVSHi2tMPloRpYE9t2UZtpJ3559Tt+PB2yrtrfY1CpVi6yiTLrxedy+n3MnT6ksE2AsYsWuadpZ\n8JP9AoGBAKYsr2VHY1hjlAxGZEc7h4+tS0z/vU2jtlNvinHc6Nt/y5m6+S7OJtPMFHrNKPo5u32Y\nucZIo14LBCZERzjKynU9KUhYv9RAvdAO3JsqxXsTGXuHhWefMY1LNocvkXe2Vp45cMyYhWuCklVK\n/k5xhpZPLOyj9KOWm9AayLyIQgoZ\n-----END PRIVATE KEY-----\n",
    "keyTransport": "-----BEGIN PRIVATE KEY-----\nmiIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCLZdlmgdxlULbM8g0bvGbzOvRP\nsCu6oH6hD9g5QtbeYcjSh5+U2m5DjmHwvPKo77gfM7zFbWdJTuqR7vLwsIWG86W1q3K8gUNur+2Y\nX8mD8poiUXF7gpMVOXxyup8iKus5v6vPeofGpsM0M/IV582abLQqGpw4QN7tNhJ68CImMbrzLGU1\nLb+xPr/fBMCjwVxqtMseRJg0C4BGCGUN+E7GeLs/TPs7xw5jmqGBFvLKAXHxzGnS0bEcggGODEQ7\nUUqCMUaf1Mid9odsp0Zz4YITggboUfBazEMoQNW499S0urcPJ0rrO7i7vjv50/RWIaRmLX3OSLv3\neI4Iv5A4h73pAgMBAAECggEAXArHDolGhltSKGbWwu6Wp5nQwWFYrmTU1/JHWh/JEpNMb75/X2EP\nF2pHPxbzvxpx36Bgz5dacKw79FnrbYOQ1ze/vgHTw6YyaT6eylLNE2O35FsUbHIePhB0HKke9AeU\nw8/MWTXVuxRXpft6qH4jYwjNuVNSvU4QJF7kuuoeEksfe5A/lTHGP1RZsKALFztwVuaOr/MhUA7N\n8LbJToyq20/KKKU/hMvJ0LvjLM4duxODdVfDEJYu8VN3GlafhCme4JLzymGRz7uH8WSGngHf4E84\nKmmz/ewkkpACNlTVPvWxwEAc6GmHObAz31R05llqlykimecXM6hQn+bjfpWmAQKBgQDN2CYVpgML\nZU52uiVz65PCB+dHQUb5iqdr88KMXkzHzVUOQlBe82dDj5h4W8FQUROJ3z7ydVNCqbr+VPVbamx3\nssJDRRwZ7IBpaLAigcUbKiJVtfPNsyCzydkB03pHadGN/99RXHsPzKxqSGDSEuMiyql63sbqux2F\n57du8BXoKQKBgQCtXQAs0v8qlcQQZIX4iL23txHk0oPBiPQFFbUJRA2zHjRqP4DIf4PGeI/P8P0X\n/DnbpNord7rDzmgkaBfZ1o98aDCak0yaZ9v6yV2G/7h4hzXRAMHsQRlWyN8BZNMazAU3JtncV3c0\nNhf99XshfSQ436arG4L4ZpSXkj9uBYjfwQKBgDhFqsOoSpTG8RhL8wkpkY8tkfBMzBZT7Uj5rmmp\nLdxBKctoHYiXiddSXiApFUPbpje+q/qkUEqdE92LZDfFdDmUyL6TGgeMO96VG/GTAEtYzWIZB7lo\nCrybpZN2OKtlJkBnfqlDWvEKxueXOcC0IRvVw1cvp7lrxbphiiftwk9hAoGBAKaN6eQmllVga2xg\nV0Guha5h6IQRJ9og7GeSMkqDojHKvAqzldOKhpyAKZJacZ3AigmWOLB4J+uEexM3GmsDsviP1No8\n1+SkEXjASuWu+ph5Nl/kvWpwJJr3AyEAr7xX9E7HOZlyQqjbq3Mmi7Rh2RH29NYA6XQigXGZZO0b\nziNBAoGAQTK9Iy3n6dRNnenwCafPcvqU+k3Tigqqml3bpPWl6zfJo91P44OSZRLqUXAXkxbigt0h\nNyLaYrkRGqCbifchWHNd2+e4SsERTUBBfV5UgSIlhZm9Ys9u0ekUUV1FbnKYMIfBYs0F/XeoEPtq\n5BjmH/05RVI9GDz9Vzi/60SdVAo=\n-----END PRIVATE KEY-----\n",
    "client_credential_token": {
        "access_token": "",
        "expires_in": 0,
        "token_type": ""
    },
    "account_request_token": {
        "access_token": "",
        "expires_in": 0,
        "token_type": ""
    },
    "payment_request_token": {
        "access_token": "",
        "expires_in": 0,
        "token_type": ""
    }
}`
	}
)

// conditionalityCheckerMock - implements model.ConditionalityChecker interface for tests
type conditionalityCheckerMock struct {
}

// IsOptional - not used in discovery test
func (c conditionalityCheckerMock) IsOptional(method, endpoint string, specification string) (bool, error) {
	return false, nil
}

// Returns IsMandatory true for POST /account-access-consents, false for all other endpoint/methods.
func (c conditionalityCheckerMock) IsMandatory(method, endpoint string, specification string) (bool, error) {
	if method == "POST" && endpoint == "/account-access-consents" {
		return true, nil
	}
	return false, nil
}

// IsOptional - not used in discovery test
func (c conditionalityCheckerMock) IsConditional(method, endpoint string, specification string) (bool, error) {
	return false, nil
}

// Returns IsPresent true for valid GET/POST/DELETE endpoints.
func (c conditionalityCheckerMock) IsPresent(method, endpoint string, specification string) (bool, error) {
	if method == "GET" || method == "POST" || method == "DELETE" {
		return true, nil
	}
	return false, nil
}

func (c conditionalityCheckerMock) MissingMandatory(endpoints []model.Input, specification string) ([]model.Input, error) {
	return []model.Input{}, nil
}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	flag.Parse()

	// silence log output when running tests...
	logrus.SetLevel(logrus.WarnLevel)

	os.Exit(m.Run())
}

func TestServer(t *testing.T) {
	// Setup Version mock
	humanVersion := "0.1.2-RC1"
	warningMsg := "Version v0.1.2 of the Conformance Suite is out-of-date, please update to v0.1.3"
	formatted := "0.1.2"
	v := &versionmock.Version{}
	v.On("GetHumanVersion").Return(humanVersion)
	v.On("UpdateWarningVersion", mock.AnythingOfType("string")).Return(warningMsg, true, nil)
	v.On("VersionFormatter", mock.AnythingOfType("string")).Return(formatted, nil)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, v)

	t.Run("NewServer() returns non-nil value", func(t *testing.T) {
		assert.NotNil(t, server)
	})

	t.Run("GET / returns index.html", func(t *testing.T) {
		code, body, _ := request(http.MethodGet, "/", nil, server)

		assert.Equal(t, true, strings.HasPrefix(body.String(), "<!DOCTYPE html>"))
		assert.Equal(t, http.StatusOK, code)
	})

	t.Run("GET /favicon.ico returns favicon.ico", func(t *testing.T) {
		code, body, _ := request(http.MethodGet, "/favicon.ico", nil, server)

		assert.NotEmpty(t, body.String())
		assert.Equal(t, http.StatusOK, code)
	})

	require.NoError(t, server.Shutdown(context.TODO()))
}

// TestServerConformanceSuiteCallback - Test that `/conformancesuite/callback` returns `./web/dist/index.html`.
func TestServerConformanceSuiteCallback(t *testing.T) {
	require := require.New(t)

	server := NewServer(nullLogger(), conditionalityCheckerMock{}, mockVersionChecker())
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	// read the file we expect to be served.
	bytes, err := ioutil.ReadFile("./web/dist/index.html")
	require.NoError(err)
	bodyExpected := string(bytes)

	code, body, headers := request(
		http.MethodGet,
		`/conformancesuite/callback`,
		nil,
		server)

	// do assertions.
	require.Equal(http.StatusOK, code)
	require.Len(headers, 5)
	require.Equal("text/html; charset=utf-8", headers["Content-Type"][0])
	require.NotNil(body)

	bodyActual := body.String()
	require.Equal(bodyExpected, bodyActual)
}

func TestServerSkipper(t *testing.T) {
	require := require.New(t)

	echo := echo.New()
	context := echo.AcquireContext()
	defer func() {
		echo.ReleaseContext(context)
	}()

	paths := map[string]bool{
		"/index.html":                false,
		"/index.js":                  false,
		"/conformancesuite/callback": false,
		"/ipa":                       false,
		"/reggaws":                   false,
		"/api":                       true,
		"/swagger":                   true,
	}
	for path, shouldSkip := range paths {
		context.SetPath(path)         // set path on the Context
		isSkipped := skipper(context) // check if the `skipper` skips the path
		require.Equal(shouldSkip, isSkipped)
	}
}

// Generic util function for making test requests.
func request(method, path string, body io.Reader, server *Server) (int, *bytes.Buffer, http.Header) {
	req := httptest.NewRequest(method, path, body)
	rec := httptest.NewRecorder()

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	server.ServeHTTP(rec, req)

	return rec.Code, rec.Body, rec.HeaderMap
}

// nullLogger - create a logger that discards output.
func nullLogger() *logrus.Entry {
	logger := logrus.New()
	logger.Out = ioutil.Discard
	return logger.WithField("app", "test")
}

// mockVersionChecker - returns mock version checker.
func mockVersionChecker() version.Checker {
	// Setup Version mock
	humanVersion := "0.1.2-RC1"
	warningMsg := "Version v0.1.2 of the Conformance Suite is out-of-date, please update to v0.1.3"
	formatted := "0.1.2"

	v := &versionmock.Version{}
	v.On("GetHumanVersion").Return(humanVersion)
	v.On("UpdateWarningVersion", mock.AnythingOfType("string")).Return(warningMsg, true, nil)
	v.On("VersionFormatter", mock.AnythingOfType("string")).Return(formatted, nil)

	return v
}
