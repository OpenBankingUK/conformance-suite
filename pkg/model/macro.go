package model

import (
	"errors"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

var macroMap = map[string]interface{}{
	"instructionIdentificationID": instructionIdentificationID,
}

func ExecuteMacro(name string, params []string) (string, error) {
	if _, fnFound := macroMap[name]; !fnFound {
		return "", errors.New("macro not found")
	}

	f := reflect.ValueOf(macroMap[name])
	if len(params) != f.Type().NumIn() {
		return "", errors.New("The number of params is not adapted.")
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result := f.Call(in)
	if len(result) < 1 {
		return "", errors.New("unable to get result from macro")
	}
	return result[0].String(), nil
}

// instructionIdentificationID is a macro used in manifests
func instructionIdentificationID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
