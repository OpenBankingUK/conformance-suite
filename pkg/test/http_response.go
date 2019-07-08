package test

import (
	resty "gopkg.in/resty.v1"
)

// CreateHTTPResponse - helper to create an http response for test cases
// body
// parameters:
//   response code:
//   data[0] response body
//   data[1] http status
//   data[x*2 - 1] = http header where x > 1
//   data[x*2] = http header value
func CreateHTTPResponse(respcode int, data ...string) *resty.Response {
	var resBody string
	headers := make(map[string]string)

	if len(data) < 2 { // if no body is provided, create one
		resBody = ""
	}

	for k, v := range data {
		if k == 1 {
			resBody = v

			continue
		}
		if k%2 == 0 {
			continue
		}
		headers[data[k-1]] = v
	}
	mockedServer, mockedServerURL := HTTPServer(respcode, resBody, headers)
	defer mockedServer.Close()
	res, _ := resty.R().Get(mockedServerURL)

	return res
}
