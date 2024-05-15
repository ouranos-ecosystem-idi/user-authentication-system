package authentication

// APIKey
// Summary: This is structure which defines the APIKey model.
type APIKey struct {
	ID              string
	APIKey          string
	ApplicationName string
	Attribute       ApplicationAttribute
}

// APIKeys
// Summary: This is structure which defines the slice of APIKey.
type APIKeys []APIKey

// ApplicationAttribute
// Summary: This is the type which defines the application attribute enum.
type ApplicationAttribute string

const (
	ApplicationAttributeDataSpace    ApplicationAttribute = "DataSpace"
	ApplicationAttributeApplication  ApplicationAttribute = "Application"
	ApplicationAttributeTraceability ApplicationAttribute = "Traceability"
)

// ContainsAPIKey
// Summary: This is the function which checks whether the API key exists in this struct slice.
// input: apiKey(string): API key
// output: (bool) true if the API key exists in this slice, false otherwise
func (ms APIKeys) ContainsAPIKey(apiKey string) bool {
	for _, m := range ms {
		if m.APIKey == apiKey {
			return true
		}
	}
	return false
}
