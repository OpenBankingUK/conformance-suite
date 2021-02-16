package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"strings"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/generation"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/manifest"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/server"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/tracer"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/version"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"gopkg.in/resty.v1"
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

			validatorEngine := discovery.NewFuncValidator(model.NewConditionalityChecker())
			testGenerator := generation.NewGenerator()
			tlsValidator := discovery.NewStdTLSValidator(tls.VersionTLS11)
			journey := server.NewJourney(logger, testGenerator, validatorEngine, tlsValidator, viper.GetBool("dynres"))

			echoServer := server.NewServer(journey, logger, ver)
			address := fmt.Sprintf("%s:%d", server.ListenHost, viper.GetInt("port"))
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
	rootCmd.PersistentFlags().Bool("disable_jws", false, "Disable JWS Signatures")
	rootCmd.PersistentFlags().Bool("dynres", false, "Use Dynamic Resource IDs - accounts")
	rootCmd.PersistentFlags().Bool("dumpcontexts", false, "Dump contexts when trace enabled")
	rootCmd.PersistentFlags().Bool("tlscheck", true, "enable tls version checking - default enabled")
	rootCmd.PersistentFlags().Bool("export_testcases", false, "Dump all testcases to console in CSV format")

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
		printConfigurationFlags()
		fmt.Fprint(os.Stderr, err)
		fmt.Fprint(os.Stderr, "\n")
		os.Exit(1)
	}
	logger.SetLevel(level)

	if viper.GetBool("export_testcases") {
		manifest.GenerateTestCaseListCSV()
		os.Exit(0)
	}

	tracer.Silent = !viper.GetBool("log_tracer")
	if viper.GetBool("log_to_file") {
		f, err := os.OpenFile("suite.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			// continue as normal
		} else {
			mw := f // io.MultiWriter(os.Stdout, f)
			logrus.SetOutput(mw)
			logger.SetFormatter(&prefixed.TextFormatter{
				DisableColors:    true,
				ForceColors:      false,
				TimestampFormat:  time.RFC3339,
				FullTimestamp:    true,
				DisableTimestamp: false,
				ForceFormatting:  true,
			})

		}
	}
	if viper.GetBool("log_http_file") {
		httpLogFile, err := os.OpenFile("http-trace.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			logrus.Warn("cannot set http trace file")
		} else {
			resty.SetLogger(httpLogFile)
		}
	}

	if viper.GetBool("disable_jws") {
		model.DisableJWS()
	}

	if viper.GetBool("dumpcontexts") {
		model.EnableContextDumps()
	}

	if viper.GetBool("tlscheck") == false {
		server.EnableTLSCheck(false)
	}

	resty.SetDebug(viper.GetBool("log_http_trace"))
	resty.SetRedirectPolicy(resty.FlexibleRedirectPolicy(15))
	printConfigurationFlags()
}

func printConfigurationFlags() {
	logger.WithFields(logrus.Fields{
		"log_level":        viper.GetString("log_level"),
		"log_tracer":       viper.GetBool("log_tracer"),
		"log_http_trace":   viper.GetBool("log_http_trace"),
		"log_http_file":    viper.GetBool("log_http_file"),
		"log_to_file":      viper.GetBool("log_to_file"),
		"port":             viper.GetInt("port"),
		"tracer.Silent":    tracer.Silent,
		"disable_jws":      viper.GetBool("disable_jws"),
		"dynres":           viper.GetBool("dynres"),
		"dumpcontexts":     viper.GetBool("dumpcontexts"),
		"tlscheck":         viper.GetBool("tlscheck"),
		"export_testcases": viper.GetString("export_testcases"),
	}).Info("configuration flags")
}
