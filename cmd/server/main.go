package main

import (
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/os"
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/server"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/tracer"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	resty "gopkg.in/resty.v1"
)

var (
	logger = logrus.StandardLogger()
)

const (
	certFile = "./certs/conformancesuite_cert.pem"
	keyFile  = "./certs/conformancesuite_key.pem"
)

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
	level, err := logrus.ParseLevel(os.GetEnvOrDefault("LOG_LEVEL", "INFO"))
	if err != nil {
		logger.Fatal(err)
	}
	logger.SetLevel(level)

	tracer.Silent = os.GetEnvOrDefault("TRACER", "off") == "off"
	resty.SetDebug(os.GetEnvOrDefault("HTTP_TRACE", "off") == "DEBUG")
}

func main() {
	logger := logger.WithField("app", "server")
	ver := version.NewBitBucket(version.BitBucketAPIRepository)

	versionInfo(ver, logger)

	echoServer := server.NewServer(logger, model.NewConditionalityChecker(), ver)
	echoServer.HideBanner = true
	server.RoutesInfo(echoServer, logger)

	logger.Info("listening on https://0.0.0.0")
	logger.Fatal(echoServer.StartTLS(":443", certFile, keyFile))
}

func versionInfo(ver version.BitBucket, logger *logrus.Entry) {
	v, err := ver.VersionFormatter(version.FullVersion)
	if err != nil {
		logger.Error(errors.Wrap(err, "version.VersionFormatter()"))
		return
	}
	msg, _, err := ver.UpdateWarningVersion(v)
	if err != nil {
		logger.Error(errors.Wrap(err, "version.UpdateWarningVersion()"))
		return
	}
	logger.Info(msg)
}
