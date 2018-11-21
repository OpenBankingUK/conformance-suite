package model

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

// Represents OB API specification.
// Fields are from the APIReference JSON-LD schema, see: https://schema.org/APIReference
// URL - URL of confluence specification file.
// SchemaVersion - URL of OpenAPI/Swagger specification file.
type specification struct {
	Identifier    string `json:"identifier,omitempty"`
	Name          string `json:"name,omitempty"`
	URL           string `json:"url,omitempty"`
	Version       string `json:"version,omitempty"`
	SchemaVersion string `json:"schemaVersion,omitempty"`
}

var specifications []specification

func init() {
	err := loadSpecifications()
	if err != nil {
		logrus.Error(err)
		os.Exit(1) // Abort if we can't read the config correctly
	}
}

// loadSpecifications - get specification data from json file
func loadSpecifications() error {
	rawjson, err := ioutil.ReadFile("../../config/specifications.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(rawjson, &specifications); err != nil {
		return err
	}

	return nil
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
