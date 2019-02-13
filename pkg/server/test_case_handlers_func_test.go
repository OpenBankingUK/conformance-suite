package server

import (
	"testing"
)

func TestGetTestCases(t *testing.T) {

	// Test Case no longer valid as initiates PSU consent flow which requires USER interface at the UI and apsps authorisation server

	// assert := assert.New(t)
	// server := NewServer(nullLogger(), conditionalityCheckerMock{}, &versionmock.Version{})
	// defer func() {
	// 	require.NoError(t, server.Shutdown(context.TODO()))
	// }()

	// discoveryModel, err := ioutil.ReadFile("../discovery/templates/ob-v3.0-ozone.json")
	// assert.NoError(err)
	// assert.NotNil(discoveryModel)

	// request(http.MethodPost, "/api/discovery-model",
	// 	strings.NewReader(string(discoveryModel)), server)

	// code, body, headers := request(http.MethodGet, "/api/test-cases",
	// 	nil, server)

	// // we should get back the test cases
	// assert.NotNil(body)
	// assert.Equal(http.StatusOK, code)
	// assert.Equal("application/json; charset=UTF-8", headers["Content-Type"][0])
	// assert.Contains(body.String(), `{"apiSpecification":{"name":"CustomTest-GetOzoneToken","url":""`)
	// assert.Contains(body.String(), `testCases":[{"@id":"#co0001","name":"Post Account Consent"`)
}
