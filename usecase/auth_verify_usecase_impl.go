package usecase

import (
	"authenticator-backend/domain/model/authentication"
	"authenticator-backend/domain/repository"
	"authenticator-backend/extension/logger"
	"authenticator-backend/usecase/input"
	"authenticator-backend/usecase/output"
)

// verifyUsecase
// Summary: This is the structure which defines the verify usecase.
type verifyUsecase struct {
	firebaseRepository repository.FirebaseRepository
	authRepository     repository.AuthRepository
}

// NewVerifyUsecase
// Summary: This is the function which creates the verify usecase.
// input: f(repository.FirebaseRepository) firebase repository
// input: a(repository.AuthRepository) auth repository
// output: (IVerifyUsecase) verify usecase
func NewVerifyUsecase(f repository.FirebaseRepository, a repository.AuthRepository) IVerifyUsecase {
	return &verifyUsecase{f, a}
}

// TokenIntrospection
// Summary: This is the function which verifies the token.
// input: input(input.VerifyTokenParam) input parameters
// output: (output.VerifyTokenResponse) output response
// output: (error) error object
func (u verifyUsecase) TokenIntrospection(input input.VerifyTokenParam) (output.VerifyTokenResponse, error) {
	claims, err := u.firebaseRepository.VerifyIDToken(input.IDToken)
	if err != nil {
		logger.Set(nil).Errorf(err.Error())

		return output.VerifyTokenResponse{}, err
	}
	return output.VerifyTokenResponse{OperatorID: &claims.OperatorID}, nil
}

// IDToken
// Summary: This is the function which verifies the ID token.
// input: input(input.VerifyIDTokenParam) input parameters
// output: (authentication.Claims) claims
// output: (error) error object
func (u verifyUsecase) IDToken(input input.VerifyIDTokenParam) (authentication.Claims, error) {
	claims, err := u.firebaseRepository.VerifyIDToken(input.IDToken)
	if err != nil {
		logger.Set(nil).Warnf(err.Error())

		return authentication.Claims{}, err
	}
	return claims, nil
}

// ApiKey
// Summary: This is the function which verifies the API key.
// input: input(input.VerifyApiKeyParam) input parameters
// output: (output.VerifyApiKeyResponse) output response
func (u verifyUsecase) ApiKey(input input.VerifyAPIKeyParam) output.VerifyApiKeyResponse {
	output := output.VerifyApiKeyResponse{
		IsAPIKeyValid:    false,
		IsIPAddressValid: false,
	}

	// 1. Check the validity of the APIKEY
	apikeys, err := u.authRepository.ListAPIKeys(repository.APIKeysParam{})
	if err != nil {
		logger.Set(nil).Warnf(err.Error())

		return output
	}
	if apikeys.ContainsAPIKey(input.APIKey) {
		output.IsAPIKeyValid = true
	}

	// 2. Check the combination of APIKEY and IP address
	cidrs, err := u.authRepository.ListCidrs(repository.APIKeyCidrsParam{
		APIKey: &input.APIKey,
	})
	if err != nil {
		logger.Set(nil).Warnf(err.Error())

		return output
	}
	if cidrs.Contains(input.IPAddress) {
		output.IsIPAddressValid = true
	}

	return output
}
