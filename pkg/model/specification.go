package model

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"

	validator "gopkg.in/go-playground/validator.v9"
)

// Represents OB API specification.
// Fields are from the APIReference JSON-LD schema, see: https://schema.org/APIReference
// URL - URL of confluence specification file.
// SchemaVersion - URL of OpenAPI/Swagger specification file.
type specification struct {
	Identifier    string `json:"identifier,omitempty" validate:"required"`
	Name          string `json:"name,omitempty" validate:"required"`
	URL           string `json:"url,omitempty" validate:"required,url"`
	Version       string `json:"version,omitempty" validate:"required"`
	SchemaVersion string `json:"schemaVersion,omitempty" validate:"required,url"`
}

var specifications []specification

// init - load and validate specification data from json file
func init() {
	err := loadDefaultSpecifications()
	if err != nil {
		logrus.Error(err)
		os.Exit(1) // Abort if we can't read the config correctly
	}
}

// loadDefaultSpecifications - load and validate specification data from json file
func loadDefaultSpecifications() error {
	file := "./config/specifications.json"
	if flag.Lookup("test.v") != nil {
		file = "../../config/specifications.json" // different path when running tests
	}

	rawjson, err := ioutil.ReadFile(file)
	if err == nil {
		err = loadSpecifications(rawjson)
	}
	return err
}

// loadDefaultSpecifications - load and validate specification data from json
func loadSpecifications(rawjson []byte) error {
	err := json.Unmarshal(rawjson, &specifications)
	if err == nil {
		validate := validator.New()
		for _, specConfig := range specifications {
			err = validate.Struct(specConfig)
			if err != nil {
				break // sufficient to report validation errors one at a time
			}
		}
	}
	return err
}

// SpecificationIdentifierFromSchemaVersion - returns specification identifier
// for given schema version URL, or nil when there is no match.
func SpecificationIdentifierFromSchemaVersion(schemaVersion string) (string, error) {
	for _, specification := range specifications {
		if specification.SchemaVersion == schemaVersion {
			return specification.Identifier, nil
		}
	}
	return "", errors.New("No specification found for schema version: " + schemaVersion)
}
