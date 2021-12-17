package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"gopkg.in/resty.v1"

	"github.com/OpenBankingUK/conformance-suite/pkg/client"
	"github.com/OpenBankingUK/conformance-suite/pkg/discovery"
	"github.com/OpenBankingUK/conformance-suite/pkg/generation"
	"github.com/OpenBankingUK/conformance-suite/pkg/model"
	"github.com/OpenBankingUK/conformance-suite/pkg/server"
	"github.com/OpenBankingUK/conformance-suite/pkg/version"

	"github.com/google/go-cmp/cmp"
)

var (
	logger = logrus.StandardLogger()
	update = flag.Bool("update", false, "update .golden files")
)

const (
	certFile = "../../certs/conformancesuite_cert.pem"
	keyFile  = "../../certs/conformancesuite_key.pem"
)

// init - this allows running the tests in debug mode, e.g.,:
//
// `LOG_HTTP_TRACE=true LOG_LEVEL=trace go test -v -count=1 -run='TestRun' ./...`
func init() {
	logger.SetNoLock()
	logger.SetFormatter(&prefixed.TextFormatter{
		DisableColors:    false,
		ForceColors:      true,
		TimestampFormat:  time.RFC3339,
		FullTimestamp:    true,
		DisableTimestamp: false,
		ForceFormatting:  true,
	})

	viper.SetDefault("LOG_LEVEL", "warn")
	viper.SetDefault("LOG_HTTP_TRACE", false)

	viper.SetEnvPrefix("")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func disableTestRun(t *testing.T) {
	debug := viper.GetBool("LOG_HTTP_TRACE")
	logLevel, err := logrus.ParseLevel(viper.GetString("LOG_LEVEL"))
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("logLevel=%+v", logLevel)
	t.Logf("debug=%+v", debug)

	logger.SetLevel(logLevel)
	resty.SetDebug(debug)

	logger := logger.WithFields(logrus.Fields{"test": "TestRun"})
	ver := version.NewGitHub(version.GitHubAPIRepository)
	validatorEngine := discovery.NewFuncValidator(model.NewConditionalityChecker())
	testGenerator := generation.NewGenerator()
	tlsValidator := discovery.NewStdTLSValidator(tls.VersionTLS11)
	journey := server.NewJourney(logger, testGenerator, validatorEngine, tlsValidator, false)

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
