package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"authenticator-backend/presentation/http/echo/handler"
	f "authenticator-backend/test/fixtures"
	mocks "authenticator-backend/test/mock"
	testhelper "authenticator-backend/test/test_helper"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// /////////////////////////////////////////////////////////////////////////////////
// Get /api/v1/ テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1. 200: 正常系：operatorの場合
// [x] 1-2. 200: 正常系：plantの場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_Get_Normal(tt *testing.T) {
	var method = "GET"
	var endPoint = "/api/v1"

	tests := []struct {
		name              string
		modifyQueryParams func(q url.Values)
	}{
		{
			name: "1-1. 200: 正常系：operatorの場合",
			modifyQueryParams: func(q url.Values) {
				q.Set("dataTarget", "operator")
			},
		},
		{
			name: "1-2. 200: 正常系：plantの場合",
			modifyQueryParams: func(q url.Values) {
				q.Set("dataTarget", "plant")
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
				test.modifyQueryParams(q)

				e := echo.New()
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath(endPoint)
				c.Set("operatorID", f.OperatorId)

				operatorHandler := new(mocks.IOperatorHandler)
				operatorHandler.On("GetOperator", mock.Anything).Return(nil)
				plantHandler := new(mocks.IPlantHandler)
				plantHandler.On("GetPlant", mock.Anything).Return(nil)
				resetHandler := new(mocks.IResetHandler)
				h := handler.NewOuranosHandler(operatorHandler, plantHandler, resetHandler)

				err := h.GetAuthInfo(c)
				assert.NoError(t, err)
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Get /api/v1/ テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1. 400: バリデーションエラー：dataTargetの値が未指定の場合
// [x] 1-2. 400: バリデーションエラー：dataTargetの値が不正の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_Get(tt *testing.T) {
	var method = "GET"
	var endPoint = "/api/v1"

	tests := []struct {
		name              string
		modifyQueryParams func(q url.Values)
		expectError       string
		expectStatus      int
	}{
		{
			name: "1-1. 400: バリデーションエラー：dataTargetが含まれない場合",
			modifyQueryParams: func(q url.Values) {
			},
			expectError:  "code=400, message={[auth] BadRequest Invalid request parameters, dataTarget: Unexpected query parameter",
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "1-2. 400: バリデーションエラー：dataTargetがoperator以外の場合",
			modifyQueryParams: func(q url.Values) {
				q.Set("dataTarget", "hoge")
			},
			expectError:  "code=400, message={[auth] BadRequest Invalid request parameters, dataTarget: Unexpected query parameter",
			expectStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				q := make(url.Values)
				test.modifyQueryParams(q)

				e := echo.New()
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath(endPoint)
				c.Set("operatorID", f.OperatorId)
				host := ""

				h := testhelper.NewMockHandler(host)

				err := h.GetAuthInfo(c)
				e.HTTPErrorHandler(err, c)
				if assert.Error(t, err) {
					fmt.Println(err)
					assert.Equal(t, test.expectStatus, rec.Code)
					assert.ErrorContains(t, err, test.expectError)
				}
			},
		)
	}
}
