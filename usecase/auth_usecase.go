package usecase

import (
	"authenticator-backend/usecase/input"
	"authenticator-backend/usecase/output"
)

// IAuthUsecase
// Summary: This is interface which defines IAuthUsecase
//
//go:generate mockery --name IAuthUsecase --output ../test/mock --case underscore
type IAuthUsecase interface {
	Login(input input.LoginParam) (output.LoginResponse, error)
	Refresh(input input.RefreshParam) (output.RefreshResponse, error)
	ChangePassword(input input.ChangePasswordParam) error
}
