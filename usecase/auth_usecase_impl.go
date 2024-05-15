package usecase

import (
	"authenticator-backend/domain/common"
	"authenticator-backend/domain/repository"
	"authenticator-backend/extension/logger"
	"authenticator-backend/usecase/input"
	"authenticator-backend/usecase/output"
)

// authUsecase
// Summary: This is the structure which defines the usecase for the Auth.
type authUsecase struct {
	firebaseRepository repository.FirebaseRepository
}

// NewAuthUsecase
// Summary: This is the function which creates the Auth usecase.
// input: (repository.FirebaseRepository) r: firebase repository
// output: (IAuthUsecase) Auth usecase
func NewAuthUsecase(r repository.FirebaseRepository) IAuthUsecase {
	return &authUsecase{r}
}

// Login
// Summary: This is the function which logs in the operator.
// input: input(input.LoginParam): input parameter
// output: (output.LoginResponse) login response
// output: (error) error object
func (u authUsecase) Login(input input.LoginParam) (output.LoginResponse, error) {
	res, err := u.firebaseRepository.SignInWithPassword(input.OperatorAccountID, input.AccountPassword)
	if err != nil {
		logger.Set(nil).Errorf(err.Error())

		return output.LoginResponse{}, err
	}
	if res.AccessToken == "" || res.RefreshToken == "" {
		// when id/pass is invalid
		logger.Set(nil).Warnf(common.Err401InvalidCredentials)

		return output.LoginResponse{}, common.NewCustomError(common.CustomErrorCode401, common.Err401InvalidCredentials, nil, common.HTTPErrorSourceAuth)
	}

	return output.LoginResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

// Refresh
// Summary: This is the function which refreshes the token.
// input: input(input.RefreshParam): input parameter
// output: (output.RefreshResponse) refresh response
// output: (error) error object
func (u authUsecase) Refresh(input input.RefreshParam) (output.RefreshResponse, error) {
	token, err := u.firebaseRepository.RefreshToken(input.RefreshToken)
	if err != nil {
		logger.Set(nil).Errorf(err.Error())

		return output.RefreshResponse{}, err
	}
	if token == "" {
		// when refresh token is invalid
		logger.Set(nil).Warnf(common.Err401InvalidCredentials)

		return output.RefreshResponse{}, common.NewCustomError(common.CustomErrorCode401, common.Err401InvalidCredentials, nil, common.HTTPErrorSourceAuth)
	}
	return output.RefreshResponse{AccessToken: token}, nil
}

// ChangePassword
// Summary: This is the function which changes the password.
// input: input(input.ChangePasswordParam): input parameter
// output: (error) error object
func (u authUsecase) ChangePassword(input input.ChangePasswordParam) error {
	return u.firebaseRepository.ChangePassword(input.UID, input.NewPassword)
}
