package main

import (
	"fmt"
	"os"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version"
	"github.com/pkg/errors"

	model "bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/tracer"
	resty "gopkg.in/resty.v1"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/server"
	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
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
	tracer.Silent = true
	resty.SetDebug(false)
}

func main() {
	logger := logrus.WithField("app", "server")
	ver := version.NewBitBucket(version.BitBucketAPIRepository)
	server := server.NewServer(logger, model.NewConditionalityChecker(), ver)
	server.HideBanner = true

	versionInfo(ver, logger)

	address := fmt.Sprintf("0.0.0.0:%s", getPort())
	logger.Infof("address -> http://%s", address)

	server.Logger.Fatal(server.Start(address))
}

func versionInfo(ver version.BitBucket, logger *logrus.Entry) {
	v, err := ver.VersionFormatter(version.FullVersion)
	if err != nil {
		logger.Error(errors.Wrap(err, "version.VersionFormatter()"))
	}
	msg, _, err := ver.UpdateWarningVersion(v)
	if err != nil {
		logger.Error(errors.Wrap(err, "version.UpdateWarningVersion()"))
	} else {
		logger.Info(msg)
	}
}
