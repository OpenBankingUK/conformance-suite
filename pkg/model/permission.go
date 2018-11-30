package model

import (
	"encoding/json"
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

// GetPermissionFromName returns a permission if a matching permission name is found
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
	if err := json.Unmarshal(permissionStaticData, &permissions); err != nil {
		return err
	}
	return nil
}

// Permission Set Handling

// PermissionSet contains a collection of permission names with the intention
// of using this structure to determine whether specific permissions are included
// in the set.
type PermissionSet struct {
	name string
	set  map[string]bool
}

// NewPermissionSet create a new permission set of an array of string permission
// names
func NewPermissionSet(name string, strPermissions []string) *PermissionSet {
	var set PermissionSet
	set.name = name
	set.set = make(map[string]bool)
	set.AddPermissions(strPermissions)
	return &set
}

// GetName returns the name associated with this permission set
func (set *PermissionSet) GetName() string {
	return set.name
}

// SetName sets the name associated with this permission set
func (set *PermissionSet) SetName(s string) {
	set.name = s
}

// Add a string to a PermissionSet
func (set *PermissionSet) Add(s string) bool {
	_, found := set.set[s]
	set.set[s] = true
	return !found // return false if already existed
}

// AddPermissions - adds permission strings from a slice
func (set *PermissionSet) AddPermissions(ss []string) {
	for _, s := range ss {
		set.Add(s)
	}
}

// Get a permission from the PermissionSet
func (set *PermissionSet) Get(s string) bool {
	_, found := set.set[s]
	return found // true if already exists
}

// Remove a value from the PermissionSet
func (set *PermissionSet) Remove(s string) {
	delete(set.set, s)
}

// GetPermissions returns a string array of the permissions in a permissionSet
func (set *PermissionSet) GetPermissions() []string {
	var result []string
	for k := range set.set {
		result = append(result, k)
	}
	return result
}

// IsSubset determines if the permissionSet passed in as a paramter
// is a subset of the target PermissionSet
func (set *PermissionSet) IsSubset(sub *PermissionSet) bool {
	for key := range sub.set {
		found := set.Get(key)
		if !found {
			return false
		}
	}
	return true
}

// Union returns a new PermissionSet named "union" which is the union
// of the receiver and parameter permissionSets
func (set *PermissionSet) Union(u *PermissionSet) *PermissionSet {
	ps := NewPermissionSet("union", []string{})
	for k := range set.set {
		ps.Add(k)
	}
	for _, v := range u.GetPermissions() {
		ps.Add(v)
	}
	return ps
}
