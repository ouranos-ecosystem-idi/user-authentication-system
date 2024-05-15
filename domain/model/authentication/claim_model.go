package authentication

import (
	"fmt"

	"firebase.google.com/go/v4/auth"
)

// Claims
// Summary: This is structure which defines the claims model.
type Claims struct {
	OperatorID string `json:"operator_id"`
	auth.Token
}

// NewClaims
// Summary: This is the function which creates the Claims model from the token.
// input: token(*auth.Token): token
// output: (Claims) Claims model
// output: (error) error object
func NewClaims(token *auth.Token) (Claims, error) {
	operatorID, ok := token.Claims["operator_id"].(string)
	if !ok {
		return Claims{}, fmt.Errorf("token does not contain 'operator_id' in claims")
	}
	return Claims{
		OperatorID: operatorID,
		Token:      *token,
	}, nil
}
