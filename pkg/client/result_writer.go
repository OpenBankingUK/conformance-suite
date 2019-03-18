package client

import (
	"fmt"
	"io"
)

// ResultWriter writes testcase results to a writer
func ResultWriter(w io.Writer, results []TestCase) {
	var passMsg = map[bool]string{true: "PASS", false: "FAIL"}
	for _, result := range results {
		fmt.Fprintf(w, "=== %s: %s\n", passMsg[result.Pass], result.Id)
		if !result.Pass {
			fmt.Fprintf(w, "\t %s\n", result.Fail)
		}
	}
}
