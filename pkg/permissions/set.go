package permissions

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

// add a string to a PermissionSet
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

// IsSubset determines if the permissionSet passed in as a parameter
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

// Equal determines if the permissionSet passed in as a parameter
// has same set of strings as the target PermissionSet
func (set *PermissionSet) Equal(other *PermissionSet) bool {
	equal := set.IsSubset(other) && other.IsSubset(set)
	return equal
}

// add returns a new PermissionSet named "union" which is the union
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

// Intersection returns a new PermissionSet named "intersection" which is the intersection
// of the receiver and parameter permissionSets
func (set *PermissionSet) Intersection(other *PermissionSet) *PermissionSet {
	ps := NewPermissionSet("intersection", []string{})
	for key := range set.set {
		found := other.Get(key)
		if found {
			ps.Add(key)
		}
	}
	return ps
}
