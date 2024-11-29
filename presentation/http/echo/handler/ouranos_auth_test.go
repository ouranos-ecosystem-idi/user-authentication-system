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
		inputFunc    func() input.ChangePasswordParam
		expectStatus int
	}{
		{
			name: "2-1. 201: 正常系(英大文字、英小文字、数字、特殊文字を各1文字以上入力)",
			inputFunc: func() input.ChangePasswordParam {
				changePasswordParam := f.NewChangePasswordParam()
				return changePasswordParam
			},
			expectStatus: http.StatusCreated,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(test.name, func(t *testing.T) {
			t.Parallel()

			inputJSON, _ := json.Marshal(test.inputFunc())
			q := make(url.Values)

			e := echo.New()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), strings.NewReader(string(inputJSON)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.SetPath(endPoint)

			operator := f.NewClaims()
			c.Set("operator", &operator)

			authUsecase := new(mocks.IAuthUsecase)
			verifyUsecase := new(mocks.IVerifyUsecase)
			authHandler := handler.NewAuthHandler(
				authUsecase,
				verifyUsecase,
			)

			authUsecase.On("ChangePassword", test.inputFunc()).Return(nil)
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
// [x] 1-3. 400: バリデーションエラー：newPasswordの桁数が7文字以下の場合
// [x] 1-4. 400: バリデーションエラー：newPasswordの桁数が21文字以上の場合
// [x] 1-5. 400: バリデーションエラー：newPasswordに英大文字を1つも含まない場合
// [x] 1-6. 400: バリデーションエラー：newPasswordに英子文字を1つも含まない場合
// [x] 1-7. 400: バリデーションエラー：newPasswordに数字を1つも含まない場合
// [x] 1-8. 400: バリデーションエラー：newPasswordに特殊文字を1つも含まない場合
// [x] 1-9. 400: バリデーションエラー：newPasswordに特殊文字を1つも含まない場合
// [x] 1-10. 500: システムエラー：変更失敗
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_Change_Abnormal(tt *testing.T) {
	var method = "POST"
	var endPoint = "/auth/change"

	tests := []struct {
		name             string
		inputFunc        func() input.ChangePasswordParam
		invalidInputFunc func() interface{}
		receive          error
		expectError      string
		expectStatus     int
	}{
		{
			name: "1-2. 400: バリデーションエラー：newPasswordの型がstring以外の場合",
			invalidInputFunc: func() interface{} {
				changePasswordInterface := f.NewChangePasswordInterface()
				changePasswordInterface.(map[string]interface{})["newPassword"] = 1
				return changePasswordInterface
			},
			receive:      nil,
			expectError:  "code=400, message={[auth] BadRequest Validation failed, newPassword: Unmarshal type error: expected=authentication.Password, got=number.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "1-3. 400: バリデーションエラー：newPasswordの桁数が7文字以下の場合",
			inputFunc: func() input.ChangePasswordParam {
				changePasswordParam := f.NewChangePasswordParam()
				changePasswordParam.NewPassword = "Abc12@"
				return changePasswordParam
			},
			receive:      nil,
			expectError:  "code=400, message={[auth] BadRequest Validation failed, newPassword: the length must be between 8 and 20.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "1-4. 400: バリデーションエラー：newPasswordの桁数が21文字以上の場合",
			inputFunc: func() input.ChangePasswordParam {
				changePasswordParam := f.NewChangePasswordParam()
				changePasswordParam.NewPassword = "AAAbbbccc111222@@@1112"
				return changePasswordParam
			},
			receive:      nil,
			expectError:  "code=400, message={[auth] BadRequest Validation failed, newPassword: the length must be between 8 and 20.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "1-6. 400: バリデーションエラー：newPasswordに英大文字を1つも含まない場合",
			inputFunc: func() input.ChangePasswordParam {
				changePasswordParam := f.NewChangePasswordParam()
				changePasswordParam.NewPassword = "aabbb11@@"
				return changePasswordParam
			},
			receive:      nil,
			expectError:  "code=400, message={[auth] BadRequest Validation failed, newPassword: must include at least one upper case letter.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "1-7. 400: バリデーションエラー：newPasswordに英小文字を1つも含まない場合",
			inputFunc: func() input.ChangePasswordParam {
				changePasswordParam := f.NewChangePasswordParam()
				changePasswordParam.NewPassword = "AABB11@@"
				return changePasswordParam
			},
			receive:      nil,
			expectError:  "code=400, message={[auth] BadRequest Validation failed, newPassword: must include at least one lower case letter.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "1-8. 400: バリデーションエラー：newPasswordに数字を1つも含まない場合",
			inputFunc: func() input.ChangePasswordParam {
				changePasswordParam := f.NewChangePasswordParam()
				changePasswordParam.NewPassword = "AAbb@@__"
				return changePasswordParam
			},
			receive:      nil,
			expectError:  "code=400, message={[auth] BadRequest Validation failed, newPassword: must include at least one digit.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "1-9. 400: バリデーションエラー：newPasswordに特殊文字を1つも含まない場合",
			inputFunc: func() input.ChangePasswordParam {
				changePasswordParam := f.NewChangePasswordParam()
				changePasswordParam.NewPassword = "xx1234Pass"
				return changePasswordParam
			},
			receive:      nil,
			expectError:  "code=400, message={[auth] BadRequest Validation failed, newPassword: must include at least one special character.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "1-10. 500: システムエラー：変更失敗",
			inputFunc: func() input.ChangePasswordParam {
				changePasswordParam := f.NewChangePasswordParam()
				changePasswordParam.NewPassword = "xx@&1234Pass"
				return changePasswordParam
			},
			receive:      fmt.Errorf("password change fail."),
			expectError:  "code=500, message={[auth] InternalServerError Unexpected error occurred",
			expectStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var inputJSON []byte
			if test.invalidInputFunc != nil {
				inputJSON, _ = json.Marshal(test.invalidInputFunc())
			} else {
				inputJSON, _ = json.Marshal(test.inputFunc())
			}

			q := make(url.Values)

			e := echo.New()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), strings.NewReader(string(inputJSON)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.SetPath(endPoint)

			operator := f.NewClaims()
			c.Set("operator", &operator)

			authUsecase := new(mocks.IAuthUsecase)
			verifyUsecase := new(mocks.IVerifyUsecase)
			authHandler := handler.NewAuthHandler(
				authUsecase,
				verifyUsecase,
			)

			if test.inputFunc != nil {
				authUsecase.On("ChangePassword", test.inputFunc()).Return(test.receive)
			}
			err := authHandler.ChangePassword(c)
			e.HTTPErrorHandler(err, c)
			if assert.Error(t, err) {
				assert.Equal(t, test.expectStatus, rec.Code)
				assert.ErrorContains(t, err, test.expectError)
			}
		})
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// POST /auth/login 正常系テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1. 201: 正常系(ログイン成功)
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_Login_Normal(tt *testing.T) {
	var method = "POST"
	var endPoint = "/auth/login"
	var dataTarget = "operator"

	tests := []struct {
		name         string
		inputFunc    func() input.LoginParam
		expectStatus int
	}{
		{
			name: "1-1. 201: 正常系(ログイン成功)",
			inputFunc: func() input.LoginParam {
				return f.NewLoginParam()
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

				inputJSON, _ := json.Marshal(test.inputFunc())
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
				authUsecase.On("Login", test.inputFunc()).Return(loginModel, nil)

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
// [x] 2-1. 400: 異常系(バリデーションエラー：operatorAccountIdの型がstring以外の場合)
// [x] 2-2. 400: 異常系(バリデーションエラー：operatorAccountIdが含まれない場合)
// [x] 2-3. 400: 異常系(バリデーションエラー：accountPasswordが含まれない場合)
// [x] 2-4. 400: 異常系(バリデーションエラー：accountPasswordが含まれない場合)
// [x] 2-5. 400: 異常系(バリデーションエラー：operatorAccountIdの型がe-mailの形式以外の場合)
// [x] 2-6. 401: 異常系(認証エラー：accountPasswordが不一致の場合)
// [x] 2-7. 500: 異常系(システムエラー：ログイン失敗)
// [x] 2-8. 503: 異常系(サービス利用不可エラー：ログイン失敗)
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_Login_Abnormal(tt *testing.T) {
	var method = "POST"
	var endPoint = "/auth/login"

	tests := []struct {
		name             string
		inputFunc        func() input.LoginParam
		invalidInputFunc func() interface{}
		receive          error
		expectError      string
		expectStatus     int
	}{
		{
			name: "2-1. 400: バリデーションエラー：operatorAccountIdの型がstring以外の場合",
			invalidInputFunc: func() interface{} {
				loginParam := f.NewLoginInterface()
				loginParam.(map[string]interface{})["operatorAccountID"] = 1
				return loginParam
			},
			receive:      nil,
			expectError:  "code=400, message={[auth] BadRequest Validation failed, operatorAccountId: Unmarshal type error: expected=string, got=number.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "2-2. 400: バリデーションエラー：operatorAccountIdが含まれない場合",
			inputFunc: func() input.LoginParam {
				loginParam := f.NewLoginParam()
				loginParam.OperatorAccountID = ""
				return loginParam
			},
			receive:      nil,
			expectError:  "code=400, message={[auth] BadRequest Validation failed, operatorAccountId: cannot be blank.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "2-3. 400: バリデーションエラー：accountPasswordの型がstring以外の場合",
			invalidInputFunc: func() interface{} {
				loginParam := f.NewLoginInterface()
				loginParam.(map[string]interface{})["accountPassword"] = 1
				return loginParam
			},
			receive:      nil,
			expectError:  "code=400, message={[auth] BadRequest Validation failed, accountPassword: Unmarshal type error: expected=string, got=number.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "2-4. 400: バリデーションエラー：accountPasswordが含まれない場合",
			inputFunc: func() input.LoginParam {
				loginParam := f.NewLoginParam()
				loginParam.AccountPassword = ""
				return loginParam
			},
			receive:      nil,
			expectError:  "code=400, message={[auth] BadRequest Validation failed, accountPassword: cannot be blank.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "2-5. 400: バリデーションエラー：operatorAccountIdの型がe-mailの形式以外の場合",
			inputFunc: func() input.LoginParam {
				loginParam := f.NewLoginParam()
				loginParam.OperatorAccountID = "xx1234Pass"
				return loginParam
			},
			receive:      nil,
			expectError:  "code=400, message={[auth] BadRequest Validation failed, operatorAccountId: must be a valid email address.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "2-6. 401: 認証エラー：accountPasswordが不一致の場合",
			inputFunc: func() input.LoginParam {
				loginParam := f.NewLoginParam()
				loginParam.AccountPassword = "xx1234Pass"
				return loginParam
			},
			receive:      common.NewCustomError(common.CustomErrorCode401, common.Err401InvalidCredentials, nil, common.HTTPErrorSourceAuth),
			expectError:  "code=401, message={[auth] Unauthorized Invalid credentials id",
			expectStatus: http.StatusUnauthorized,
		},
		{
			name: "2-7. 500: システムエラー：ログイン失敗",
			inputFunc: func() input.LoginParam {
				loginParam := f.NewLoginParam()
				return loginParam
			},
			receive:      fmt.Errorf("password change fail."),
			expectError:  "code=500, message={[auth] InternalServerError Unexpected error occurred",
			expectStatus: http.StatusInternalServerError,
		},
		{
			name: "2-8. 503: サービス利用不可エラー：ログイン失敗",
			inputFunc: func() input.LoginParam {
				loginParam := f.NewLoginParam()
				return loginParam
			},
			receive:      common.NewCustomError(common.CustomErrorCode503, common.Err503OuterService, nil, common.HTTPErrorSourceAuth),
			expectError:  "code=503, message={[auth] ServiceUnavailable Unexpected error occurred in outer service",
			expectStatus: http.StatusServiceUnavailable,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				var inputJSON []byte
				if test.invalidInputFunc != nil {
					inputJSON, _ = json.Marshal(test.invalidInputFunc())
				} else {
					inputJSON, _ = json.Marshal(test.inputFunc())
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
				if test.inputFunc != nil {
					authUsecase.On("Login", test.inputFunc()).Return(loginModel, test.receive)
				}

				err := authHandler.Login(c)
				e.HTTPErrorHandler(err, c)
				if assert.Error(t, err) {
					assert.Equal(t, test.expectStatus, rec.Code)
					assert.ErrorContains(t, err, test.expectError)
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// POST /auth/refresh 正常系テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1. 201: 正常系(リフレッシュトークン更新)
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_Refresh_Normal(tt *testing.T) {
	var method = "POST"
	var endPoint = "/auth/refresh"

	tests := []struct {
		name         string
		inputFunc    func() input.RefreshParam
		expectStatus int
	}{
		{
			name: "1-1. 201: 正常系(リフレッシュトークン更新)",
			inputFunc: func() input.RefreshParam {
				return f.NewRefreshParam()
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

				inputJSON, _ := json.Marshal(test.inputFunc())
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
				authUsecase.On("Refresh", test.inputFunc()).Return(refreshModel, nil)

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
// [x] 2-1. 400: バリデーションエラー：refreshTokenの型がstring以外の場合
// [x] 2-2. 400: バリデーションエラー：refreshTokenが含まれない場合
// [x] 2-3. 401: 認証エラー：リフレッシュ失敗)
// [x] 2-4. 500: システムエラー：リフレッシュ失敗)
// [x] 2-5. 503: サービス利用不可エラー：ログイン失敗)
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_Refresh_Abnormal(tt *testing.T) {
	var method = "POST"
	var endPoint = "/auth/refresh"

	tests := []struct {
		name             string
		inputFunc        func() input.RefreshParam
		invalidInputFunc func() interface{}
		receive          error
		expectError      string
		expectStatus     int
	}{
		{
			name: "2-1. 400: バリデーションエラー：refreshTokenの型がstring以外の場合",
			invalidInputFunc: func() interface{} {
				refreshInterface := f.NewRefreshInterface()
				refreshInterface.(map[string]interface{})["refreshToken"] = 1
				return refreshInterface
			},
			receive:      nil,
			expectError:  "code=400, message={[auth] BadRequest Validation failed, refreshToken: Unmarshal type error: expected=string, got=number.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "2-2. 400: バリデーションエラー：refreshTokenが含まれない場合",
			inputFunc: func() input.RefreshParam {
				refreshParam := f.NewRefreshParam()
				refreshParam.RefreshToken = ""
				return refreshParam
			},
			receive:      nil,
			expectError:  "code=400, message={[auth] BadRequest Validation failed, refreshToken: cannot be blank.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name:         "2-3. 401: 認証エラー：リフレッシュ失敗",
			inputFunc:    func() input.RefreshParam { return f.NewRefreshParam() },
			receive:      common.NewCustomError(common.CustomErrorCode401, common.Err401InvalidCredentials, nil, common.HTTPErrorSourceAuth),
			expectError:  "code=401, message={[auth] Unauthorized Invalid credentials id",
			expectStatus: http.StatusUnauthorized,
		},
		{
			name:         "2-4. 500: システムエラー：リフレッシュ失敗",
			inputFunc:    func() input.RefreshParam { return f.NewRefreshParam() },
			receive:      fmt.Errorf("RefreshToken get fail."),
			expectError:  "code=500, message={[auth] InternalServerError Unexpected error occurred",
			expectStatus: http.StatusInternalServerError,
		},
		{
			name:         "2-5. 503: サービス利用不可エラー：リフレッシュ失敗",
			inputFunc:    func() input.RefreshParam { return f.NewRefreshParam() },
			receive:      common.NewCustomError(common.CustomErrorCode503, common.Err503OuterService, nil, common.HTTPErrorSourceAuth),
			expectError:  "code=503, message={[auth] ServiceUnavailable Unexpected error occurred in outer service",
			expectStatus: http.StatusServiceUnavailable,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				var inputJSON []byte
				if test.invalidInputFunc != nil {
					inputJSON, _ = json.Marshal(test.invalidInputFunc())
				} else {
					inputJSON, _ = json.Marshal(test.inputFunc())
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
				if test.inputFunc != nil {
					authUsecase.On("Refresh", test.inputFunc()).Return(refreshModel, test.receive)
				}

				err := authHandler.Refresh(c)
				e.HTTPErrorHandler(err, c)
				if assert.Error(t, err) {
					assert.Equal(t, test.expectStatus, rec.Code)
					assert.ErrorContains(t, err, test.expectError)
				}
			},
		)
	}
}
