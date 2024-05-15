package repository

import "authenticator-backend/domain/model/authentication"

// AuthRepository
// Summary: This is interface which defines the functions for the authentication repository.
//
//go:generate mockery --name AuthRepository --output ../../test/mock --case underscore
type AuthRepository interface {
	ListAPIKeys(param APIKeysParam) (authentication.APIKeys, error)
	ListAPIKeyOperators(param APIKeyOperatorsParam) (authentication.APIKeyOperators, error)
	ListCidrs(param APIKeyCidrsParam) (authentication.Cidrs, error)
}

// APIKeysParam
// Summary: This is the structure which defines the parameters for the ListAPIKeys Method.
type APIKeysParam struct {
	Attributes []authentication.ApplicationAttribute
}

// APIKeyOperatorsParam
// Summary: This is the structure which defines the parameters for the ListAPIKeyOperators Method.
type APIKeyOperatorsParam struct {
	APIKey *string
}

// APIKeyCidrsParam
// Summary: This is the structure which defines the parameters for the ListCidrs Method.
type APIKeyCidrsParam struct {
	APIKey *string
}
