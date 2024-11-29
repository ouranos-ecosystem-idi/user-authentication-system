package repository_test

import (
	"authenticator-backend/domain/model/authentication"
	"authenticator-backend/infrastructure/firebase/repository"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
)

// /////////////////////////////////////////////////////////////////////////////////
// Firebase SignInWithPassword テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1: 正常系：正常返却の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Firebase_SignInWithPassword(tt *testing.T) {

	tests := []struct {
		name         string
		inputIdPPath string
		inputSecPath string
		receiveBody  string
		expect       authentication.LoginResult
	}{
		{
			name:         "1-1: 正常系",
			inputIdPPath: "identitytoolkit.googleapis.com/v1/accounts:signInWithPassword",
			inputSecPath: "securetoken.googleapis.com/v1/token",
			receiveBody: `{
				"kind": "kind",
				"localId": "localId",
				"email": "aaa@aaa.com",
				"displayName": "displayName",
				"idToken": "idToken",
				"registered": true,
				"refreshToken": "refreshToken",
				"expiresIn": "3600"
			}`,
			expect: authentication.LoginResult{
				AccessToken:  "idToken",
				RefreshToken: "refreshToken",
			},
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if strings.HasSuffix(r.URL.Path, test.inputIdPPath) {
						code, err := w.Write([]byte(test.receiveBody))
						if err != nil {
							w.WriteHeader(code)
						}
					} else {
						w.WriteHeader(http.StatusBadRequest)
						code, err := w.Write([]byte("Bad Request"))
						if err != nil {
							w.WriteHeader(code)
						}
					}
				})
				ts := httptest.NewServer(handler)
				defer ts.Close()

				r := repository.NewFirebase(nil, fmt.Sprintf("%s/%s", ts.URL, test.inputIdPPath), "aaa", "apikey", fmt.Sprintf("%s/%s", ts.URL, test.inputSecPath))
				actual, err := r.SignInWithPassword("aaa@aaa.com", "password")
				if assert.NoError(t, err) {
					assert.Equal(t, test.expect.AccessToken, actual.AccessToken)
					assert.Equal(t, test.expect.RefreshToken, actual.RefreshToken)
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Firebase SignInWithPassword テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1: 異常系：400の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Firebase_SignInWithPassword_Abnormal(tt *testing.T) {

	tests := []struct {
		name         string
		inputIdPPath string
		inputSecPath string
	}{
		{
			name:         "2-1: 異常系：503の場合",
			inputIdPPath: "identitytoolkit.googleapis.com/v1/accounts:signInWithPassword",
			inputSecPath: "securetoken.googleapis.com/v1/token",
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
					code, err := w.Write([]byte("Bad Request"))
					if err != nil {
						w.WriteHeader(code)
					}
				})
				ts := httptest.NewServer(handler)
				defer ts.Close()

				r := repository.NewFirebase(nil, fmt.Sprintf("%s/%s", ts.URL, test.inputIdPPath), "aaa", "apikey", fmt.Sprintf("%s/%s", ts.URL, test.inputSecPath))
				_, err := r.SignInWithPassword("aaa@aaa.com", "password")
				assert.Error(t, err)
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Firebase RefreshToken テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1: 正常系：正常返却の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Firebase_RefreshToken(tt *testing.T) {

	tests := []struct {
		name         string
		inputIdPPath string
		inputSecPath string
		receiveBody  string
		expect       string
	}{
		{
			name:         "1-1: 正常系",
			inputIdPPath: "identitytoolkit.googleapis.com/v1/accounts:signInWithPassword",
			inputSecPath: "securetoken.googleapis.com/v1/token",
			receiveBody: `{
				"access_token": "accessToken",
				"expires_in": "3600",
				"token_type": "tokenType",
				"refresh_token": "refreshToken",
				"id_token": "idToken",
				"user_id": "userId",
				"project_id": "projectId"
			}`,
			expect: "accessToken",
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if strings.HasSuffix(r.URL.Path, test.inputSecPath) {
						code, err := w.Write([]byte(test.receiveBody))
						if err != nil {
							w.WriteHeader(code)
						}
					} else {
						w.WriteHeader(http.StatusBadRequest)
						code, err := w.Write([]byte("Bad Request"))
						if err != nil {
							w.WriteHeader(code)
						}
					}
				})
				ts := httptest.NewServer(handler)
				defer ts.Close()
				r := repository.NewFirebase(nil, fmt.Sprintf("%s/%s", ts.URL, test.inputIdPPath), "aaa", "apikey", fmt.Sprintf("%s/%s", ts.URL, test.inputSecPath))
				actual, err := r.RefreshToken("token")
				if assert.NoError(t, err) {
					assert.Equal(t, test.expect, actual)
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Firebase RefreshToken テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1: 異常系：400の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Firebase_RefreshToken_Abnormal(tt *testing.T) {

	tests := []struct {
		name         string
		inputIdPPath string
		inputSecPath string
	}{
		{
			name:         "2-1: 異常系：400の場合",
			inputIdPPath: "identitytoolkit.googleapis.com/v1/accounts:signInWithPassword",
			inputSecPath: "securetoken.googleapis.com/v1/token",
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
					code, err := w.Write([]byte("Bad Request"))
					if err != nil {
						w.WriteHeader(code)
					}
				})
				ts := httptest.NewServer(handler)
				defer ts.Close()

				r := repository.NewFirebase(nil, fmt.Sprintf("%s/%s", ts.URL, test.inputIdPPath), "aaa", "apikey", fmt.Sprintf("%s/%s", ts.URL, test.inputSecPath))
				_, err := r.RefreshToken("token")
				assert.Error(t, err)
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Firebase VerifyIDToken テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1: 正常系：正常返却の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Firebase_VerifyIDToken(tt *testing.T) {

	tests := []struct {
		name           string
		inputProjectID string
		inputClaim     func(projectID string, operatorID string) string
		receiveBody    string
		expect         string
	}{
		{
			name:           "1-1: 正常系",
			inputProjectID: "local",
			inputClaim: func(projectID string, operatorID string) string {
				// ペイロードを設定
				token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
					"iss":         fmt.Sprintf("https://securetoken.google.com/%s", projectID),
					"aud":         projectID,
					"sub":         "test",
					"iat":         time.Now().Unix(),
					"exp":         time.Now().Add(time.Hour * 1).Unix(),
					"operator_id": operatorID,
				})
				key, _ := rsa.GenerateKey(rand.Reader, 256*8)
				idToken, _ := token.SignedString(key)
				return idToken
			},
			receiveBody: `{
				"users":[
					{"localId": "test"}
				]
			}`,
			expect: "b39e6248-c888-56ca-d9d0-89de1b1adc8e",
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if strings.HasSuffix(r.URL.Path, fmt.Sprintf("%s/accounts:lookup", test.inputProjectID)) {
						code, err := w.Write([]byte(test.receiveBody))
						if err != nil {
							w.WriteHeader(code)
						}
					} else {
						w.WriteHeader(http.StatusBadRequest)
						code, err := w.Write([]byte("Bad Request"))
						if err != nil {
							w.WriteHeader(code)
						}
					}
				})
				ts := httptest.NewServer(handler)
				defer ts.Close()
				conf := &firebase.Config{ProjectID: test.inputProjectID}
				os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", strings.Replace(ts.URL, "http://", "", 1))
				ctx := context.Background()
				app, _ := firebase.NewApp(ctx, conf, option.WithoutAuthentication())
				authCli, _ := app.Auth(ctx)
				r := repository.NewFirebase(authCli, ts.URL, "aaa", "apikey", ts.URL)
				actual, err := r.VerifyIDToken(test.inputClaim(test.inputProjectID, test.expect))
				if assert.NoError(t, err) {
					assert.Equal(t, test.expect, actual.OperatorID)
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Firebase VerifyIDToken テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1: 異常系：対象ユーザーなしの場合
// [x] 2-2: 異常系：Claim不正の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Firebase_VerifyIDToken_Abnormal(tt *testing.T) {

	tests := []struct {
		name           string
		inputProjectID string
		inputClaim     func(projectID string, operatorID string) string
		receiveBody    string
		expect         error
	}{
		{
			name:           "2-1: 異常系：対象ユーザーなしの場合",
			inputProjectID: "local",
			inputClaim: func(projectID string, operatorID string) string {
				// ペイロードを設定
				token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
					"iss":         fmt.Sprintf("https://securetoken.google.com/%s", projectID),
					"aud":         projectID,
					"sub":         "test",
					"iat":         time.Now().Unix(),
					"exp":         time.Now().Add(time.Hour * 1).Unix(),
					"operator_id": operatorID,
				})
				key, _ := rsa.GenerateKey(rand.Reader, 256*8)
				idToken, _ := token.SignedString(key)
				return idToken
			},
			receiveBody: `{
				"users":[]
			}`,
			expect: fmt.Errorf("no user exists with the uid: \"test\""),
		},
		{
			name:           "2-2: 異常系：Claim不正の場合",
			inputProjectID: "local",
			inputClaim: func(projectID string, operatorID string) string {
				// ペイロードを設定
				token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
					"iss": fmt.Sprintf("https://securetoken.google.com/%s", projectID),
					"aud": projectID,
					"sub": "test",
					"iat": time.Now().Unix(),
					"exp": time.Now().Add(time.Hour * 1).Unix(),
				})
				key, _ := rsa.GenerateKey(rand.Reader, 256*8)
				idToken, _ := token.SignedString(key)
				return idToken
			},
			receiveBody: `{
				"users":[
					{"localId": "test"}
				]
			}`,
			expect: fmt.Errorf("token does not contain 'operator_id' in claims"),
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if strings.HasSuffix(r.URL.Path, fmt.Sprintf("%s/accounts:lookup", test.inputProjectID)) {
						code, err := w.Write([]byte(test.receiveBody))
						if err != nil {
							w.WriteHeader(code)
						}
					} else {
						w.WriteHeader(http.StatusBadRequest)
						code, err := w.Write([]byte("Bad Request"))
						if err != nil {
							w.WriteHeader(code)
						}
					}
				})
				ts := httptest.NewServer(handler)
				defer ts.Close()
				conf := &firebase.Config{ProjectID: test.inputProjectID}
				os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", strings.Replace(ts.URL, "http://", "", 1))
				ctx := context.Background()
				app, _ := firebase.NewApp(ctx, conf, option.WithoutAuthentication())
				authCli, _ := app.Auth(ctx)
				r := repository.NewFirebase(authCli, ts.URL, "aaa", "apikey", ts.URL)
				_, err := r.VerifyIDToken(test.inputClaim(test.inputProjectID, ""))
				if assert.Error(t, err) {
					assert.Equal(t, test.expect.Error(), err.Error())
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Firebase ChangePassword テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1: 正常系：正常返却の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Firebase_ChangePassword(tt *testing.T) {

	tests := []struct {
		name           string
		inputProjectID string
		receiveBody    string
		expect         string
	}{
		{
			name:           "1-1: 正常系",
			inputProjectID: "local",
			receiveBody: `{
				"users":[
					{"localId": "test"}
				]
			}`,
			expect: "b39e6248-c888-56ca-d9d0-89de1b1adc8e",
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if strings.HasSuffix(r.URL.Path, fmt.Sprintf("%s/accounts:lookup", test.inputProjectID)) {
						code, err := w.Write([]byte(test.receiveBody))
						if err != nil {
							w.WriteHeader(code)
						}
					} else if strings.HasSuffix(r.URL.Path, fmt.Sprintf("%s/accounts:update", test.inputProjectID)) {
						code, err := w.Write([]byte(""))
						if err != nil {
							w.WriteHeader(code)
						}
					} else {
						w.WriteHeader(http.StatusBadRequest)
						code, err := w.Write([]byte("Bad Request"))
						if err != nil {
							w.WriteHeader(code)
						}
					}
				})
				ts := httptest.NewServer(handler)
				defer ts.Close()
				conf := &firebase.Config{ProjectID: test.inputProjectID}
				os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", strings.Replace(ts.URL, "http://", "", 1))
				ctx := context.Background()
				app, _ := firebase.NewApp(ctx, conf, option.WithoutAuthentication())
				authCli, _ := app.Auth(ctx)
				r := repository.NewFirebase(authCli, ts.URL, "aaa", "apikey", ts.URL)
				err := r.ChangePassword("test", "newpass")
				assert.NoError(t, err)
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Firebase ChangePassword テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1: 異常系：400の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Firebase_ChangePassword_Abnormal(tt *testing.T) {

	tests := []struct {
		name           string
		inputProjectID string
	}{
		{
			name:           "2-1: 異常系：400の場合",
			inputProjectID: "local",
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
					code, err := w.Write([]byte("Bad Request"))
					if err != nil {
						w.WriteHeader(code)
					}
				})
				ts := httptest.NewServer(handler)
				defer ts.Close()
				conf := &firebase.Config{ProjectID: test.inputProjectID}
				os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", strings.Replace(ts.URL, "http://", "", 1))
				ctx := context.Background()
				app, _ := firebase.NewApp(ctx, conf, option.WithoutAuthentication())
				authCli, _ := app.Auth(ctx)
				r := repository.NewFirebase(authCli, ts.URL, "aaa", "apikey", ts.URL)
				err := r.ChangePassword("test", "newpass")
				assert.Error(t, err)
			},
		)
	}
}
