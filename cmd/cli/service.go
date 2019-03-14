package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"time"
)

// Service is a gateway to backend services provided by FCS
type Service interface {
	Version() (VersionResponse, error)
	SetDiscoveryModel(filename string) error
	SetConfig(filename string) error
	TestCases() error
	RunTests(resultChan chan<- TestCase, endedChan chan<- struct{}) error
}

const (
	setDiscoveryModelPath = "/api/discovery-model"
	setConfigPath         = "/api/config/global"
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

func newService(host, wsHost string, conn *Connection) Service {
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

func (s service) SetDiscoveryModel(filename string) error {
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

func (s service) SetConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return errors.Wrap(err, "setting config")
	}

	response, err := s.conn.Post(s.host+setConfigPath, "application/json", file)
	if err != nil {
		return errors.Wrap(err, "setting config")
	}

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code setting config %d", response.StatusCode)
	}

	return nil
}

func (s service) TestCases() error {
	response, err := s.conn.Get(s.host + generateTestCases)
	if err != nil {
		return errors.Wrap(err, "generating test cases")
	}

	if response.StatusCode != http.StatusOK {
		var responseBody []byte
		_, err = response.Body.Read(responseBody)
		if err != nil {
			return errors.Wrap(err, "reading error response from generating test cases")
		}

		return errors.Errorf("unexpected status code generating test cases: %d, %s", response.StatusCode, responseBody)
	}

	return nil
}

func (s service) RunTests(resultChan chan<- TestCase, endedChan chan<- struct{}) error {
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
