package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/internal/pkg/version"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/server"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/tracer"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	resty "gopkg.in/resty.v1"
)

const (
	certFile = "./certs/conformancesuite_cert.pem"
	keyFile  = "./certs/conformancesuite_key.pem"
)

var (
	logger  = logrus.StandardLogger()
	rootCmd = &cobra.Command{
		Use:   "fcs_server",
		Short: "Functional Conformance Suite Server",
		Long: `A Fast and Flexible tool that enables implementers to test
interfaces and data endpoints against the Functional API
standard built with love by Open Banking and friends in Go.

Complete documentation is available at https://bitbucket.org/openbankingteam/conformance-suite`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := logger.WithField("app", "server")
			ver := version.NewBitBucket(version.BitBucketAPIRepository)

			printVersionInfo(ver, logger)

			echoServer := server.NewServer(logger, model.NewConditionalityChecker(), ver)
			echoServer.HideBanner = true
			server.PrintRoutesInfo(echoServer, logger)

			address := fmt.Sprintf("0.0.0.0:%d", viper.GetInt("port"))
			logger.Infof("listening on https://%s", address)
			return echoServer.StartTLS(address, certFile, keyFile)
		},
	}
)

func printVersionInfo(ver version.BitBucket, logger *logrus.Entry) {
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

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err)
		fmt.Fprint(os.Stderr, "\n")
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("log_level", "INFO", "Log level")
	rootCmd.PersistentFlags().Bool("log_tracer", false, "Enable tracer logging")
	rootCmd.PersistentFlags().Bool("log_http_trace", false, "Enable HTTP logging")
	rootCmd.PersistentFlags().Int("port", 8443, "Server port")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		fmt.Fprint(os.Stderr, err)
		fmt.Fprint(os.Stderr, "\n")
		os.Exit(1)
	}

	viper.SetEnvPrefix("")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	logger.SetNoLock()
	logger.SetFormatter(&prefixed.TextFormatter{
		DisableColors:    false,
		ForceColors:      true,
		TimestampFormat:  time.RFC3339,
		FullTimestamp:    true,
		DisableTimestamp: false,
		ForceFormatting:  true,
	})
	level, err := logrus.ParseLevel(viper.GetString("log_level"))
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		fmt.Fprint(os.Stderr, "\n")
		os.Exit(1)
	}
	logger.SetLevel(level)
	logrus.SetLevel(level)

	tracer.Silent = viper.GetBool("log_tracer")
	resty.SetDebug(viper.GetBool("log_http_trace"))
}
