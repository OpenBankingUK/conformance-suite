package main

import (
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/test"
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/client"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/server"
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

const (
	certFile = "../../certs/conformancesuite_cert.pem"
	keyFile  = "../../certs/conformancesuite_key.pem"
)

var update = flag.Bool("update", false, "update .golden files")

func TestRun(t *testing.T) {
	logger := test.NullLogger()

	freePort, err := getFreePort()
	require.NoError(t, err)

	ver := version.NewBitBucket(version.BitBucketAPIRepository)
	validatorEngine := discovery.NewFuncValidator(model.NewConditionalityChecker())
	testGenerator := generation.NewGenerator()
	journey := server.NewJourney(logger, testGenerator, validatorEngine)
	address := fmt.Sprintf("%s:%d", "127.0.0.1", freePort)

	echoServer := server.NewServer(journey, logger, ver)

	go func() {
		logger.Debugf("starting server %s", address)
		errEcho := echoServer.StartTLS(address, certFile, keyFile)
		require.NoError(t, errEcho)
	}()
	defer func() {
		errEcho := echoServer.Shutdown(context.TODO())
		if errEcho != nil {
			require.NoError(t, errEcho)
		}
	}()

	waitForServerReady(t, address)

	insecureConn, err := client.NewConnection()
	if err == client.ErrInsecure {
		logger.Println("server's certificate chain and host name not verified")
	} else {
		require.NoError(t, err)
	}
	service := client.NewService("https://"+address, "wss://"+address, insecureConn)

	goldenFile := filepath.Join("testdata", "ozone-results.golden")

	results, err := service.Run(
		"../discovery/templates/ob-v3.1-ozone-headless.json",
		"../../config/config-ozone.json",
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
		t.Log(string(expected))
		t.Log(w.String())
		t.Log(cmp.Diff(string(expected), w.String()))
		t.Fail()
	}
}

// GetFreePort asks the kernel for a free open port that is ready to use.
func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	port := l.Addr().(*net.TCPAddr).Port
	err = l.Close()
	return port, err
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
