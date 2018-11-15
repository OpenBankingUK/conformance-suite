package main

import (
	"fmt"
	"os"
	"time"

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
	logrus.SetLevel(logrus.InfoLevel)
	logrus.StandardLogger().SetNoLock()
}

func main() {
	server := server.NewServer()
	server.HideBanner = true

	address := fmt.Sprintf("0.0.0.0:%s", getPort())
	logrus.Infof("address -> http://%s", address)

	server.Logger.Fatal(server.Start(address))
}
