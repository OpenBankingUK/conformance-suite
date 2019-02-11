package model

// Registry - contains a list of component
type Registry struct {
	Entries map[string]interface{}
}

// NewRegistry -
func NewRegistry() Registry {
	entries := make(map[string]interface{})
	return Registry{Entries: entries}
}

// Add - adds a name to the registry
func (r *Registry) Add(name string, item interface{}) {
	r.Entries[name] = item
}

// Get - gets and item from the refistry
func (r *Registry) Get(name string) (interface{}, bool) {
	item, exists := r.Entries[name]
	return item, exists
}
