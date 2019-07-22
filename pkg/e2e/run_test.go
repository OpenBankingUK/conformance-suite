package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"gopkg.in/resty.v1"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/client"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/server"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/test"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/version"

	"github.com/google/go-cmp/cmp"
)

const (
	certFile = "../../certs/conformancesuite_cert.pem"
	keyFile  = "../../certs/conformancesuite_key.pem"
)

var update = flag.Bool("update", false, "update .golden files")

// init - this allows running the tests in debug mode, e.g.,:
//
// `LOG_HTTP_TRACE=true go test -v -count=1 -run='TestRun' ./...`
func init() {
	viper.SetEnvPrefix("")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func TestRun(t *testing.T) {
	debug := viper.GetBool("LOG_HTTP_TRACE")
	resty.SetDebug(debug)

	logger := test.NullLogger()

	ver := version.NewBitBucket(version.BitBucketAPIRepository)
	validatorEngine := discovery.NewFuncValidator(model.NewConditionalityChecker())
	testGenerator := generation.NewGenerator()
	journey := server.NewJourney(logger, testGenerator, validatorEngine)

	echoServer := server.NewServer(journey, logger, ver)

	go func() {
		require.EqualError(t, echoServer.StartTLS(":0", certFile, keyFile), "http: Server closed")
	}()
	time.Sleep(100 * time.Millisecond)

	defer func() {
		errEcho := echoServer.Shutdown(context.TODO())
		require.NoError(t, errEcho)
	}()

	tcpAddr, ok := echoServer.TLSListener.Addr().(*net.TCPAddr)
	require.True(t, ok)
	serverHost := fmt.Sprintf("localhost:%d", tcpAddr.Port)
	waitForServerReady(t, serverHost)

	insecureConn, err := client.NewConnection()
	if err == client.ErrInsecure {
		logger.Println("server's certificate chain and host name not verified")
	} else {
		require.NoError(t, err)
	}
	service := client.NewService("https://"+serverHost, "wss://"+serverHost, insecureConn)

	goldenFile := filepath.Join("testdata", "ozone-results.golden")

	results, err := service.Run(
		"../discovery/templates/ob-v3.1-ozone-headless.json",
		"../../config/config-ozone-run_test.json",
		"../../config/report.json")
	require.NoError(t, err)

	w := bytes.NewBufferString("")
	client.ResultWriter(w, results)

	if *update {
		t.Log("update golden file")
		require.NoError(t, ioutil.WriteFile(goldenFile, w.Bytes(), 0644), "failed to update golden file")
	}

	expected, err := ioutil.ReadFile(goldenFile)
	require.NoError(t, err, "failed reading .golden")

	if string(expected) != w.String() {
		t.Logf("expected=%q", string(expected))
		t.Logf("actual=%q", w.String())

		t.Log(cmp.Diff(string(expected), w.String()))
		t.Fail()
	}
}

func waitForServerReady(t *testing.T, address string) {
	c, err := client.NewConnection()
	require.Error(t, err)
	ready := false
	for !ready {
		r, err := c.Get("https://" + address + "/api/ping")
		if err == nil && r.StatusCode == http.StatusOK {
			t.Log(r.StatusCode)
			ready = true
		} else {
			time.Sleep(time.Millisecond * 100)
		}
	}
}
