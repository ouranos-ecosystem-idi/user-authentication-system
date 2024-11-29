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

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// /////////////////////////////////////////////////////////////////////////////////
// POST /dataReset テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1. 200: 正常系
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_Reset_Normal2(tt *testing.T) {
	var method = "POST"
	var endPoint = "/dataReset"

	tests := []struct {
		name         string
		receive      error
		expectStatus int
	}{
		{
			name:         "2-1. 200: 正常系",
			expectStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(test.name, func(t *testing.T) {
			t.Parallel()

			q := make(url.Values)

			e := echo.New()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.SetPath(endPoint)

			c.Set("operatorID", f.OperatorId)

			operatorUsecase := new(mocks.IOperatorUsecase)
			operatorHandler := handler.NewOperatorHandler(operatorUsecase)
			plantUsecase := new(mocks.IPlantUsecase)
			plantHandler := handler.NewPlantHandler(plantUsecase)
			resetUsecase := new(mocks.IResetUsecase)
			resetHandler := handler.NewResetHandler(resetUsecase)
			resetUsecase.On("Reset", mock.Anything).Return(test.receive)
			h := handler.NewOuranosHandler(operatorHandler, plantHandler, resetHandler)
			err := h.Reset(c)
			if assert.NoError(t, err) {
				assert.Equal(t, test.expectStatus, rec.Code)
				resetUsecase.AssertExpectations(t)
			}
		})
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// POST /dataReset テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1. 200: 正常系
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_Reset_Normal(tt *testing.T) {
	var method = "POST"
	var endPoint = "/dataReset"

	tests := []struct {
		name         string
		receive      error
		expectStatus int
	}{
		{
			name:         "2-1. 200: 正常系",
			expectStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(test.name, func(t *testing.T) {
			t.Parallel()

			q := make(url.Values)

			e := echo.New()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.SetPath(endPoint)

			c.Set("operatorID", f.OperatorId)

			resetUsecase := new(mocks.IResetUsecase)
			resetHandler := handler.NewResetHandler(
				resetUsecase,
			)

			resetUsecase.On("Reset", mock.Anything).Return(test.receive)
			err := resetHandler.Reset(c)
			if assert.NoError(t, err) {
				assert.Equal(t, test.expectStatus, rec.Code)
				resetUsecase.AssertExpectations(t)
			}
		})
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// POST /dataReset テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1. 500: 異常系
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_Reset_Abnormal(tt *testing.T) {
	var method = "POST"
	var endPoint = "/dataReset"

	tests := []struct {
		name         string
		receive      error
		expectError  string
		expectStatus int
	}{
		{
			name:         "2-1. 500: 異常系",
			receive:      fmt.Errorf("Access Error"),
			expectError:  "code=500, message={[auth] InternalServerError Unexpected error occurred",
			expectStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(test.name, func(t *testing.T) {
			t.Parallel()

			q := make(url.Values)

			e := echo.New()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.SetPath(endPoint)

			c.Set("operatorID", f.OperatorId)

			resetUsecase := new(mocks.IResetUsecase)
			resetHandler := handler.NewResetHandler(
				resetUsecase,
			)

			resetUsecase.On("Reset", mock.Anything).Return(test.receive)
			err := resetHandler.Reset(c)
			e.HTTPErrorHandler(err, c)
			if assert.Error(t, err) {
				assert.Equal(t, test.expectStatus, rec.Code)
				assert.ErrorContains(t, err, test.expectError)
			}
		})
	}
}
