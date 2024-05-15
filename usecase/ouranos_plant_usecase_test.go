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

// TestProjectUsecase_ListPlants
// Summary: This is abnormal test class which confirm the operation of API #4 GetPlant.
// Target: ouranos_plant_usecase.go
// TestPattern:
// [x] 1-1. 200: 全項目応答
func TestProjectUsecase_ListPlants(tt *testing.T) {

	var method = "GET"
	var endPoint = "/auth/v1/authInfo"
	var dataTarget = "plant"

	plant := traceability.PlantEntityModels{
		{
			PlantID:       uuid.MustParse(f.PlantID),
			OperatorID:    uuid.MustParse(f.OperatorID),
			PlantName:     f.OperatorName,
			PlantAddress:  f.OperatorAddress,
			OpenPlantID:   f.OpenOperatorID,
			GlobalPlantID: &f.GlobalOperatorID,
			DeletedAt: gorm.DeletedAt{
				Time:  time.Now(),
				Valid: true,
			},
			CreatedAt:     time.Now(),
			CreatedUserID: f.OperatorID,
			UpdatedAt:     time.Now(),
			UpdatedUserID: f.OperatorID,
		},
	}

	expected := traceability.PlantModels{
		{
			PlantID:      uuid.MustParse(f.PlantID),
			OperatorID:   uuid.MustParse(f.OperatorID),
			PlantName:    f.OperatorName,
			PlantAddress: f.OperatorAddress,
			OpenPlantID:  f.OpenOperatorID,
			PlantAttribute: traceability.PlantAttribute{
				GlobalPlantID: &f.GlobalOperatorID,
			},
		},
	}

	tests := []struct {
		name    string
		input   traceability.GetPlantModel
		receive traceability.PlantEntityModels
		expect  traceability.PlantModels
	}{
		{
			name:    "1-1. 200: 全項目応答",
			input:   f.NewGetPlantModel(),
			receive: plant,
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
				ouranosRepositoryMock.On("ListPlantsByOperatorID", mock.Anything).Return(test.receive, nil)
				plantUsecase := usecase.NewPlantUsecase(ouranosRepositoryMock)

				actual, err := plantUsecase.ListPlants(test.input)
				if assert.NoError(t, err) {
					// 実際のレスポンスと期待されるレスポンスを比較
					// 順番が実行ごとに異なるため、順不同で中身を比較
					assert.ElementsMatch(t, test.expect, actual, f.AssertMessage)
				}
			},
		)
	}
}

// TestProjectUsecase_ListPlants_Abnormal
// Summary: This is abnormal test class which confirm the operation of API #4 GetPlant.
// Target: ouranos_plant_usecase.go
// TestPattern:
// [x] 2-1. 400: データ取得エラー
func TestProjectUsecase_ListPlants_Abnormal(tt *testing.T) {

	var method = "GET"
	var endPoint = "/auth/v1/authInfo"
	var dataTarget = "operator"

	dsResGetError := fmt.Errorf("DB AccessError")

	tests := []struct {
		name    string
		input   traceability.GetPlantModel
		receive error
		expect  error
	}{
		{
			name:    "2-1. 400: データ取得エラー",
			input:   f.NewGetPlantModel(),
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
				ouranosRepositoryMock.On("ListPlantsByOperatorID", mock.Anything).Return(traceability.PlantEntityModels{}, test.receive)

				plantUsecase := usecase.NewPlantUsecase(ouranosRepositoryMock)

				_, err := plantUsecase.ListPlants(test.input)
				if assert.Error(t, err) {
					// 実際のレスポンスと期待されるレスポンスを比較
					// 順番が実行ごとに異なるため、順不同で中身を比較
					assert.Equal(t, test.expect.Error(), err.Error())
				}
			},
		)
	}
}

// TestProjectUsecase_PutPlant
// Summary: This is normal test class which confirm the operation of API #3 PutPlant.
// Target: ouranos_plant_usecase.go
// TestPattern:
// [x] 1-1. 200: 全項目応答(新規)
// [x] 1-2. 200: 全項目応答(更新)
func TestProjectUsecase_PutPlant(tt *testing.T) {

	var method = "PUT"
	var endPoint = "/auth/v1/authInfo"
	var dataTarget = "plant"

	plant := traceability.PlantEntityModel{
		PlantID:       uuid.MustParse(f.PlantID),
		OperatorID:    uuid.MustParse(f.OperatorID),
		PlantName:     f.OperatorName,
		PlantAddress:  f.OperatorAddress,
		OpenPlantID:   f.OpenOperatorID,
		GlobalPlantID: &f.GlobalOperatorID,
		DeletedAt: gorm.DeletedAt{
			Time:  time.Now(),
			Valid: true,
		},
		CreatedAt:     time.Now(),
		CreatedUserID: f.OperatorID,
		UpdatedAt:     time.Now(),
		UpdatedUserID: f.OperatorID,
	}

	expected := traceability.PlantModel{
		PlantID:      uuid.MustParse(f.OperatorID),
		OperatorID:   uuid.MustParse(f.OperatorID),
		PlantName:    f.OperatorName,
		PlantAddress: f.OperatorAddress,
		OpenPlantID:  f.OpenOperatorID,
		PlantAttribute: traceability.PlantAttribute{
			GlobalPlantID: &f.GlobalOperatorID,
		},
	}

	tests := []struct {
		name    string
		input   traceability.PlantModel
		receive traceability.PlantEntityModel
		expect  traceability.PlantModel
	}{
		{
			name:    "1-1. 200: 全項目応答(新規)",
			input:   f.NewPlantModel(false),
			receive: plant,
			expect:  expected,
		},
		{
			name:    "1-2. 200: 全項目応答(更新)",
			input:   f.NewPlantModel(true),
			receive: plant,
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
				if test.name == "1-1. 200: 全項目応答(新規)" {
					ouranosRepositoryMock.On("CreatePlant", mock.Anything).Return(test.receive, nil)
				} else {
					ouranosRepositoryMock.On("GetPlant", mock.Anything, mock.Anything).Return(test.receive, nil)
					ouranosRepositoryMock.On("UpdatePlant", mock.Anything).Return(test.receive, nil)
				}

				plantUsecase := usecase.NewPlantUsecase(ouranosRepositoryMock)

				actual, err := plantUsecase.PutPlant(test.input)
				if assert.NoError(t, err) {
					// 実際のレスポンスと期待されるレスポンスを比較
					// 順番が実行ごとに異なるため、順不同で中身を比較
					assert.Equal(t, test.expect.OperatorID, actual.OperatorID, f.AssertMessage)
					assert.Equal(t, test.expect.PlantName, actual.PlantName, f.AssertMessage)
					assert.Equal(t, test.expect.PlantAddress, actual.PlantAddress, f.AssertMessage)
					assert.Equal(t, test.expect.OpenPlantID, actual.OpenPlantID, f.AssertMessage)
					assert.Equal(t, test.expect.PlantAttribute.GlobalPlantID, actual.PlantAttribute.GlobalPlantID, f.AssertMessage)
				}
			},
		)
	}
}

// TestProjectUsecase_PutPlant_Abnormal
// Summary: This is abnormal test class which confirm the operation of API #3 PutPlant.
// Target: ouranos_plant_usecase.go
// TestPattern:
// [x] 2-1. 400: データ登録エラー(新規)
// [x] 2-2. 400: データ取得エラー(更新)
// [x] 2-3. 400: データ更新エラー(更新)
func TestProjectUsecase_PutPlant_Abnormal(tt *testing.T) {

	var method = "PUT"
	var endPoint = "/auth/v1/authInfo"
	var dataTarget = "plant"

	plant := traceability.PlantEntityModel{
		PlantID:       uuid.MustParse(f.PlantID),
		OperatorID:    uuid.MustParse(f.OperatorID),
		PlantName:     f.OperatorName,
		PlantAddress:  f.OperatorAddress,
		OpenPlantID:   f.OpenOperatorID,
		GlobalPlantID: &f.GlobalOperatorID,
		DeletedAt: gorm.DeletedAt{
			Time:  time.Now(),
			Valid: true,
		},
		CreatedAt:     time.Now(),
		CreatedUserID: f.OperatorID,
		UpdatedAt:     time.Now(),
		UpdatedUserID: f.OperatorID,
	}

	dsResGetError := fmt.Errorf("DB AccessError")

	tests := []struct {
		name         string
		input        traceability.PlantModel
		receive      traceability.PlantEntityModel
		receiveError error
		expect       error
	}{
		{
			name:         "2-1. 400: データ登録エラー(新規)",
			input:        f.NewPlantModel(false),
			receive:      traceability.PlantEntityModel{},
			receiveError: dsResGetError,
			expect:       dsResGetError,
		},
		{
			name:         "2-2. 400: データ取得エラー(更新)",
			input:        f.NewPlantModel(true),
			receive:      traceability.PlantEntityModel{},
			receiveError: dsResGetError,
			expect:       dsResGetError,
		},
		{
			name:         "2-3. 400: データ更新エラー(更新)",
			input:        f.NewPlantModel(true),
			receive:      plant,
			receiveError: dsResGetError,
			expect:       dsResGetError,
		},
		{
			name:    "2-4. 400: データ登録エラー(新規)",
			input:   f.NewPlantModel(false),
			receive: traceability.PlantEntityModel{},
			receiveError: &pgconn.PgError{
				Severity:         "ERROR",
				Code:             "23505",
				Message:          "unique_open_plant_id_operator_id",
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
			expect: fmt.Errorf("%s: %s is already exists.", "openPlantId", f.OpenPlantID),
		},
		{
			name:    "2-5. 400: データ登録エラー(新規)",
			input:   f.NewPlantModel(false),
			receive: traceability.PlantEntityModel{},
			receiveError: &pgconn.PgError{
				Severity:         "ERROR",
				Code:             "23504",
				Message:          "unique_open_plant_id_operator_id",
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
			expect: fmt.Errorf("ERROR: unique_open_plant_id_operator_id (SQLSTATE 23504)"),
		},
		{
			name:    "2-6. 400: データ更新エラー(更新)",
			input:   f.NewPlantModel(true),
			receive: traceability.PlantEntityModel{},
			receiveError: &pgconn.PgError{
				Severity:         "ERROR",
				Code:             "23505",
				Message:          "unique_global_plant_id_operator_id",
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
			expect: fmt.Errorf("%s: %s is already exists.", "globalPlantId", f.GlobalPlantID),
		},
		{
			name:    "2-7. 400: データ更新エラー(更新)",
			input:   f.NewPlantModel(true),
			receive: traceability.PlantEntityModel{},
			receiveError: &pgconn.PgError{
				Severity:         "ERROR",
				Code:             "23504",
				Message:          "unique_open_plant_id_operator_id",
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
			expect: fmt.Errorf("ERROR: unique_open_plant_id_operator_id (SQLSTATE 23504)"),
		},
		{
			name:         "2-8. 400: データ取得エラー(更新)",
			input:        f.NewPlantModel(true),
			receive:      traceability.PlantEntityModel{},
			receiveError: gorm.ErrRecordNotFound,
			expect:       common.NewCustomError(common.CustomErrorCode404, "record not found", nil, common.HTTPErrorSourceAuth),
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
				if test.name == "2-1. 400: データ登録エラー(新規)" {
					ouranosRepositoryMock.On("CreatePlant", mock.Anything).Return(traceability.PlantEntityModel{}, test.receiveError)
				} else if test.name == "2-2. 400: データ取得エラー(更新)" {
					ouranosRepositoryMock.On("GetPlant", mock.Anything, mock.Anything).Return(traceability.PlantEntityModel{}, test.receiveError)
				} else if test.name == "2-3. 400: データ更新エラー(更新)" {
					ouranosRepositoryMock.On("GetPlant", mock.Anything, mock.Anything).Return(test.receive, nil)
					ouranosRepositoryMock.On("UpdatePlant", mock.Anything).Return(traceability.PlantEntityModel{}, test.receiveError)
				} else if test.name == "2-4. 400: データ登録エラー(新規)" {
					ouranosRepositoryMock.On("CreatePlant", mock.Anything).Return(traceability.PlantEntityModel{}, test.receiveError)
				} else if test.name == "2-5. 400: データ登録エラー(新規)" {
					ouranosRepositoryMock.On("CreatePlant", mock.Anything).Return(traceability.PlantEntityModel{}, test.receiveError)
				} else if test.name == "2-6. 400: データ更新エラー(更新)" {
					ouranosRepositoryMock.On("GetPlant", mock.Anything, mock.Anything).Return(test.receive, nil)
					ouranosRepositoryMock.On("UpdatePlant", mock.Anything).Return(traceability.PlantEntityModel{}, test.receiveError)
				} else if test.name == "2-7. 400: データ更新エラー(更新)" {
					ouranosRepositoryMock.On("GetPlant", mock.Anything, mock.Anything).Return(test.receive, nil)
					ouranosRepositoryMock.On("UpdatePlant", mock.Anything).Return(traceability.PlantEntityModel{}, test.receiveError)
				} else if test.name == "2-8. 400: データ取得エラー(更新)" {
					ouranosRepositoryMock.On("GetPlant", mock.Anything, mock.Anything).Return(traceability.PlantEntityModel{}, test.receiveError)
				}
				plantUsecase := usecase.NewPlantUsecase(ouranosRepositoryMock)

				_, err := plantUsecase.PutPlant(test.input)
				if assert.Error(t, err) {
					// 実際のレスポンスと期待されるレスポンスを比較
					// 順番が実行ごとに異なるため、順不同で中身を比較
					assert.Equal(t, test.expect.Error(), err.Error())
				}
			},
		)
	}
}
