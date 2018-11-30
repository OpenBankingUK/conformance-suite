package proxy

import (
	"fmt"
	"net/http"

	"github.com/fatih/color"
	"github.com/go-openapi/errors"
)

// Reporter Interface - types implementing this interface can be wired to receive
type Reporter interface {
	Success(req *http.Request)
	Error(req *http.Request, err error)
	Warning(req *http.Request, msg string)
	Report()
}

// LogReporter - type for "Logging" using the Reporter Interface
// Naming likey needs tidying a bit
type LogReporter struct {
}

// Success - log successful event
func (r *LogReporter) Success(req *http.Request) {
	fmt.Fprintf(color.Output, "%s %s %s\n",
		color.GreenString("✔"), req.Method, req.URL,
	)
}

// Error - log error event
func (r *LogReporter) Error(req *http.Request, err error) {
	fmt.Fprintf(color.Output, "%s %s %s",
		color.RedString("✗"), req.Method, req.URL,
	)
	if cErr, ok := err.(*errors.CompositeError); ok {
		for i, err := range cErr.Errors {
			fmt.Printf("  %d) %s\n", i+1, err)
		}

	} else {
		fmt.Printf("  => %s\n", err)
	}
}

// Warning - log warning event
func (r *LogReporter) Warning(req *http.Request, msg string) {
	fmt.Fprintf(color.Output, "%s %s %s\n",
		color.YellowString("!"), req.Method, req.URL,
	)
	fmt.Printf("  WARNING: %s\n", msg)
}

// Report - generic implemetnation
func (r *LogReporter) Report() {}
