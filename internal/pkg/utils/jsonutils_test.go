package pkgutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type teststruct struct {
	ID           string   `json:"@id"`
	Type         []string `json:"@type,omitempty"`
	Name         string   `json:"name"`
	Purpose      string   `json:"purpose"`
	Specref      string   `json:"specref"`
	Speclocation string   `json:"speclocation"`
}

// Simple testcase that capture the json representation of a Go struct in a byte slice
// the prints out the resulting bytes as a string
// also check the length of the resulting byte slice to see if its the expected length
func TestDumpStruct(t *testing.T) {
	tst := teststruct{ID: "123", Type: []string{"ABC"}, Name: "XYZ", Purpose: "test", Specref: "My reference", Speclocation: "Mylocation"}
	result := DumpJSON(tst)
	fmt.Println(string(result))
	assert.Equal(t, len(result), 163) // check the resulting byte slice is the expected length
}
