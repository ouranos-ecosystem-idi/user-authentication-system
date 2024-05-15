package datastore

import (
	"authenticator-backend/domain/model/authentication"
	"authenticator-backend/domain/repository"
	"authenticator-backend/extension/logger"

	"gorm.io/gorm"
)

// authRepository
// Summary:This is structure which defines the repository for the authentication.
type authRepository struct {
	db *gorm.DB
}

// NewAuthRepository
// Summary: This is the function which creates the authentication repository.
// input: db(*gorm.DB): gorm db
// output: (authRepository) authentication repository
func NewAuthRepository(
	db *gorm.DB,
) repository.AuthRepository {
	return &authRepository{db}
}

// ListAPIKeys
// Summary: This is the function which lists the api keys.
// input: param(APIKeysParam): apikeys param
// output: (authentication.APIKeys) apikeys
// output: (error) error object
func (r *authRepository) ListAPIKeys(param repository.APIKeysParam) (authentication.APIKeys, error) {
	var apikeys []authentication.APIKey

	query := r.db.Table("api_keys")
	if param.Attributes != nil {
		query = query.Where("application_attribute in (?)", param.Attributes)
	}
	if err := query.Find(&apikeys).Error; err != nil {
		logger.Set(nil).Errorf(err.Error())

		return nil, err
	}
	return apikeys, nil
}

// ListAPIKeyOperators
// Summary: This is the function which lists the api key operators.
// input: param(APIKeyOperatorsParam): apikey operators param
// output: (authentication.APIKeyOperators) apikey operators
// output: (error) error object
func (r *authRepository) ListAPIKeyOperators(param repository.APIKeyOperatorsParam) (authentication.APIKeyOperators, error) {
	var apikeyOperators authentication.APIKeyOperators

	query := r.db.Table("apikey_operators")
	if param.APIKey != nil {
		query = query.Where("api_key = ?", *param.APIKey)
	}
	if err := query.Find(&apikeyOperators).Error; err != nil {
		logger.Set(nil).Errorf(err.Error())

		return nil, err
	}
	return apikeyOperators, nil
}

// ListCidrs
// Summary: This is the function which lists the cidrs.
// input: param(APIKeyCidrsParam): apikey cidrs param
// output: (authentication.Cidrs) cidrs
// output: (error) error object
func (r *authRepository) ListCidrs(param repository.APIKeyCidrsParam) (authentication.Cidrs, error) {
	var cidrs authentication.Cidrs

	query := r.db.Table("cidrs")
	if param.APIKey != nil {
		query = query.Where("api_key = ?", *param.APIKey)
	}
	if err := query.Find(&cidrs).Error; err != nil {
		logger.Set(nil).Errorf(err.Error())

		return nil, err
	}
	return cidrs, nil

}
