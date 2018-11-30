package pkgutils

import (
	"encoding/json"
)

// jsonutils contains common functions that manipulate JSON

// DumpJSON - output formatted json from a go struct to byte array which can
// then be output to the console for example
//
// e.g.
// fmt.Println(string(DumpJSON(mystruct)))
//
func DumpJSON(i interface{}) []byte {
	var model []byte
	model, _ = json.MarshalIndent(i, "", "    ")
	return model
}
