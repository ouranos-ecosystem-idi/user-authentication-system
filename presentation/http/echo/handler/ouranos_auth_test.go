package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"authenticator-backend/domain/common"
	"authenticator-backend/domain/model/authentication"
	"authenticator-backend/presentation/http/echo/handler"
	f "authenticator-backend/test/fixtures"
	mocks "authenticator-backend/test/mock"
	"authenticator-backend/usecase/input"
	"authenticator-backend/usecase/output"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// /////////////////////////////////////////////////////////////////////////////////
// POST /auth/change テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1. 201: 正常系(英大文字、英小文字、数字、特殊文字を各1文字以上入力)
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_Change_Normal(tt *testing.T) {
	var method = "POST"
	var endPoint = "/auth/change"

	tests := []struct {
		name         string
		modifyInput  func(i *input.ChangePasswordParam)
		expectStatus int
	}{
		{
			name: "2-1. 201: 正常系(英大文字、英小文字、数字、特殊文字を各1文字以上入力)",
			modifyInput: func(i *input.ChangePasswordParam) {
				i.NewPassword = authentication.Password("xx@&1234Pass")
			},
			expectStatus: http.StatusCreated,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var newPassword authentication.Password = authentication.Password("xx@&1234Pass")

			input := input.ChangePasswordParam{
				UID:         f.Claims.UID,
				NewPassword: newPassword,
			}
			test.modifyInput(&input)
			inputJSON, _ := json.Marshal(input)
			q := make(url.Values)

			e := echo.New()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), strings.NewReader(string(inputJSON)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.SetPath(endPoint)

			operator := &f.Claims
			c.Set("operator", operator)

			authUsecase := new(mocks.IAuthUsecase)
			verifyUsecase := new(mocks.IVerifyUsecase)
			authHandler := handler.NewAuthHandler(
				authUsecase,
				verifyUsecase,
			)

			authUsecase.On("ChangePassword", input).Return(nil)
			err := authHandler.ChangePassword(c)
			if assert.NoError(t, err) {
				assert.Equal(t, test.expectStatus, rec.Code)
				authUsecase.AssertExpectations(t)
			}
		})
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// POST /auth/change テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-2. 400: バリデーションエラー：newPasswordの型がstring以外の場合
// [x] 1-9. 400: バリデーションエラー：newPasswordに特殊文字を1つも含まない場合
// [x] 2-2. 500: システムエラー：変更失敗
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_Change_Abnormal(tt *testing.T) {
	var method = "POST"
	var endPoint = "/auth/change"

	tests := []struct {
		name         string
		modifyInput  func(i *input.ChangePasswordParam)
		invalidInput any
		receive      error
		expectError  string
	}{
		{
			name: "1-2. 400: バリデーションエラー：newPasswordの型がstring以外の場合",
			invalidInput: struct {
				UID         int
				NewPassword authentication.Password `json:"newPassword"`
			}{
				1,
				"a",
			},
			receive:     nil,
			expectError: "code=400, message={[auth] BadRequest Validation failed, UID: Unmarshal type error: expected=string, got=number.",
		},
		{
			name: "1-9. 400: バリデーションエラー：newPasswordに特殊文字を1つも含まない場合",
			modifyInput: func(i *input.ChangePasswordParam) {
				i.NewPassword = authentication.Password("xx1234Pass")
			},
			receive:     nil,
			expectError: "code=400, message={[auth] BadRequest Validation failed, newPassword: must include at least one special character.",
		},
		{
			name:        "2-2. 500: システムエラー：変更失敗",
			modifyInput: func(i *input.ChangePasswordParam) {},
			receive:     fmt.Errorf("password change fail."),
			expectError: "code=500, message={[auth] InternalServerError Unexpected error occurred",
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var newPassword authentication.Password = authentication.Password("xx@&1234Pass")

			input := input.ChangePasswordParam{
				UID:         f.Claims.UID,
				NewPassword: newPassword,
			}
			var inputJSON []byte
			if test.invalidInput != nil {
				inputJSON, _ = json.Marshal(test.invalidInput)
			} else {
				test.modifyInput(&input)
				inputJSON, _ = json.Marshal(input)
			}

			q := make(url.Values)

			e := echo.New()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), strings.NewReader(string(inputJSON)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.SetPath(endPoint)

			operator := &f.Claims
			c.Set("operator", operator)

			authUsecase := new(mocks.IAuthUsecase)
			verifyUsecase := new(mocks.IVerifyUsecase)
			authHandler := handler.NewAuthHandler(
				authUsecase,
				verifyUsecase,
			)

			authUsecase.On("ChangePassword", input).Return(test.receive)
			err := authHandler.ChangePassword(c)
			if assert.Error(t, err) {
				assert.ErrorContains(t, err, test.expectError)
			}
		})
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// POST /auth/login 正常系テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1. 201: 正常系(ログイン成功)
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_Login_Normal(tt *testing.T) {
	var method = "POST"
	var endPoint = "/auth/login"
	var dataTarget = "operator"

	tests := []struct {
		name         string
		modifyInput  func(i *input.LoginParam)
		expectStatus int
	}{
		{
			name: "2-1. 201: 正常系(ログイン成功)",
			modifyInput: func(i *input.LoginParam) {
				i.OperatorAccountID = f.Email
				i.AccountPassword = "xx@&1234Pass"
			},
			expectStatus: http.StatusCreated,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				input := input.LoginParam{
					OperatorAccountID: f.Email,
					AccountPassword:   "xx@&1234Pass",
				}
				test.modifyInput(&input)
				inputJSON, _ := json.Marshal(input)
				loginModel := output.LoginResponse{}

				q := make(url.Values)
				q.Set("dataTarget", dataTarget)

				e := echo.New()
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), strings.NewReader(string(inputJSON)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath(endPoint)

				authUsecase := new(mocks.IAuthUsecase)
				verifyUsecase := new(mocks.IVerifyUsecase)
				authHandler := handler.NewAuthHandler(
					authUsecase,
					verifyUsecase,
				)
				authUsecase.On("Login", input).Return(loginModel, nil)

				err := authHandler.Login(c)
				if assert.NoError(t, err) {
					assert.Equal(t, test.expectStatus, rec.Code)
					authUsecase.AssertExpectations(t)
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// POST /auth/login テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1. 400: 異常系(バリデーションエラー：リクエストパラメータ変換不正)
// [x] 2-2. 400: 異常系(バリデーションエラー：フォーマット不正)
// [x] 2-3. 500: 異常系(システムエラー：ログイン失敗)
// [x] 2-4. 401: 異常系(認証エラー：ログイン失敗)
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_Login_Abnormal(tt *testing.T) {
	var method = "POST"
	var endPoint = "/auth/login"

	tests := []struct {
		name         string
		modifyInput  func(i *input.LoginParam)
		invalidInput any
		receive      error
		expectError  string
	}{
		{
			name: "2-1. 400: バリデーションエラー：リクエストパラメータ変換不正",
			invalidInput: struct {
				OperatorAccountID int
				AccountPassword   string
			}{
				1,
				"xx@&1234Pass",
			},
			receive:     nil,
			expectError: "code=400, message={[auth] BadRequest Validation failed, operatorAccountId: Unmarshal type error: expected=string, got=number.",
		},
		{
			name: "2-2. 400: バリデーションエラー：フォーマット不正",
			modifyInput: func(i *input.LoginParam) {
				i.OperatorAccountID = "xx1234Pass"
			},
			receive:     nil,
			expectError: "code=400, message={[auth] BadRequest Validation failed, operatorAccountId: must be a valid email address.",
		},
		{
			name:        "2-3. 500: システムエラー：ログイン失敗",
			modifyInput: func(i *input.LoginParam) {},
			receive:     fmt.Errorf("password change fail."),
			expectError: "code=500, message={[auth] InternalServerError Unexpected error occurred",
		},
		{
			name:        "2-4. 401: 認証エラー：ログイン失敗",
			modifyInput: func(i *input.LoginParam) {},
			receive:     common.NewCustomError(common.CustomErrorCode401, common.Err401InvalidCredentials, nil, common.HTTPErrorSourceAuth),
			expectError: "code=401, message={[auth] Unauthorized Invalid credentials id",
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				input := input.LoginParam{
					OperatorAccountID: f.Email,
					AccountPassword:   "xx@&1234Pass",
				}

				var inputJSON []byte
				if test.invalidInput != nil {
					inputJSON, _ = json.Marshal(test.invalidInput)
				} else {
					test.modifyInput(&input)
					inputJSON, _ = json.Marshal(input)
				}

				loginModel := output.LoginResponse{}

				q := make(url.Values)

				e := echo.New()
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), strings.NewReader(string(inputJSON)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath(endPoint)

				authUsecase := new(mocks.IAuthUsecase)
				verifyUsecase := new(mocks.IVerifyUsecase)
				authHandler := handler.NewAuthHandler(
					authUsecase,
					verifyUsecase,
				)
				authUsecase.On("Login", input).Return(loginModel, test.receive)

				err := authHandler.Login(c)
				if assert.Error(t, err) {
					assert.ErrorContains(t, err, test.expectError)
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// POST /auth/refresh 正常系テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1. 201: 正常系(リフレッシュトークン更新)
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_Refresh_Normal(tt *testing.T) {
	var method = "POST"
	var endPoint = "/auth/refresh"

	tests := []struct {
		name         string
		modifyInput  func(i *input.RefreshParam)
		expectStatus int
	}{
		{
			name: "2-1. 201: 正常系(リフレッシュトークン更新)",
			modifyInput: func(i *input.RefreshParam) {
				i.RefreshToken = f.Token
			},
			expectStatus: http.StatusCreated,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				input := input.RefreshParam{
					RefreshToken: f.Token,
				}
				test.modifyInput(&input)
				inputJSON, _ := json.Marshal(input)
				refreshModel := output.RefreshResponse{}

				q := make(url.Values)

				e := echo.New()
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), strings.NewReader(string(inputJSON)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath(endPoint)

				authUsecase := new(mocks.IAuthUsecase)
				verifyUsecase := new(mocks.IVerifyUsecase)
				authHandler := handler.NewAuthHandler(
					authUsecase,
					verifyUsecase,
				)
				authUsecase.On("Refresh", input).Return(refreshModel, nil)

				err := authHandler.Refresh(c)
				if assert.NoError(t, err) {
					assert.Equal(t, test.expectStatus, rec.Code)
					authUsecase.AssertExpectations(t)
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// POST /auth/refresh テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1. 400: 異常系(バリデーションエラー：リクエストパラメータ変換不正)
// [x] 2-2. 400: 異常系(バリデーションエラー：フォーマット不正)
// [x] 2-3. 500: 異常系(システムエラー：ログイン失敗)
// [x] 2-4. 401: 異常系(認証エラー：ログイン失敗)
// /////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////
// POST /auth/refresh 正常系テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1. 201: 正常系(リフレッシュトークン更新)
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_Refresh_Abnormal(tt *testing.T) {
	var method = "POST"
	var endPoint = "/auth/refresh"

	tests := []struct {
		name         string
		modifyInput  func(i *input.RefreshParam)
		invalidInput any
		receive      error
		expectError  string
	}{
		{
			name: "2-1. 400: バリデーションエラー：リクエストパラメータ変換不正",
			invalidInput: struct {
				RefreshToken int
			}{
				1,
			},
			receive:     nil,
			expectError: "code=400, message={[auth] BadRequest Validation failed, refreshToken: Unmarshal type error: expected=string, got=number.",
		},
		{
			name: "2-2. 400: バリデーションエラー：リクエストパラメータ変換不正",
			modifyInput: func(i *input.RefreshParam) {
				i.RefreshToken = ""
			},
			receive:     nil,
			expectError: "code=400, message={[auth] BadRequest Validation failed, refreshToken: cannot be blank.",
		},
		{
			name: "2-3. 401: 認証エラー：リフレッシュ失敗",
			modifyInput: func(i *input.RefreshParam) {
				i.RefreshToken = f.Token
			},
			receive:     common.NewCustomError(common.CustomErrorCode401, common.Err401InvalidCredentials, nil, common.HTTPErrorSourceAuth),
			expectError: "code=401, message={[auth] Unauthorized Invalid credentials id",
		},
		{
			name: "2-4. 500: システムエラー：リフレッシュ失敗",
			modifyInput: func(i *input.RefreshParam) {
				i.RefreshToken = f.Token
			},
			receive:     fmt.Errorf("RefreshToken get fail."),
			expectError: "code=500, message={[auth] InternalServerError Unexpected error occurred",
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				input := input.RefreshParam{
					RefreshToken: f.Token,
				}
				var inputJSON []byte
				if test.invalidInput != nil {
					inputJSON, _ = json.Marshal(test.invalidInput)
				} else {
					test.modifyInput(&input)
					inputJSON, _ = json.Marshal(input)
				}
				refreshModel := output.RefreshResponse{}

				q := make(url.Values)

				e := echo.New()
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), strings.NewReader(string(inputJSON)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath(endPoint)

				authUsecase := new(mocks.IAuthUsecase)
				verifyUsecase := new(mocks.IVerifyUsecase)
				authHandler := handler.NewAuthHandler(
					authUsecase,
					verifyUsecase,
				)
				authUsecase.On("Refresh", input).Return(refreshModel, test.receive)

				err := authHandler.Refresh(c)
				if assert.Error(t, err) {
					assert.ErrorContains(t, err, test.expectError)
				}
			},
		)
	}
}
