package usecase

import (
	"authenticator-backend/domain/model/authentication"
	"authenticator-backend/usecase/input"
	"authenticator-backend/usecase/output"
)

// IVerifyUsecase
// Summary: This is interface which defines IVerifyUsecase
//
//go:generate mockery --name IVerifyUsecase --output ../test/mock --case underscore
type IVerifyUsecase interface {
	TokenIntrospection(input input.VerifyTokenParam) (output.VerifyTokenResponse, error)
	IDToken(input input.VerifyIDTokenParam) (authentication.Claims, error)
	ApiKey(input input.VerifyAPIKeyParam) output.VerifyApiKeyResponse
}
