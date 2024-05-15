package authentication

// APIKeyOperator
// Summary: This is structure which defines the APIkeyOperator model.
type APIKeyOperator struct {
	APIKey     string `json:"api_key"`
	OperatorID string `json:"operator_id"`
}

// APIKeyOperators
// Summary: This is structure which defines the slice of APIkeyOperator.
type APIKeyOperators []APIKeyOperator

// GetOperatorIds
// Summary: This is the function which returns the operator ids of the struct.
// output: ([]string) operator ids
func (ms APIKeyOperators) GetOperatorIds() []string {
	var operatorIDs = make([]string, len(ms))
	for i, m := range ms {
		operatorIDs[i] = m.OperatorID
	}
	return operatorIDs
}
