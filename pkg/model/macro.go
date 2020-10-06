package model

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
)

var macroMap = map[string]interface{}{
	"instructionIdentificationID": instructionIdentificationID,
	"nextDayDateTime":             nextDayDateTime,
	"nextDayDateTimeHour":         nextDayDateTimeHour,
}

// AddMacro inserts the provided macro in the map where they are held.
// It is not expected to be called concurrently.
func AddMacro(name string, macro interface{}) {
	macroMap[name] = macro
}

// ExecuteMacro calls a macro by `name`, with parameters to be passed using `params`. `params` is a collection of strings
// that get passed as is. Type assertions will need be performed in the macro implementation.
func ExecuteMacro(name string, params []string) (string, error) {
	macro, found := macroMap[name]
	if !found {
		return "", errors.New("macro not found")
	}

	f := reflect.ValueOf(macro)
	if len(params) != f.Type().NumIn() {
		return "", errors.New("the number of params is not adapted")
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

func nextDayDateTime(format string) string {
	// In the tests which use generated times, there must be no assertion
	// on the timestamp's actual value, (e.g. checking if time == 2022-01-01T12:00:00Z).
	nextDay := time.Now().UTC().Add(24 * time.Hour)
	return nextDay.Format(format)
}

func nextDayDateTimeHour(format string) string {
	nextDay := roundDownToHour(time.Now().UTC().Add(24 * time.Hour))
	return nextDay.Format(format)
}

func roundDownToHour(t time.Time) time.Time {
	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		0, 0, 0,
		t.Location(),
	)
}
