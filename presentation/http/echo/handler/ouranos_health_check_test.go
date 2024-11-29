package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"authenticator-backend/presentation/http/echo/handler"
	mocks "authenticator-backend/test/mock"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// /////////////////////////////////////////////////////////////////////////////////
// GET /api/v1/authInfo/health テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1. 200: 正常系
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectHandler_HealthCheck_Normal(tt *testing.T) {
	var method = "GET"
	var endPoint = "/api/v1/authInfo/health"

	tests := []struct {
		name         string
		expectStatus int
		expectBody   string
	}{
		{
			name:         "1-1. 200: 正常系",
			expectStatus: http.StatusOK,
			expectBody:   "{\"isSystemHealthy\":true}\n",
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(test.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(method, endPoint, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.SetPath(endPoint)

			operatorHandler := new(mocks.IOperatorHandler)
			operatorHandler.On("PutOperator", mock.Anything).Return(nil)
			plantHandler := new(mocks.IPlantHandler)
			plantHandler.On("PutPlant", mock.Anything).Return(nil)
			resetHandler := new(mocks.IResetHandler)
			h := handler.NewOuranosHandler(operatorHandler, plantHandler, resetHandler)

			err := h.HealthCheck(c)
			if assert.NoError(t, err) {
				assert.Equal(t, test.expectStatus, rec.Code)
				assert.Equal(t, test.expectBody, rec.Body.String())
			}
		})
	}
}
