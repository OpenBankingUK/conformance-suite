package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/appconfig"
	"bitbucket.org/openbankingteam/conformance-suite/lib/server"
	"bitbucket.org/openbankingteam/conformance-suite/proxy"
	"github.com/go-openapi/loads"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

const (
	defaultPort      = "8080"
	defaultConfigDir = "config"
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
		DisableColors: false,
		// TimestampFormat: "2006-01-02 15:04:05.000",
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
		ForceFormatting: true,
	})
	logrus.SetLevel(logrus.DebugLevel)
}

// Run the proxy at the address specified by "bind"
// Requests get sent to the target server identifyed by proxy.Target()
// configure some channels to handle shutdown/interrupts
func serve(proxy *proxy.Proxy, bind string) error {
	s := http.Server{
		Addr:    bind,
		Handler: proxy.Router(),
	}

	// Run server s.ListenAndServe on a goroutine
	errChannel := make(chan error)
	go func() {
		logrus.WithFields(logrus.Fields{
			"bind":   bind,
			"target": proxy.Target()}).
			Info("ObProxy is listening:")
		errChannel <- s.ListenAndServe()
	}()

	// handle interrupted
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)

	select {
	case err := <-errChannel: // Error from listen&serve - exit
		return err
	case s := <-sigChannel: // Interrupted - exit
		logrus.Warningf("%s", s)
		return nil
	}
}

// main - kick off proxy by:
// loading the spec,
// creating a new proxy configured with
//    - bind address
//    - swagger specification location
//    - target host (aspsp resource server)
//    - verbosity
// configure an default logreport
func main() {
	go startServeWeb()          // run what was the previous 'main' on another goroutine
	time.Sleep(1 * time.Second) // wait for server to initialise

	logrus.Println("OB Logging Proxy")

	appconfig, err := appconfig.LoadAppConfiguration(defaultConfigDir)
	if err != nil {
		log.Fatal(err)
	}

	appconfig.PrintAppConfig()
	doc, err := loads.Spec(appconfig.Spec)
	if err != nil {
		log.Fatal(err)
	}

	proxy, err := proxy.New(doc.Spec(), &proxy.LogReporter{},
		proxy.WithTarget(appconfig.TargetHost),
		proxy.WithVerbose(appconfig.Verbose),
		proxy.WithAppConfig(appconfig),
	)
	if err != nil {
		log.Fatal(err)
	}

	// start serving the proxy - and don't return unless there is a problem/exit
	if err := serve(proxy, appconfig.Bind); err != nil {
		log.Println(err)
	}

	// Report PendingOperations - part of shutdown tidyup
	logrus.Debugln("Pending Operations:")
	logrus.Debugln("------------------")
	for i, op := range proxy.PendingOperations() {
		logrus.Debugf("%03d) id=%s", i+1, op.ID)
	}
}

// Start the web server that will serve the Vue.js Single Page Application
func startServeWeb() {
	address := fmt.Sprintf("localhost:%s", getPort())
	logrus.Debugf("address=http://%s", address)

	e := server.NewServer()
	e.Logger.Fatal(e.Start(address))
}
