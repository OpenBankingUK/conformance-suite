package main

import (
	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/os"
	"fmt"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version"
	"github.com/pkg/errors"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/tracer"
	"gopkg.in/resty.v1"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/server"
	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
)

const defaultPort = "8080"

func init() {
	logrus.SetFormatter(&prefixed.TextFormatter{
		DisableColors:    false,
		ForceColors:      true,
		TimestampFormat:  time.RFC3339,
		FullTimestamp:    true,
		DisableTimestamp: false,
		ForceFormatting:  true,
	})
	logLevel := logrusLogLevel(os.GetEnvOrDefault("LOG_LEVEL", "INFO"))
	logrus.SetLevel(logLevel)
	logrus.StandardLogger().SetNoLock()
	tracer.Silent = true
	resty.SetDebug(false)
}

func main() {
	logger := logrus.StandardLogger().WithField("app", "server")
	ver := version.NewBitBucket(version.BitBucketAPIRepository)

	echoServer := server.NewServer(logger, model.NewConditionalityChecker(), ver)
	echoServer.HideBanner = true

	versionInfo(ver, logger)

	address := fmt.Sprintf("0.0.0.0:%s", os.GetEnvOrDefault("PORT", defaultPort))
	logger.Infof("address -> http://%s", address)
	server.RoutesInfo(echoServer, logger)

	logger.Fatal(echoServer.Start(address))
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

func logrusLogLevel(level string) logrus.Level {
	switch level {
	case "PANIC":
		return logrus.PanicLevel
	case "FATAL":
		return logrus.FatalLevel
	case "ERROR":
		return logrus.ErrorLevel
	case "WARN":
		return logrus.WarnLevel
	case "DEBUG":
		return logrus.DebugLevel
	case "INFO":
		fallthrough
	default:
		return logrus.InfoLevel
	}
}
