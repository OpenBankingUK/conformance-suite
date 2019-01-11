package tracer

/*
The tracer module contain a simple application trace facility whose purpose is to capture
error and execution messages to add detail to the reporting output.

The basic naive implementation consists of:-
 - an ApplicationTrace channel which receives a string message
 - a goroutine which listens to the channel and outputs anything received to console - initial implementation
 - AppInfo which add a trace message
 - AppEntry - which effectively indents the messages for easy trace reading in the initial implementation
 - AppExit - which reduces an indentation level on function return

 So the idea is to produce and initial trace which some structure as a starting point for a richer execution flow capture
 to feed reporting, so a user can follow what executed, and knows exactly where an error occurred and why.

*/
import (
	"fmt"
	"strings"
)

// Tracer info capture interface - capturing testcase execution info - expects to be extended/improved over time
type Tracer interface {
	AppMsg(objType, msg, serialisedObj string) // general application flow information
	AppErr(objType, msg, serialisedObj string) // application errors
	AppEntry(objType, msg string)              // function entry
	AppExit(objType, msg string)               // function exit
}

var (
	// ApplicationTrace channel
	applicationTrace chan string
	// Silent running
	Silent = true
	// Current text indent
	indents = 0
)

// Model Package Initialisation
// Core tracing channel
func init() {
	applicationTrace = make(chan string)
	go appTraceLisener()
}

// simple trace listener to dump out to console
func appTraceLisener() {
	for {
		msg := <-applicationTrace
		fmt.Println(msg)
	}
}

// AppMsg - generic trace fuction
func AppMsg(objtype, msg, objdump string) {
	if Silent {
		return
	}
	applicationTrace <- fmt.Sprintf("%s[%s]MSG %s", indent(), objtype, msg)
}

// AppErr - generic trace errpr fuction
func AppErr(objtype, msg, objdump string) {
	if Silent {
		return
	}
	applicationTrace <- fmt.Sprintf("%s[%s]ERR %s", indent(), objtype, msg)
}

// AppEntry - application level trace
func AppEntry(objtype, msg string) {
	AppMsg(objtype, msg, "")
	addIndent(1)
}

// AppExit - application level trace
func AppExit(objtype, msg string) {
	addIndent(-1)
	AppMsg(objtype, msg, "")
}

// Indent -
func indent() string {
	if indents < 0 {
		return ""
	}
	return strings.Repeat(" ", indents*3)
}

// addIndent -
func addIndent(i int) {
	indents += i
}