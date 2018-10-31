package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/lib/server"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	ginkgo "github.com/onsi/ginkgo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	appConfigJSON = `{
    "softwareStatementId": "5b5a2008b093465496d238fc",
    "keyId": "d6c3f49c-7112-4c5c-9c9d-84926e992c74",
    "targetHost": "https://rs.aspsp.ob.forgerock.financial:443",
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
)

func TestServer(t *testing.T) {
	ginkgo.RunSpecs(t, "Server Suite")
}

// Generic util function for making test requests.
func request(method, path string, body io.Reader, s *server.Server) (int, *bytes.Buffer, http.Header) {
	req := httptest.NewRequest(method, path, body)
	rec := httptest.NewRecorder()

	s.ServeHTTP(rec, req)

	return rec.Code, rec.Body, rec.HeaderMap
}

var _ bool = ginkgo.Describe("Server", func() {
	var (
		s *server.Server
	)

	ginkgo.BeforeEach(func() {
		s = server.NewServer()
	})

	ginkgo.AfterEach(func() {
		if err := s.Shutdown(nil); err != nil {
			logrus.Errorln("AfterEach -> Shutdown err=", err)
		}
	})

	ginkgo.It("NewServer() returns non-nil value", func() {
		assert.NotNil(ginkgo.GinkgoT(), s)
	})

	ginkgo.Describe("/", func() {
		ginkgo.Context("GET", func() {
			ginkgo.It("Returns index.html", func() {
				code, body, _ := request(http.MethodGet, "/", nil, s)

				assert.Equal(ginkgo.GinkgoT(), true, strings.HasPrefix(body.String(), "<!DOCTYPE html>"))
				assert.Equal(ginkgo.GinkgoT(), http.StatusOK, code)
			})

			ginkgo.It("Returns favicon.ico", func() {
				code, body, _ := request(http.MethodGet, "/favicon.ico", nil, s)

				assert.NotEmpty(ginkgo.GinkgoT(), body.String())
				assert.Equal(ginkgo.GinkgoT(), http.StatusOK, code)
			})

			ginkgo.It(`Returns {"message":"Not Found"} when file does not exist`, func() {
				code, body, _ := request(http.MethodGet, "/NotFound.ico", nil, s)

				assert.Equal(ginkgo.GinkgoT(), http.StatusNotFound, code)
				assert.Equal(ginkgo.GinkgoT(), `{"message":"Not Found"}`, body.String())
			})
		})
	})

	ginkgo.Describe("/api/health", func() {
		ginkgo.Context("GET", func() {
			ginkgo.It("When successful returns OK", func() {
				code, body, _ := request(http.MethodGet, "/api/health", nil, s)

				assert.Equal(ginkgo.GinkgoT(), "OK", body.String())
				assert.Equal(ginkgo.GinkgoT(), http.StatusOK, code)
			})
		})
	})

	ginkgo.Describe("/api/validation-runs", func() {
		ginkgo.Context("POST", func() {
			ginkgo.It("When successful returns validation run ID in JSON and nil error", func() {
				code, body, headerMap := request(http.MethodPost, "/api/validation-runs", nil, s)

				assert.NotNil(ginkgo.GinkgoT(), body)
				var responseBody server.ValidationRunsResponse
				json.Unmarshal(body.Bytes(), &responseBody)
				id, err := uuid.Parse(responseBody.ID)
				assert.NoError(ginkgo.GinkgoT(), err)
				assert.Equal(ginkgo.GinkgoT(), id.String(), responseBody.ID)

				assert.Equal(ginkgo.GinkgoT(), echo.MIMEApplicationJSONCharsetUTF8, headerMap.Get(echo.HeaderContentType))

				assert.Equal(ginkgo.GinkgoT(), http.StatusAccepted, code)
			})
		})
	})

	ginkgo.Describe("/api/validation-runs/${id}", func() {
		ginkgo.Context("GET", func() {
			ginkgo.It("When succesful returns OK ", func() {
				id := "c243d5b6-32f0-45ce-a516-1fc6bb6c3c9a"
				code, body, headerMap := request(
					http.MethodGet,
					fmt.Sprintf("/api/validation-runs/%s", id),
					nil,
					s,
				)

				assert.NotNil(ginkgo.GinkgoT(), body)
				var responseBody server.ValidationRunsIDResponse
				json.Unmarshal(body.Bytes(), &responseBody)
				assert.Equal(ginkgo.GinkgoT(), id, responseBody.Status)

				assert.Equal(ginkgo.GinkgoT(), echo.MIMEApplicationJSONCharsetUTF8, headerMap.Get(echo.HeaderContentType))

				assert.Equal(ginkgo.GinkgoT(), http.StatusOK, code)
			})
		})
	})

	ginkgo.Describe("/api/config", func() {
		ginkgo.Context("POST", func() {
			ginkgo.It("Can POST config", func() {
				// assert server isn't started before call
				frontendProxy, _ := url.Parse("http://0.0.0.0:8989/open-banking/v2.0/accounts")
				_, err := http.Get(frontendProxy.String())
				assert.Error(ginkgo.GinkgoT(), err)

				// create the request to post the config
				// this should start the proxy
				req := httptest.NewRequest(http.MethodPost, "/api/config", strings.NewReader(appConfigJSON))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				rec := httptest.NewRecorder()

				// do the request
				s.ServeHTTP(rec, req)

				assert.NotNil(ginkgo.GinkgoT(), rec.Body)
				assert.Equal(ginkgo.GinkgoT(), appConfigJSON, rec.Body.String())
				assert.Equal(ginkgo.GinkgoT(), http.StatusOK, rec.Code)

				// check the proxy is up now, we should hit the forgerock server
				resp, err := http.Get(frontendProxy.String())
				body, err := ioutil.ReadAll(resp.Body)
				assert.NoError(ginkgo.GinkgoT(), err)
				assert.Equal(ginkgo.GinkgoT(), http.StatusBadRequest, resp.StatusCode)

				// assert that the body matches a certain regex
				assert.Regexp(
					ginkgo.GinkgoT(),
					regexp.MustCompile(`^{"Code":"OBRI.FR.Request.Invalid","Id":".*","Message":"An error happened when parsing the request arguments","Errors":\[{"ErrorCode":"UK.OBIE.Header.Missing","Message":"Missing request header 'x-fapi-financial-id' for method parameter of type String","Url":"https://docs.ob.forgerock.financial/errors#UK.OBIE.Header.Missing"}\]}$`),
					string(body),
				)
			})

			ginkgo.It("cannot POST config twice without first deleting it", func() {
				// assert server isn't started before call
				frontendProxy, _ := url.Parse("http://0.0.0.0:8989/open-banking/v2.0/accounts")
				_, err := http.Get(frontendProxy.String())
				assert.Error(ginkgo.GinkgoT(), err)

				// create the request to post the config
				// this should start the proxy
				req := httptest.NewRequest(http.MethodPost, "/api/config", strings.NewReader(appConfigJSON))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()

				// do the request
				s.ServeHTTP(rec, req)

				assert.NotNil(ginkgo.GinkgoT(), rec.Body)
				assert.Equal(ginkgo.GinkgoT(), appConfigJSON, rec.Body.String())
				assert.Equal(ginkgo.GinkgoT(), http.StatusOK, rec.Code)

				// create another request to POST the config again
				// this should fail because a DELETE need to happen first.
				req = httptest.NewRequest(http.MethodPost, "/api/config", strings.NewReader(appConfigJSON))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec = httptest.NewRecorder()
				// do the request
				s.ServeHTTP(rec, req)

				assert.NotNil(ginkgo.GinkgoT(), rec.Body)
				assert.Equal(
					ginkgo.GinkgoT(),
					"{\n    \"error\": \"listen tcp :8989: bind: address already in use\"\n}",
					rec.Body.String(),
				)
				assert.Equal(ginkgo.GinkgoT(), http.StatusBadRequest, rec.Code)
			})
		})

		ginkgo.Context("DELETE", func() {
			ginkgo.It("DELETE stops the proxy", func() {
				// assert server isn't started before call
				frontendProxy, _ := url.Parse("http://0.0.0.0:8989/open-banking/v2.0/accounts")
				_, err := http.Get(frontendProxy.String())
				assert.Error(ginkgo.GinkgoT(), err)

				// create the request to post the config
				// this should start the proxy
				req := httptest.NewRequest(http.MethodPost, "/api/config", strings.NewReader(appConfigJSON))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()

				// do the request
				s.ServeHTTP(rec, req)

				assert.NotNil(ginkgo.GinkgoT(), rec.Body)
				assert.Equal(ginkgo.GinkgoT(), appConfigJSON, rec.Body.String())
				assert.Equal(ginkgo.GinkgoT(), http.StatusOK, rec.Code)

				// create request to delete config
				req = httptest.NewRequest(http.MethodDelete, "/api/config", nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec = httptest.NewRecorder()
				// do the request
				s.ServeHTTP(rec, req)

				assert.NotNil(ginkgo.GinkgoT(), rec.Body)
				assert.Equal(
					ginkgo.GinkgoT(),
					"",
					rec.Body.String(),
				)
				assert.Equal(ginkgo.GinkgoT(), http.StatusOK, rec.Code)

				// call proxy and assert it is no longer up
				// check the proxy is up now, we should hit the forgerock server
				resp, err := http.Get(frontendProxy.String())
				assert.Equal(
					ginkgo.GinkgoT(),
					`Get http://0.0.0.0:8989/open-banking/v2.0/accounts: dial tcp 0.0.0.0:8989: connect: connection refused`,
					err.Error(),
				)
				assert.Nil(ginkgo.GinkgoT(), resp)
			})
		})
	})
})
