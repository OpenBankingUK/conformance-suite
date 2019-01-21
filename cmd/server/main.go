package main

import (
	"fmt"
	"os"
	"time"

	model "bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/tracer"
	resty "gopkg.in/resty.v1"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/server"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

const (
	defaultPort = "8080"
)

func getPort() string {
	port := os.Getenv("PORT")
	if port != "" {
		return port
	}

	return defaultPort
}

func init() {
	logrus.SetFormatter(&prefixed.TextFormatter{
		DisableColors:    false,
		ForceColors:      true,
		TimestampFormat:  time.RFC3339,
		FullTimestamp:    true,
		DisableTimestamp: false,
		ForceFormatting:  true,
	})
	logrus.SetLevel(logrus.DebugLevel)
	logrus.StandardLogger().SetNoLock()
	tracer.Silent = false
	resty.SetDebug(true)
}

func main() {
	logger := logrus.WithField("app", "server")
	server := server.NewServer(logger, model.NewConditionalityChecker())
	server.HideBanner = true

	address := fmt.Sprintf("0.0.0.0:%s", getPort())
	logrus.Infof("address -> http://%s", address)

	server.Logger.Fatal(server.Start(address))
}
