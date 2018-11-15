package model

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

// Permission holds endpoint permission data
type Permission struct {
	Permission          string   `json:"permission,omitempty"`
	Endpoints           []string `json:"endpoints,omitempty"`
	Default             bool     `json:"default,omitempty"`
	RequiredPermissions []string `json:"required_permissions,omitempty"`
	Optional            []string `json:"optional,omitempty"`
}

// EndpointConditionality - Store of endpoint conditionality
var permissions []Permission

func init() {
	err := loadPermissions()
	if err != nil {
		logrus.Error(err)
		os.Exit(1) // Abort if we can't read the config
	}
}

// GetPermissionsForEndpoint returns a list of permissions that are accepted by the specified endpoint
// no indication of whats mandatory/optional is given, you have to examine the individual permissions
// returned for that information
// if not entries are found a permission array with zero entries is returned
func GetPermissionsForEndpoint(endpoint string) []Permission {
	var endpointPermissions = []Permission{}
	for _, p := range permissions {
		for _, e := range p.Endpoints {
			if e == endpoint {
				endpointPermissions = append(endpointPermissions, p)
			}
		}
	}
	return endpointPermissions
}

// GetPermissionFromName returns a permission is a matching permission name is found
// or and empty permission if an entry with a matching name is not found
func GetPermissionFromName(name string) Permission {
	for _, p := range permissions {
		if name == p.Permission {
			return p
		}
	}
	return Permission{}
}

// loads permission data into modal permissions array structure
func loadPermissions() error {
	rawjson, _ := ioutil.ReadFile("../../pkg/model/permissions.json")
	err := json.Unmarshal(rawjson, &permissions)
	if err != nil {
		return err
	}
	return nil
}
