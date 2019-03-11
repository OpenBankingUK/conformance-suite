package main

import (
	"crypto/tls"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type Connection struct {
	*http.Client
}

var errInsecure = errors.New("this client connection is insecure")

func NewConnection() (*Connection, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &Connection{
		&http.Client{
			Transport: tr,
			Timeout:   5 * time.Minute,
		},
	}, errInsecure
}
