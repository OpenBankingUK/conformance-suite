package test

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

// NullLogger - create a logger that discards output.
func NullLogger() *logrus.Entry {
	logger := logrus.New()
	logger.Out = ioutil.Discard
	return logger.WithField("app", "test")
}
