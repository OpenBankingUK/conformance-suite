package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/server/models"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/version/mocks"
)

func TestServerImportHandlersPostImportReview(t *testing.T) {
	require := test.NewRequire(t)

	server := NewServer(testJourney(), nullLogger(), &mocks.Version{})
	defer func() {
		require.NoError(server.Shutdown(context.TODO()))
	}()
	require.NotNil(server)

	report, err := ioutil.ReadFile("./testdata/report.zip")
	require.NoError(err)

	importRequest := models.ImportRequest{
		Report: string(report),
	}
	requestJSON, err := json.MarshalIndent(importRequest, marshalIndentPrefix, marshalIndentindent)
	require.NoError(err)
	require.NotNil(requestJSON)

	// make the request
	code, body, headers := request(http.MethodPost, "/api/import/review", bytes.NewReader(requestJSON), server)

	// do assertions
	response := models.ImportReviewResponse{}
	responseJSON, err := json.MarshalIndent(response, marshalIndentPrefix, marshalIndentindent)
	require.NoError(err)
	require.NotNil(responseJSON)

	require.NotNil(body)
	bodyExpected := string(responseJSON)
	bodyActual := body.String()
	require.JSONEq(bodyExpected, bodyActual)

	require.Equal(http.StatusOK, code)
	// No gzip compression on this route
	require.Equal(http.Header{
		"Content-Type": []string{
			"application/json; charset=UTF-8",
		},
	}, headers)
}

func TestServerImportHandlersPostImportRerun(t *testing.T) {
	// TODO(mbana): need to implement rerun functionality and write tests.
	t.Skip()
}
