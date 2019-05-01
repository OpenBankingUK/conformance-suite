package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// Service is a gateway to backend services provided by FCS
type Service interface {
	Version() (VersionResponse, error)
	Run(discoveryFile, configFile, exportConfig string) ([]TestCase, error)
}

const (
	setDiscoveryModelPath = "/api/discovery-model"
	setConfigPath         = "/api/config/global"
	exportReport          = "/api/export"
	generateTestCases     = "/api/test-cases"
	runTestCases          = "/api/run"
	runTestCasesResultsWS = "/api/run/ws"
	versionPath           = "/api/version"
)

// service is the implementation using HTTP client for consuming FCS services
type service struct {
	conn   *Connection
	host   string
	wsHost string
}

func NewService(host, wsHost string, conn *Connection) service {
	return service{
		conn:   conn,
		host:   host,
		wsHost: wsHost,
	}
}

type VersionResponse struct {
	Version string `json:"version"`
	Message string `json:"message"`
	Update  bool   `json:"update"`
}

func (s service) Run(discovery, config, report string) ([]TestCase, error) {
	err := s.setDiscoveryModel(discovery)
	if err != nil {
		return nil, err
	}

	err = s.setConfig(config)
	if err != nil {
		return nil, err
	}

	err = s.TestCases()
	if err != nil {
		return nil, err
	}

	resultsChan := make(chan TestCase)
	endedChan := make(chan struct{})

	err = s.runTests(resultsChan, endedChan)
	if err != nil {
		return nil, err
	}

	results, err := aggregateResults(resultsChan, endedChan)
	if err != nil {
		return nil, err
	}

	err = s.exportReport(report)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func aggregateResults(resultChan chan TestCase, endedChan chan struct{}) ([]TestCase, error) {
	var results []TestCase
	const timeoutRunningTests = 5 * time.Minute

	deadline := time.NewTicker(timeoutRunningTests)
	defer deadline.Stop()
	for {
		select {
		case result := <-resultChan:
			results = append(results, result)

		case <-endedChan:
			return results, nil

		case <-deadline.C:
			return nil, errors.New("timout running tests")
		}
	}
}

func (s service) Version() (VersionResponse, error) {
	response, err := s.conn.Get(s.host + versionPath)
	if err != nil {
		return VersionResponse{}, errors.Wrap(err, "getting version")
	}

	if response.StatusCode != http.StatusOK {
		return VersionResponse{}, fmt.Errorf("unexpected status code from getting version %d", response.StatusCode)
	}

	versionResponse := VersionResponse{}
	err = json.NewDecoder(response.Body).Decode(&versionResponse)
	if err != nil {
		return VersionResponse{}, err
	}

	return versionResponse, nil
}

func (s service) setDiscoveryModel(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return errors.Wrap(err, "setting discovery model")
	}

	response, err := s.conn.Post(s.host+setDiscoveryModelPath, "application/json", file)
	if err != nil {
		return errors.Wrap(err, "setting discovery model")
	}

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code setting discovery model %d", response.StatusCode)
	}

	return nil
}

func (s service) setConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return errors.Wrap(err, "setting config")
	}

	response, err := s.conn.Post(s.host+setConfigPath, "application/json", file)
	if err != nil {
		return errors.Wrap(err, "setting config")
	}

	if response.StatusCode != http.StatusCreated {
		responseBody, err := ioutil.ReadAll(response.Body)
		defer response.Body.Close()
		if err != nil {
			return errors.Wrap(err, "reading error response from setting config")
		}

		return fmt.Errorf("unexpected status code setting config %d, %s", response.StatusCode, string(responseBody))
	}

	return nil
}

func (s service) TestCases() error {
	response, err := s.conn.Get(s.host + generateTestCases)
	if err != nil {
		return errors.Wrap(err, "generating test cases")
	}

	if response.StatusCode != http.StatusOK {
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return errors.Wrap(err, "reading error response from generating test cases")
		}

		return errors.Errorf("unexpected status code generating test cases: %d, %s", response.StatusCode, string(responseBody))
	}

	return nil
}

func (s service) runTests(resultChan chan<- TestCase, endedChan chan<- struct{}) error {
	err := s.handleResults(resultChan, endedChan)
	if err != nil {
		return errors.Wrap(err, "running test cases")
	}

	response, err := s.conn.Post(s.host+runTestCases, "application/test", nil)
	if err != nil {
		return errors.Wrap(err, "running test cases")
	}

	if response.StatusCode != http.StatusCreated {
		var responseBody []byte
		_, err = response.Body.Read(responseBody)
		if err != nil {
			return errors.Wrap(err, "reading error response from running test cases")
		}

		return errors.Errorf(" unexpected status code generating test cases: %d, %s", response.StatusCode, responseBody)
	}

	return nil
}

func (s service) exportReport(reportConfig string) error {
	file, err := os.Open(reportConfig)
	if err != nil {
		return errors.Wrap(err, "export report")
	}

	response, err := s.conn.Post(s.host+exportReport, "application/json", file)
	if err != nil {
		return errors.Wrap(err, "export config")
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code export report config %d", response.StatusCode)
	}

	return nil
}

func (s service) handleResults(resultChan chan<- TestCase, endedChan chan<- struct{}) error {
	c, err := s.wsDialer()
	if err != nil {
		return err
	}

	cleanup := func() {
		close(resultChan)
		errCleanup := c.Close()
		if errCleanup != nil {
			fmt.Printf("error closing ws: %v\n", errCleanup)
		}
	}

	msgProcessor := newMsgProcessor(msgHandlersChain(resultChan, endedChan))
	go func() {
		defer cleanup()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				return
			}
			errProcessor := msgProcessor.process(message)
			if errProcessor != nil {
				endedChan <- struct{}{}
				return
			}
		}
	}()

	return nil
}

func (s service) wsDialer() (*websocket.Conn, error) {
	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 5 * time.Second,
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
	}
	c, _, err := dialer.Dial(s.wsHost+runTestCasesResultsWS, nil)
	return c, err
}
