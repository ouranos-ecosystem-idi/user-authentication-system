package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"authenticator-backend/domain/common"
	f "authenticator-backend/test/fixtures"
	mocks "authenticator-backend/test/mock"
	"authenticator-backend/usecase/input"
	"authenticator-backend/usecase/output"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// /////////////////////////////////////////////////////////////////////////////////
// POST /api/v1/systemAuth/token のテストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1. 200: 正常系
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_SystemAuthToken_Normal(tt *testing.T) {
	var method = "POST"
	var endPoint = "/api/v1/systemAuth/token"

	tests := []struct {
		name         string
		input        input.VerifyTokenParam
		expectStatus int
	}{
		{
			name: "1-1. 200: 正常系",
			input: input.VerifyTokenParam{
				IDToken: f.Token,
			},
			expectStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()
				inputJSON, _ := json.Marshal(test.input)
				q := make(url.Values)

				e := echo.New()
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), strings.NewReader(string(inputJSON)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath(endPoint)
				c.Set("operatorID", f.OperatorId)

				authUsecase := new(mocks.IAuthUsecase)
				verifyUsecase := new(mocks.IVerifyUsecase)
				res := output.VerifyTokenResponse{
					OperatorID: &f.OperatorId,
				}
				authHandler := NewAuthHandler(
					authUsecase,
					verifyUsecase,
				)
				verifyUsecase.On("TokenIntrospection", test.input).Return(res, nil)

				err := authHandler.TokenIntrospection(c)
				if assert.NoError(t, err) {
					assert.Equal(t, test.expectStatus, rec.Code)
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// POST /api/v1/systemAuth/token のテストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1. 400: バリデーションエラー: idTokenが含まれていない場合
// [x] 1-2. 400: バリデーションエラー: idTokenがstring形式でない場合
// [x] 1-3. 500: システムエラー: トークン処理異常の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_SystemAuthToken(tt *testing.T) {
	var method = "POST"
	var endPoint = "/api/v1/systemAuth/token"

	tests := []struct {
		name         string
		input        input.VerifyTokenParam
		invalidInput any
		receive      error
		expectError  string
		expectStatus int
	}{
		{
			name: "1-1. 400: バリデーションエラー：idTokenが含まれていない場合",
			input: input.VerifyTokenParam{
				IDToken: "",
			},
			expectError:  "code=400, message={[auth] BadRequest Validation failed, idToken: cannot be blank.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "1-2. 400: バリデーションエラー：idTokenがstring形式でない場合",
			invalidInput: struct {
				IdToken int
			}{
				1,
			},
			expectError:  "code=400, message={[auth] BadRequest Invalid request parameters",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "1-3. 503: 外部システムエラー：接続エラー",
			input: input.VerifyTokenParam{
				IDToken: f.Token,
			},
			receive:      common.NewCustomError(common.CustomErrorCode503, common.Err503OuterService, nil, common.HTTPErrorSourceAuth),
			expectError:  "code=503, message={[auth] ServiceUnavailable",
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
				if test.invalidInput != nil {
					inputJSON, _ = json.Marshal(test.invalidInput)
				} else {
					inputJSON, _ = json.Marshal(test.input)
				}
				q := make(url.Values)

				e := echo.New()
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), strings.NewReader(string(inputJSON)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath(endPoint)
				c.Set("operatorID", f.OperatorId)

				authUsecase := new(mocks.IAuthUsecase)
				verifyUsecase := new(mocks.IVerifyUsecase)
				res := output.VerifyTokenResponse{
					OperatorID: &f.OperatorId,
				}
				authHandler := NewAuthHandler(
					authUsecase,
					verifyUsecase,
				)
				verifyUsecase.On("TokenIntrospection", test.input).Return(res, test.receive)

				err := authHandler.TokenIntrospection(c)
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
// POST /api/v1/systemAuth/apiKey のテストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1. 200: 正常系
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_SystemAuthApiKey_Normal(tt *testing.T) {
	var method = "POST"
	var endPoint = "/api/v1/systemAuth/apiKey"

	tests := []struct {
		name         string
		input        input.VerifyAPIKeyParam
		expectStatus int
	}{
		{
			name: "1-1. 200: 正常系",
			input: input.VerifyAPIKeyParam{
				IPAddress: "127.0.0.1",
				APIKey:    f.ApiKey,
			},
			expectStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()
				inputJSON, _ := json.Marshal(test.input)
				q := make(url.Values)

				e := echo.New()
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), strings.NewReader(string(inputJSON)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath(endPoint)
				c.Set("operatorID", f.OperatorId)

				authUsecase := new(mocks.IAuthUsecase)
				verifyUsecase := new(mocks.IVerifyUsecase)
				verifyUsecase.On("ApiKey", mock.Anything).Return(output.VerifyApiKeyResponse{})
				authHandler := NewAuthHandler(
					authUsecase,
					verifyUsecase,
				)

				err := authHandler.ApiKey(c)
				if assert.NoError(t, err) {
					assert.Equal(t, test.expectStatus, rec.Code)
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// POST /api/v1/systemAuth/apiKey のテストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1. 400: バリデーションエラー: IPAddressが含まれていない場合
// [x] 1-3. 400: バリデーションエラー: APIKeyが含まれていない場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_SystemAuthApiKey(tt *testing.T) {
	var method = "POST"
	var endPoint = "/api/v1/systemAuth/apiKey"

	tests := []struct {
		name         string
		input        input.VerifyAPIKeyParam
		invalidInput any
		expectError  string
		expectStatus int
	}{
		{
			name: "1-1. 400: バリデーションエラー：IPAddressが含まれていない場合",
			input: input.VerifyAPIKeyParam{
				IPAddress: "",
				APIKey:    f.ApiKey,
			},
			expectError:  "code=400, message={[auth] BadRequest Validation failed, ipAddress: cannot be blank.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "1-2. 400: バリデーションエラー：APIKeyが含まれていない場合",
			input: input.VerifyAPIKeyParam{
				IPAddress: f.IpAddress,
				APIKey:    "",
			},
			expectError:  "code=400, message={[auth] BadRequest Validation failed, apiKey: cannot be blank.",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "1-3. 400: バリデーションエラー：リクエスト不正の場合",
			invalidInput: struct {
				IpAddress int
				ApiKey    int
			}{
				1,
				1,
			},
			expectError:  "code=400, message={[auth] BadRequest Invalid request parameters",
			expectStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()
				var inputJSON []byte
				if test.invalidInput != nil {
					inputJSON, _ = json.Marshal(test.invalidInput)
				} else {
					inputJSON, _ = json.Marshal(test.input)
				}
				q := make(url.Values)

				e := echo.New()
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), strings.NewReader(string(inputJSON)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath(endPoint)
				c.Set("operatorID", f.OperatorId)

				authUsecase := new(mocks.IAuthUsecase)
				verifyUsecase := new(mocks.IVerifyUsecase)

				authHandler := NewAuthHandler(
					authUsecase,
					verifyUsecase,
				)

				err := authHandler.ApiKey(c)
				e.HTTPErrorHandler(err, c)
				if assert.Error(t, err) {
					assert.Equal(t, test.expectStatus, rec.Code)
					assert.ErrorContains(t, err, test.expectError)
				}
			},
		)
	}
}
