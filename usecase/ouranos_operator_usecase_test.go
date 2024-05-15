package usecase_test

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"authenticator-backend/domain/common"
	"authenticator-backend/domain/model/traceability"
	f "authenticator-backend/test/fixtures"
	mocks "authenticator-backend/test/mock"
	"authenticator-backend/usecase"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// TestProjectUsecase_GetOperator
// Summary: This is normal test class which confirm the operation of API #2 GetOperator.
// Target: ouranos_operator_usecase.go
// TestPattern:
// [x] 1-1. 200: 全項目応答(公開事業者識別子)
// [x] 1-2. 200: 全項目応答(内部事業者識別子)
func TestProjectUsecase_GetOperator(tt *testing.T) {

	var method = "GET"
	var endPoint = "/auth/v1/authInfo"
	var dataTarget = "operator"

	operator := traceability.OperatorEntityModel{
		OperatorID:       uuid.MustParse(f.OperatorID),
		OperatorName:     f.OperatorName,
		OperatorAddress:  f.OperatorAddress,
		OpenOperatorID:   f.OpenOperatorID,
		GlobalOperatorID: &f.GlobalOperatorID,
		DeletedAt: gorm.DeletedAt{
			Time:  time.Now(),
			Valid: true,
		},
		CreatedAt:         time.Now(),
		CreatedOperatorID: f.OperatorID,
		UpdatedAt:         time.Now(),
		UpdatedOperatorID: f.OperatorID,
	}

	expected := traceability.OperatorModel{
		OperatorID:      uuid.MustParse(f.OperatorID),
		OperatorName:    f.OperatorName,
		OperatorAddress: f.OperatorAddress,
		OpenOperatorID:  f.OpenOperatorID,
		OperatorAttribute: traceability.OperatorAttribute{
			GlobalOperatorID: &f.GlobalOperatorID,
		},
	}

	tests := []struct {
		name               string
		openOperatorSearch bool
		input              traceability.GetOperatorInput
		receive            traceability.OperatorEntityModel
		expect             traceability.OperatorModel
	}{
		{
			name:               "1-1. 200: 全項目応答(公開事業者識別子)",
			openOperatorSearch: true,
			input:              f.NewGetOperatorInput(true),
			receive:            operator,
			expect:             expected,
		},
		{
			name:               "1-2. 200: 全項目応答(内部事業者識別子)",
			openOperatorSearch: false,
			input:              f.NewGetOperatorInput(false),
			receive:            operator,
			expect:             expected,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				q := make(url.Values)
				q.Set("dataTarget", dataTarget)

				e := echo.New()
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath(endPoint)
				c.Set("operatorID", f.OperatorID)

				ouranosRepositoryMock := new(mocks.OuranosRepository)
				if test.openOperatorSearch {
					ouranosRepositoryMock.On("GetOperatorByOpenOperatorID", mock.Anything).Return(test.receive, nil)
				} else {
					ouranosRepositoryMock.On("GetOperator", mock.Anything).Return(test.receive, nil)
				}
				operatorUsecase := usecase.NewOperatorUsecase(ouranosRepositoryMock)

				actual, err := operatorUsecase.GetOperator(test.input)
				if assert.NoError(t, err) {
					// 実際のレスポンスと期待されるレスポンスを比較
					// 順番が実行ごとに異なるため、順不同で中身を比較
					assert.Equal(t, test.expect.OperatorID, actual.OperatorID, f.AssertMessage)
					assert.Equal(t, test.expect.OperatorName, actual.OperatorName, f.AssertMessage)
					assert.Equal(t, test.expect.OperatorAddress, actual.OperatorAddress, f.AssertMessage)
					assert.Equal(t, test.expect.OpenOperatorID, actual.OpenOperatorID, f.AssertMessage)
					assert.Equal(t, test.expect.OperatorAttribute.GlobalOperatorID, actual.OperatorAttribute.GlobalOperatorID, f.AssertMessage)
				}
			},
		)
	}
}

// TestProjectUsecase_GetOperator_Abnormal
// Summary: This is abnormal test class which confirm the operation of API #2 GetOperator.
// Target: ouranos_operator_usecase.go
// TestPattern:
// [x] 2-1. 400: データ取得エラー(公開事業者識別子)
// [x] 2-2. 400: データ取得エラー(内部事業者識別子)
func TestProjectUsecase_GetOperator_Abnormal(tt *testing.T) {

	var method = "GET"
	var endPoint = "/auth/v1/authInfo"
	var dataTarget = "operator"

	dsResGetError := fmt.Errorf("DB AccessError")

	tests := []struct {
		name               string
		openOperatorSearch bool
		input              traceability.GetOperatorInput
		receive            error
		expect             error
	}{
		{
			name:               "2-1. 400: データ取得エラー(公開事業者識別子)",
			openOperatorSearch: true,
			input:              f.NewGetOperatorInput(true),
			receive:            dsResGetError,
			expect:             dsResGetError,
		},
		{
			name:               "2-2. 400: データ取得エラー(内部事業者識別子)",
			openOperatorSearch: false,
			input:              f.NewGetOperatorInput(false),
			receive:            dsResGetError,
			expect:             dsResGetError,
		},
		{
			name:               "2-3. 400: データ取得エラー(公開事業者識別子)",
			openOperatorSearch: true,
			input:              f.NewGetOperatorInput(true),
			receive:            gorm.ErrRecordNotFound,
			expect:             gorm.ErrRecordNotFound,
		},
		{
			name:               "2-4. 400: データ取得エラー(内部事業者識別子)",
			openOperatorSearch: false,
			input:              f.NewGetOperatorInput(false),
			receive:            gorm.ErrRecordNotFound,
			expect:             gorm.ErrRecordNotFound,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				q := make(url.Values)
				q.Set("dataTarget", dataTarget)

				e := echo.New()
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath(endPoint)
				c.Set("operatorID", f.OperatorID)

				ouranosRepositoryMock := new(mocks.OuranosRepository)
				if test.openOperatorSearch {
					ouranosRepositoryMock.On("GetOperatorByOpenOperatorID", mock.Anything).Return(traceability.OperatorEntityModel{}, test.receive)
				} else {
					ouranosRepositoryMock.On("GetOperator", mock.Anything).Return(traceability.OperatorEntityModel{}, test.receive)
				}

				operatorUsecase := usecase.NewOperatorUsecase(ouranosRepositoryMock)

				_, err := operatorUsecase.GetOperator(test.input)
				if assert.Error(t, err) {
					// 実際のレスポンスと期待されるレスポンスを比較
					// 順番が実行ごとに異なるため、順不同で中身を比較
					assert.Equal(t, test.expect.Error(), err.Error())
				}
			},
		)
	}
}

// TestProjectUsecase_GetOperators
// Summary: This is normal test class which confirm the operation of API #17 GetOperators.
// Target: ouranos_operator_usecase.go
// TestPattern:
// [x] 1-1. 200: 全項目応答
func TestProjectUsecase_GetOperators(tt *testing.T) {

	var method = "GET"
	var endPoint = "/auth/v1/authInfo"
	var dataTarget = "operator"

	operator := traceability.OperatorEntityModels{
		{
			OperatorID:       uuid.MustParse(f.OperatorID),
			OperatorName:     f.OperatorName,
			OperatorAddress:  f.OperatorAddress,
			OpenOperatorID:   f.OpenOperatorID,
			GlobalOperatorID: &f.GlobalOperatorID,
			DeletedAt: gorm.DeletedAt{
				Time:  time.Now(),
				Valid: true,
			},
			CreatedAt:         time.Now(),
			CreatedOperatorID: f.OperatorID,
			UpdatedAt:         time.Now(),
			UpdatedOperatorID: f.OperatorID,
		},
	}

	expected := traceability.OperatorModels{
		{
			OperatorID:      uuid.MustParse(f.OperatorID),
			OperatorName:    f.OperatorName,
			OperatorAddress: f.OperatorAddress,
			OpenOperatorID:  f.OpenOperatorID,
			OperatorAttribute: traceability.OperatorAttribute{
				GlobalOperatorID: &f.GlobalOperatorID,
			},
		},
	}

	tests := []struct {
		name    string
		input   traceability.GetOperatorsInput
		receive traceability.OperatorEntityModels
		expect  traceability.OperatorModels
	}{
		{
			name:    "1-1. 200: 全項目応答",
			input:   f.NewGetOperatorsInput(),
			receive: operator,
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
				q.Set("dataTarget", dataTarget)

				e := echo.New()
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath(endPoint)
				c.Set("operatorID", f.OperatorID)

				ouranosRepositoryMock := new(mocks.OuranosRepository)
				ouranosRepositoryMock.On("GetOperators", mock.Anything).Return(test.receive, nil)
				operatorUsecase := usecase.NewOperatorUsecase(ouranosRepositoryMock)

				actual, err := operatorUsecase.GetOperators(test.input)
				if assert.NoError(t, err) {
					// 実際のレスポンスと期待されるレスポンスを比較
					// 順番が実行ごとに異なるため、順不同で中身を比較
					assert.ElementsMatch(t, test.expect, actual, f.AssertMessage)
				}
			},
		)
	}
}

// TestProjectUsecase_GetOperators_Abnormal
// Summary: This is abnormal test class which confirm the operation of API #17 GetOperators.
// Target: ouranos_operator_usecase.go
// TestPattern:
// [x] 2-1. 400: データ取得エラー
func TestProjectUsecase_GetOperators_Abnormal(tt *testing.T) {

	var method = "GET"
	var endPoint = "/auth/v1/authInfo"
	var dataTarget = "operator"

	dsResGetError := fmt.Errorf("DB AccessError")

	tests := []struct {
		name    string
		input   traceability.GetOperatorsInput
		receive error
		expect  error
	}{
		{
			name:    "2-1. 400: データ取得エラー",
			input:   f.NewGetOperatorsInput(),
			receive: dsResGetError,
			expect:  dsResGetError,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				q := make(url.Values)
				q.Set("dataTarget", dataTarget)

				e := echo.New()
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath(endPoint)
				c.Set("operatorID", f.OperatorID)

				ouranosRepositoryMock := new(mocks.OuranosRepository)
				ouranosRepositoryMock.On("GetOperators", mock.Anything).Return(traceability.OperatorEntityModels{}, test.receive)

				operatorUsecase := usecase.NewOperatorUsecase(ouranosRepositoryMock)

				_, err := operatorUsecase.GetOperators(test.input)
				if assert.Error(t, err) {
					// 実際のレスポンスと期待されるレスポンスを比較
					// 順番が実行ごとに異なるため、順不同で中身を比較
					assert.Equal(t, test.expect.Error(), err.Error())
				}
			},
		)
	}
}

// TestProjectUsecase_PutOperator
// Summary: This is normal test class which confirm the operation of API #1 PutOperator.
// Target: ouranos_operator_usecase.go
// TestPattern:
// [x] 1-1. 200: 全項目応答
func TestProjectUsecase_PutOperator(tt *testing.T) {

	var method = "PUT"
	var endPoint = "/auth/v1/authInfo"
	var dataTarget = "operator"

	operator := traceability.OperatorEntityModel{
		OperatorID:       uuid.MustParse(f.OperatorID),
		OperatorName:     f.OperatorName,
		OperatorAddress:  f.OperatorAddress,
		OpenOperatorID:   f.OpenOperatorID,
		GlobalOperatorID: &f.GlobalOperatorID,
		DeletedAt: gorm.DeletedAt{
			Time:  time.Now(),
			Valid: true,
		},
		CreatedAt:         time.Now(),
		CreatedOperatorID: f.OperatorID,
		UpdatedAt:         time.Now(),
		UpdatedOperatorID: f.OperatorID,
	}

	tests := []struct {
		name    string
		input   traceability.OperatorModel
		receive traceability.OperatorEntityModel
		expect  traceability.OperatorModel
	}{
		{
			name:    "1-1. 200: 全項目応答",
			input:   f.NewOperatorModel(f.OperatorID, f.OpenOperatorID),
			receive: operator,
			expect:  f.NewOperatorModel(f.OperatorID, f.OpenOperatorID),
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				t.Parallel()

				q := make(url.Values)
				q.Set("dataTarget", dataTarget)

				e := echo.New()
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath(endPoint)
				c.Set("operatorID", f.OperatorID)

				ouranosRepositoryMock := new(mocks.OuranosRepository)
				ouranosRepositoryMock.On("GetOperator", mock.Anything).Return(test.receive, nil)
				ouranosRepositoryMock.On("PutOperator", mock.Anything).Return(test.receive, nil)

				operatorUsecase := usecase.NewOperatorUsecase(ouranosRepositoryMock)

				actual, err := operatorUsecase.PutOperator(test.input)
				if assert.NoError(t, err) {
					// 実際のレスポンスと期待されるレスポンスを比較
					// 順番が実行ごとに異なるため、順不同で中身を比較
					assert.Equal(t, test.expect.OperatorID, actual.OperatorID, f.AssertMessage)
					assert.Equal(t, test.expect.OperatorName, actual.OperatorName, f.AssertMessage)
					assert.Equal(t, test.expect.OperatorAddress, actual.OperatorAddress, f.AssertMessage)
					assert.Equal(t, test.expect.OpenOperatorID, actual.OpenOperatorID, f.AssertMessage)
					assert.Equal(t, test.expect.OperatorAttribute.GlobalOperatorID, actual.OperatorAttribute.GlobalOperatorID, f.AssertMessage)
				}
			},
		)
	}
}

// TestProjectUsecase_PutOperator_Abnormal
// Summary: This is abnormal test class which confirm the operation of API #1 PutOperator.
// Target: ouranos_operator_usecase.go
// TestPattern:
// [x] 2-1. 400: データ取得エラー
// [x] 2-2. 400: データ更新エラー
func TestProjectUsecase_PutOperator_Abnormal(tt *testing.T) {

	var method = "PUT"
	var endPoint = "/auth/v1/authInfo"
	var dataTarget = "operator"

	operator := traceability.OperatorEntityModel{
		OperatorID:       uuid.MustParse(f.OperatorID),
		OperatorName:     f.OperatorName,
		OperatorAddress:  f.OperatorAddress,
		OpenOperatorID:   f.OpenOperatorID,
		GlobalOperatorID: &f.GlobalOperatorID,
		DeletedAt: gorm.DeletedAt{
			Time:  time.Now(),
			Valid: true,
		},
		CreatedAt:         time.Now(),
		CreatedOperatorID: f.OperatorID,
		UpdatedAt:         time.Now(),
		UpdatedOperatorID: f.OperatorID,
	}

	dsResGetError := fmt.Errorf("DB AccessError")

	tests := []struct {
		name         string
		input        traceability.OperatorModel
		receive      *traceability.OperatorEntityModel
		receiveError error
		expect       error
	}{
		{
			name:         "2-1. 400: データ取得エラー",
			input:        f.NewOperatorModel(f.OperatorID, f.OpenOperatorID),
			receive:      nil,
			receiveError: dsResGetError,
			expect:       dsResGetError,
		},
		{
			name:         "2-2. 400: データ更新エラー",
			input:        f.NewOperatorModel(f.OperatorID, f.OpenOperatorID),
			receive:      &operator,
			receiveError: dsResGetError,
			expect:       dsResGetError,
		},
		{
			name:         "2-3. 400: データ更新エラー",
			input:        f.NewOperatorModel(f.OperatorID, "AAAA-12345"),
			receive:      &operator,
			receiveError: common.NewCustomError(common.CustomErrorCode400, fmt.Errorf(common.FieldIsImutable("openOperatorID")).Error(), nil, common.HTTPErrorSourceAuth),
			expect:       common.NewCustomError(common.CustomErrorCode400, fmt.Errorf(common.FieldIsImutable("openOperatorID")).Error(), nil, common.HTTPErrorSourceAuth),
		},
		{
			name:    "2-4. 500: データ更新エラー",
			input:   f.NewOperatorModel(f.OperatorID, f.OpenOperatorID),
			receive: &operator,
			receiveError: &pgconn.PgError{
				Severity:         "ERROR",
				Code:             "23505",
				Message:          "unique_global_operator_id",
				Detail:           "",
				Hint:             "",
				Position:         1,
				InternalPosition: 1,
				InternalQuery:    "",
				Where:            "",
				SchemaName:       "",
				TableName:        "",
				ColumnName:       "",
				DataTypeName:     "",
				ConstraintName:   "",
				File:             "",
				Line:             1,
				Routine:          "",
			},
			expect: common.NewCustomError(common.CustomErrorCode400, "globalOperatorId: GlobalOperatorId is already exists.", nil, common.HTTPErrorSourceAuth),
		},
		{
			name:    "2-4. 500: データ更新エラー",
			input:   f.NewOperatorModel(f.OperatorID, f.OpenOperatorID),
			receive: &operator,
			receiveError: &pgconn.PgError{
				Severity:         "ERROR",
				Code:             "23504",
				Message:          "unique_global_operator_id",
				Detail:           "",
				Hint:             "",
				Position:         1,
				InternalPosition: 1,
				InternalQuery:    "",
				Where:            "",
				SchemaName:       "",
				TableName:        "",
				ColumnName:       "",
				DataTypeName:     "",
				ConstraintName:   "",
				File:             "",
				Line:             1,
				Routine:          "",
			},
			expect: &pgconn.PgError{
				Severity:         "ERROR",
				Code:             "23504",
				Message:          "unique_global_operator_id",
				Detail:           "",
				Hint:             "",
				Position:         1,
				InternalPosition: 1,
				InternalQuery:    "",
				Where:            "",
				SchemaName:       "",
				TableName:        "",
				ColumnName:       "",
				DataTypeName:     "",
				ConstraintName:   "",
				File:             "",
				Line:             1,
				Routine:          "",
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
				q.Set("dataTarget", dataTarget)

				e := echo.New()
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(method, endPoint+"?"+q.Encode(), nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath(endPoint)
				c.Set("operatorID", f.OperatorID)

				ouranosRepositoryMock := new(mocks.OuranosRepository)
				if test.receive != nil {
					ouranosRepositoryMock.On("GetOperator", mock.Anything).Return(*test.receive, nil)
					ouranosRepositoryMock.On("PutOperator", mock.Anything).Return(traceability.OperatorEntityModel{}, test.receiveError)
				} else {
					ouranosRepositoryMock.On("GetOperator", mock.Anything).Return(traceability.OperatorEntityModel{}, test.receiveError)
				}

				operatorUsecase := usecase.NewOperatorUsecase(ouranosRepositoryMock)

				_, err := operatorUsecase.PutOperator(test.input)
				if assert.Error(t, err) {
					// 実際のレスポンスと期待されるレスポンスを比較
					// 順番が実行ごとに異なるため、順不同で中身を比較
					assert.Equal(t, test.expect.Error(), err.Error())
				}
			},
		)
	}
}
