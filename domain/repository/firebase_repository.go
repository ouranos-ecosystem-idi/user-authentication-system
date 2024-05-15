package repository

import "authenticator-backend/domain/model/authentication"

// FirebaseRepository
// Summary: This is interface which defines FirebaseRepository　functions.
//
//go:generate mockery --name FirebaseRepository --output ../../test/mock --case underscore
type FirebaseRepository interface {
	SignInWithPassword(email string, password string) (authentication.LoginResult, error)
	VerifyIDToken(idToken string) (authentication.Claims, error)
	RefreshToken(refreshToken string) (string, error)
	ChangePassword(uid string, newPassword authentication.Password) error
}
