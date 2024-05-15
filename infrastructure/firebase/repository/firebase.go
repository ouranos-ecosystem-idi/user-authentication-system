package repository

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"authenticator-backend/domain/common"
	"authenticator-backend/domain/model/authentication"
	"authenticator-backend/extension/logger"
	"authenticator-backend/infrastructure/firebase/entity"

	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
)

const errorCodeUnavailable = "UNAVAILABLE"

// firebaseRepository
// Summary: This struct is the repository for the firebase.
type firebaseRepository struct {
	cli                   *auth.Client
	signInWithPasswordURL string
	idpApikey             string
	secureTokenApiKey     string
	secureTokenApi        string
}

// NewFirebase
// Summary: This is the function which creates the firebase repository.
// input: cli(*auth.Client) auth client
// input: signInWithPasswordURL(string) sign in with password URL
// input: idpApikey(string) idp api key
// input: secureTokenApiKey(string) secure token api key
// input: secureTokenApi(string) secure token api
// output: (firebaseRepository) firebase repository
func NewFirebase(
	cli *auth.Client,
	signInWithPasswordURL string,
	idpApikey string,
	secureTokenApiKey string,
	secureTokenApi string,
) firebaseRepository {
	return firebaseRepository{
		cli,
		signInWithPasswordURL,
		idpApikey,
		secureTokenApiKey,
		secureTokenApi,
	}
}

// SignInWithPassword
// Summary: This is the function which signs in with email and password.
// input: email(string) email
// input: password(string) password
// output: (authentication.LoginResult) login result
// output: (error) error object
func (r firebaseRepository) SignInWithPassword(email string, password string) (authentication.LoginResult, error) {
	reqBody := map[string]interface{}{
		"email":             email,
		"password":          password,
		"returnSecureToken": true,
	}
	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		logger.Set(nil).Errorf(err.Error())

		return authentication.LoginResult{}, err
	}

	request, err := http.NewRequest(http.MethodPost, r.signInWithPasswordURL, strings.NewReader(string(reqBodyJson)))
	if err != nil {
		logger.Set(nil).Errorf(err.Error())

		return authentication.LoginResult{}, err
	}
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	values := url.Values{}
	values.Add("key", r.idpApikey)
	request.URL.RawQuery = values.Encode()

	client := new(http.Client)

	response, err := client.Do(request)
	if err != nil {
		logger.Set(nil).Errorf(err.Error())

		return authentication.LoginResult{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Set(nil).Errorf(err.Error())

		return authentication.LoginResult{}, err
	}
	defer response.Body.Close()

	var loginResponse entity.LoginResponse
	err = json.Unmarshal(body, &loginResponse)
	if err != nil {
		logger.Set(nil).Errorf(err.Error())

		return authentication.LoginResult{}, err
	}

	loginOutput := authentication.LoginResult{
		AccessToken:  loginResponse.IDToken,
		RefreshToken: loginResponse.RefreshToken,
	}

	return loginOutput, nil
}

// RefreshToken
// Summary: This is the function which refreshes the token.
// input: refreshToken(string) refresh token
// output: (string) access token
// output: (error) error object
func (r firebaseRepository) RefreshToken(refreshToken string) (string, error) {
	param := url.Values{}
	param.Add("key", r.idpApikey)
	requestURL, err := url.Parse(r.secureTokenApi)
	if err != nil {
		logger.Set(nil).Errorf(err.Error())

		return "", err
	}
	requestURL.RawQuery = param.Encode()

	formData := url.Values{}
	formData.Add("key", r.secureTokenApiKey)
	formData.Add("grant_type", "refresh_token")
	formData.Add("refresh_token", refreshToken)

	request, err := http.NewRequest(http.MethodPost, requestURL.String(), strings.NewReader(formData.Encode()))
	if err != nil {
		logger.Set(nil).Errorf(err.Error())

		return "", err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := new(http.Client)

	response, err := client.Do(request)
	if err != nil {
		logger.Set(nil).Errorf(err.Error())

		return "", err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Set(nil).Errorf(err.Error())

		return "", err
	}
	defer response.Body.Close()

	var refreshResponse entity.RefreshResponse
	err = json.Unmarshal(body, &refreshResponse)
	if err != nil {
		logger.Set(nil).Errorf(err.Error())

		return "", err
	}

	return refreshResponse.AccessToken, nil
}

// VerifyIDToken
// Summary: This is the function which verifies the ID token.
// input: idToken(string) id token
// output: (authentication.Claims) claims
// output: (error) error object
func (r firebaseRepository) VerifyIDToken(idToken string) (authentication.Claims, error) {
	ctx := context.Background()
	token, err := r.cli.VerifyIDToken(ctx, idToken)
	if err != nil {
		if errorCode, ok := extractErrorCode(err); ok {
			if errorCode == errorCodeUnavailable {
				logger.Set(nil).Errorf(err.Error())

				return authentication.Claims{}, common.NewCustomError(common.CustomErrorCode503, common.Err503OuterService, nil, common.HTTPErrorSourceAuth)
			}
		}
		logger.Set(nil).Errorf(err.Error())

		return authentication.Claims{}, err
	}
	claims, err := authentication.NewClaims(token)
	if err != nil {
		logger.Set(nil).Errorf(err.Error())

		return authentication.Claims{}, err
	}
	return claims, nil
}

// ChangePassword
// Summary: This is the function which changes the password.
// input: uid(string) firebase UID
// input: newPassword(authentication.Password) new password
// output: (error) error object
func (r firebaseRepository) ChangePassword(uid string, newPassword authentication.Password) error {
	ctx := context.Background()

	_, err := r.cli.UpdateUser(ctx, uid, (&auth.UserToUpdate{}).Password(newPassword.ToString()))
	if err != nil {
		logger.Set(nil).Errorf(err.Error())

		return err
	}
	return nil
}

// extractErrorCode
// Summary: This is the function which extracts the error code.
// input: err(error) error object
// output: (string) error code
// output: (bool) true if the error code is extracted, false otherwise
func extractErrorCode(err error) (string, bool) {
	val := reflect.ValueOf(err)
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	if val.Kind() == reflect.Struct {
		errorCodeField := val.FieldByName("ErrorCode")
		if errorCodeField.IsValid() && errorCodeField.Kind() == reflect.String {
			return errorCodeField.String(), true
		}
	}
	return "", false
}
