package server_test

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	// silence log output when running tests...
	logrus.SetLevel(logrus.WarnLevel)

	os.Exit(m.Run())
}
