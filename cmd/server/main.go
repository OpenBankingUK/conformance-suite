package main

import (
	"fmt"
	"os"
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
	logrus.SetLevel(logrus.InfoLevel)
	logrus.StandardLogger().SetNoLock()
	tracer.Silent = true
	resty.SetDebug(false)
}

func main() {
	logger := logrus.WithField("app", "server")
	ver := version.NewBitBucket(version.BitBucketAPIRepository)

	echoServer := server.NewServer(logger, model.NewConditionalityChecker(), ver)
	echoServer.HideBanner = true

	versionInfo(ver, logger)

	address := fmt.Sprintf("0.0.0.0:%s", getEnvOrDefault("PORT", defaultPort))
	logger.Infof("address -> http://%s", address)
	server.RoutesInfo(echoServer, logger)

	logger.Fatal(echoServer.Start(address))
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

func getEnvOrDefault(key, defaultValue string) string {
	port, found := os.LookupEnv(key)
	if !found {
		return defaultValue
	}
	return port
}
