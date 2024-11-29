package usecase_test

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"

	"authenticator-backend/domain/common"
	"authenticator-backend/domain/model/authentication"
	f "authenticator-backend/test/fixtures"
	mocks "authenticator-backend/test/mock"
	"authenticator-backend/usecase"
	"authenticator-backend/usecase/input"
	"authenticator-backend/usecase/output"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestProjectUsecase_TokenIntrospection
// Summary: This is normal test class which confirm the operation of API TokenIntrospection.
// Target: auth_verify_usecase_impl.go
// TestPattern:
// [x] 1-1. 200: 正常系
func TestProjectUsecase_TokenIntrospection(tt *testing.T) {

	var method = "GET"
	var endPoint = "/tokenIntrospection"

	res := authentication.Claims{
		OperatorID: "e03cc699-7234-31ed-86be-cc18c92208e5",
	}
	expected := output.VerifyTokenResponse{
		OperatorID: common.StringPtr("e03cc699-7234-31ed-86be-cc18c92208e5"),
	}
	tests := []struct {
		name    string
		input   input.VerifyTokenParam
		receive authentication.Claims
		expect  output.VerifyTokenResponse
	}{
		{
			name:    "1-1. 200: 正常系",
			input:   f.NewInputVerifyTokenParam(),
			receive: res,
			expect:  expected,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				q := make(url.Values)

				e := echo.New()
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath(endPoint)

				firebaseRepositoryMock := new(mocks.FirebaseRepository)
				authRepositoryMock := new(mocks.AuthRepository)
				firebaseRepositoryMock.On("VerifyIDToken", mock.Anything).Return(test.receive, nil)
				verifyUsecase := usecase.NewVerifyUsecase(firebaseRepositoryMock, authRepositoryMock)

				actual, err := verifyUsecase.TokenIntrospection(test.input)
				if assert.NoError(t, err) {
					// 実際のレスポンスと期待されるレスポンスを比較
					// 順番が実行ごとに異なるため、順不同で中身を比較
					assert.Equal(t, test.expect.OperatorID, actual.OperatorID, f.AssertMessage)
				}
			},
		)
	}
}

// TestProjectUsecase_TokenIntrospection_Abnormal
// Summary: This is abnormal test class which confirm the operation of API TokenIntrospection.
// Target: auth_verify_usecase_impl.go
// TestPattern:
// [x] 2-1. 500: 検証処理エラー
func TestProjectUsecase_TokenIntrospection_Abnormal(tt *testing.T) {

	var method = "GET"
	var endPoint = "/tokenIntrospection"

	tests := []struct {
		name         string
		input        input.VerifyTokenParam
		receive      authentication.Claims
		receiveError error
		expect       error
	}{
		{
			name:         "2-1. 500: 検証処理エラー",
			input:        f.NewInputVerifyTokenParam(),
			receive:      authentication.Claims{},
			receiveError: fmt.Errorf("検証処理エラー"),
			expect:       fmt.Errorf("検証処理エラー"),
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				q := make(url.Values)

				e := echo.New()
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath(endPoint)

				firebaseRepositoryMock := new(mocks.FirebaseRepository)
				authRepositoryMock := new(mocks.AuthRepository)
				firebaseRepositoryMock.On("VerifyIDToken", mock.Anything).Return(test.receive, test.receiveError)
				verifyUsecase := usecase.NewVerifyUsecase(firebaseRepositoryMock, authRepositoryMock)

				_, err := verifyUsecase.TokenIntrospection(test.input)
				if assert.Error(t, err) {
					// 実際のレスポンスと期待されるレスポンスを比較
					// 順番が実行ごとに異なるため、順不同で中身を比較
					assert.Equal(t, test.expect.Error(), err.Error())
				}
			},
		)
	}
}

// TestProjectUsecase_IDToken
// Summary: This is normal test class which confirm the operation of API IDToken.
// Target: auth_verify_usecase_impl.go
// TestPattern:
// [x] 1-1. 200: 正常系
func TestProjectUsecase_IDToken(tt *testing.T) {

	var method = "GET"
	var endPoint = "/tokenIntrospection"

	res := authentication.Claims{
		OperatorID: "e03cc699-7234-31ed-86be-cc18c92208e5",
	}
	tests := []struct {
		name    string
		input   input.VerifyIDTokenParam
		receive authentication.Claims
		expect  authentication.Claims
	}{
		{
			name:    "1-1. 200: 正常系",
			input:   f.NewInputVerifyIDTokenParam(),
			receive: res,
			expect:  res,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				q := make(url.Values)

				e := echo.New()
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath(endPoint)

				firebaseRepositoryMock := new(mocks.FirebaseRepository)
				authRepositoryMock := new(mocks.AuthRepository)
				firebaseRepositoryMock.On("VerifyIDToken", mock.Anything).Return(test.receive, nil)
				verifyUsecase := usecase.NewVerifyUsecase(firebaseRepositoryMock, authRepositoryMock)

				actual, err := verifyUsecase.IDToken(test.input)
				if assert.NoError(t, err) {
					// 実際のレスポンスと期待されるレスポンスを比較
					// 順番が実行ごとに異なるため、順不同で中身を比較
					assert.Equal(t, test.expect.OperatorID, actual.OperatorID, f.AssertMessage)
				}
			},
		)
	}
}

// TestProjectUsecase_IDToken_Abnormal
// Summary: This is abnormal test class which confirm the operation of API IDToken.
// Target: auth_verify_usecase_impl.go
// TestPattern:
// [x] 2-1. 500: 検証処理エラー
func TestProjectUsecase_IDToken_Abnormal(tt *testing.T) {

	var method = "GET"
	var endPoint = "/token"

	tests := []struct {
		name         string
		input        input.VerifyIDTokenParam
		receive      authentication.Claims
		receiveError error
		expect       error
	}{
		{
			name:         "2-1. 500: 検証処理エラー",
			input:        f.NewInputVerifyIDTokenParam(),
			receive:      authentication.Claims{},
			receiveError: fmt.Errorf("検証処理エラー"),
			expect:       fmt.Errorf("検証処理エラー"),
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				q := make(url.Values)

				e := echo.New()
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath(endPoint)

				firebaseRepositoryMock := new(mocks.FirebaseRepository)
				authRepositoryMock := new(mocks.AuthRepository)
				firebaseRepositoryMock.On("VerifyIDToken", mock.Anything).Return(test.receive, test.receiveError)
				verifyUsecase := usecase.NewVerifyUsecase(firebaseRepositoryMock, authRepositoryMock)

				_, err := verifyUsecase.IDToken(test.input)
				if assert.Error(t, err) {
					// 実際のレスポンスと期待されるレスポンスを比較
					// 順番が実行ごとに異なるため、順不同で中身を比較
					assert.Equal(t, test.expect.Error(), err.Error())
				}
			},
		)
	}
}

// TestProjectUsecase_ApiKey
// Summary: This is normal test class which confirm the operation of API KEY.
// Target: auth_verify_usecase_impl.go
// TestPattern:
// [x] 1-1. 200: 両方OK
// [x] 1-2. 200: APIKEYのみNG
// [x] 1-3. 200: IPアドレスのみNG
// [x] 1-4. 200: 両方NG
func TestProjectUsecase_ApiKey(tt *testing.T) {

	var method = "GET"
	var endPoint = "/apikey"

	resKeys := authentication.APIKeys{
		authentication.APIKey{
			APIKey: "36cfd2a8-9f45-0766-77d3-7098c1336a32",
		},
	}
	resCidrs := authentication.Cidrs{
		&authentication.Cidr{
			Cidr: "127.0.0.1/32",
		},
	}
	tests := []struct {
		name         string
		inputFunc    func() input.VerifyAPIKeyParam
		receiveKeys  authentication.APIKeys
		receiveCidrs authentication.Cidrs
		expect       output.VerifyApiKeyResponse
	}{
		{
			name: "1-1. 200: 両方OK",
			inputFunc: func() input.VerifyAPIKeyParam {
				return f.NewInputVerifyAPIKeyParam()
			},
			receiveKeys:  resKeys,
			receiveCidrs: resCidrs,
			expect: output.VerifyApiKeyResponse{
				IsAPIKeyValid:    true,
				IsIPAddressValid: true,
			},
		},
		{
			name: "1-2. 200: APIKEYのみNG",
			inputFunc: func() input.VerifyAPIKeyParam {
				InputVerifyAPIKeyParam := f.NewInputVerifyAPIKeyParam()
				InputVerifyAPIKeyParam.APIKey = "APIKEY2"
				return InputVerifyAPIKeyParam
			},
			receiveKeys:  resKeys,
			receiveCidrs: resCidrs,
			expect: output.VerifyApiKeyResponse{
				IsAPIKeyValid:    false,
				IsIPAddressValid: true,
			},
		},
		{
			name: "1-3. 200: IPアドレスのみNG",
			inputFunc: func() input.VerifyAPIKeyParam {
				InputVerifyAPIKeyParam := f.NewInputVerifyAPIKeyParam()
				InputVerifyAPIKeyParam.IPAddress = "127.0.0.2"
				return InputVerifyAPIKeyParam
			},
			receiveKeys:  resKeys,
			receiveCidrs: resCidrs,
			expect: output.VerifyApiKeyResponse{
				IsAPIKeyValid:    true,
				IsIPAddressValid: false,
			},
		},
		{
			name: "1-4. 200: 両方NG",
			inputFunc: func() input.VerifyAPIKeyParam {
				InputVerifyAPIKeyParam := f.NewInputVerifyAPIKeyParam()
				InputVerifyAPIKeyParam.APIKey = "APIKEY2"
				InputVerifyAPIKeyParam.IPAddress = "127.0.0.2"
				return InputVerifyAPIKeyParam
			},
			receiveKeys:  resKeys,
			receiveCidrs: resCidrs,
			expect: output.VerifyApiKeyResponse{
				IsAPIKeyValid:    false,
				IsIPAddressValid: false,
			},
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				q := make(url.Values)

				e := echo.New()
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath(endPoint)

				firebaseRepositoryMock := new(mocks.FirebaseRepository)
				authRepositoryMock := new(mocks.AuthRepository)
				authRepositoryMock.On("ListAPIKeys", mock.Anything).Return(test.receiveKeys, nil)
				authRepositoryMock.On("ListCidrs", mock.Anything).Return(test.receiveCidrs, nil)
				verifyUsecase := usecase.NewVerifyUsecase(firebaseRepositoryMock, authRepositoryMock)

				actual := verifyUsecase.ApiKey(test.inputFunc())
				assert.Equal(t, test.expect.IsAPIKeyValid, actual.IsAPIKeyValid, f.AssertMessage)
				assert.Equal(t, test.expect.IsIPAddressValid, actual.IsIPAddressValid, f.AssertMessage)
			},
		)
	}
}

// TestProjectUsecase_ApiKey_Abnormal
// Summary: This is abnormal test class which confirm the operation of API KEY.
// Target: auth_verify_usecase_impl.go
// TestPattern:
// [x] 2-1. 500: 検証処理エラー(Key)
// [x] 2-1. 500: 検証処理エラー(Cidrs)
func TestProjectUsecase_ApiKey_Abnormal(tt *testing.T) {

	var method = "GET"
	var endPoint = "/token"

	resKeys := authentication.APIKeys{
		authentication.APIKey{
			APIKey: "36cfd2a8-9f45-0766-77d3-7098c1336a32",
		},
	}
	tests := []struct {
		name              string
		inputFunc         func() input.VerifyAPIKeyParam
		receiveKeys       authentication.APIKeys
		receiveKeysError  error
		receiveCidrs      authentication.Cidrs
		receiveCidrsError error
		expect            output.VerifyApiKeyResponse
	}{
		{
			name: "2-1. 500: 検証処理エラー(Key)",
			inputFunc: func() input.VerifyAPIKeyParam {
				return f.NewInputVerifyAPIKeyParam()
			},
			receiveKeys:      authentication.APIKeys{},
			receiveKeysError: fmt.Errorf("検証処理エラー"),
			receiveCidrs:     authentication.Cidrs{},
			expect: output.VerifyApiKeyResponse{
				IsAPIKeyValid:    false,
				IsIPAddressValid: false,
			},
		},
		{
			name: "2-2. 500: 検証処理エラー(Cidrs)",
			inputFunc: func() input.VerifyAPIKeyParam {
				return f.NewInputVerifyAPIKeyParam()
			},
			receiveKeys:       resKeys,
			receiveCidrs:      authentication.Cidrs{},
			receiveCidrsError: fmt.Errorf("検証処理エラー"),
			expect: output.VerifyApiKeyResponse{
				IsAPIKeyValid:    true,
				IsIPAddressValid: false,
			},
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				q := make(url.Values)

				e := echo.New()
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath(endPoint)

				firebaseRepositoryMock := new(mocks.FirebaseRepository)
				authRepositoryMock := new(mocks.AuthRepository)
				authRepositoryMock.On("ListAPIKeys", mock.Anything).Return(test.receiveKeys, test.receiveKeysError)
				authRepositoryMock.On("ListCidrs", mock.Anything).Return(test.receiveCidrs, test.receiveCidrsError)
				verifyUsecase := usecase.NewVerifyUsecase(firebaseRepositoryMock, authRepositoryMock)

				actual := verifyUsecase.ApiKey(test.inputFunc())
				assert.Equal(t, test.expect.IsAPIKeyValid, actual.IsAPIKeyValid, f.AssertMessage)
				assert.Equal(t, test.expect.IsIPAddressValid, actual.IsIPAddressValid, f.AssertMessage)
			},
		)
	}
}
