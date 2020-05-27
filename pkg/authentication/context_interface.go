package authentication

// ContextInterface - avoid cycling dependency to `model.Context`.
type ContextInterface interface {
	// GetString get the string value associated with key
	GetString(key string) (string, error)
	// Get the key form the Context map - currently assumes value converts easily to a string!
	Get(key string) (interface{}, bool)
	GetStringSlice(key string) ([]string, error)
}
