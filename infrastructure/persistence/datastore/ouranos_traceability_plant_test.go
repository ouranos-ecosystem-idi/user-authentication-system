package datastore_test

import (
	"authenticator-backend/domain/common"
	"authenticator-backend/domain/model/traceability"
	"authenticator-backend/infrastructure/persistence/datastore"
	testhelper "authenticator-backend/test/test_helper"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// /////////////////////////////////////////////////////////////////////////////////
// Plant CreatePlant テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1: 正常系：生成成功の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Plant_CreatePlant(tt *testing.T) {

	plantUUID, _ := uuid.Parse("00000000-0000-0000-0000-000000000103")
	operatorUUID, _ := uuid.Parse("b39e6248-c888-56ca-d9d0-89de1b1adc8e")

	tests := []struct {
		name   string
		input  traceability.PlantEntityModel
		expect traceability.PlantEntityModel
	}{
		{
			name: "1-1: 正常系：生成成功の場合",
			input: traceability.PlantEntityModel{
				PlantID:       plantUUID,
				OperatorID:    operatorUUID,
				PlantName:     "事業所A3",
				PlantAddress:  "工場A3",
				OpenPlantID:   "000000000003",
				GlobalPlantID: common.StringPtr("000000000003"),
				DeletedAt:     gorm.DeletedAt{},
				CreatedAt:     time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
				CreatedUserID: "test",
				UpdatedAt:     time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
				UpdatedUserID: "test",
			},
			expect: traceability.PlantEntityModel{
				PlantID:       plantUUID,
				OperatorID:    operatorUUID,
				PlantName:     "事業所A3",
				PlantAddress:  "工場A3",
				OpenPlantID:   "000000000003",
				GlobalPlantID: common.StringPtr("000000000003"),
				DeletedAt:     gorm.DeletedAt{},
				CreatedAt:     time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
				CreatedUserID: "test",
				UpdatedAt:     time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
				UpdatedUserID: "test",
			},
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				db, err := testhelper.NewMockDB()
				if err != nil {
					assert.Fail(t, err.Error())
				}
				r := datastore.NewOuranosRepository(db)
				actual, err := r.CreatePlant(test.input)
				if assert.NoError(t, err) {
					assert.Equal(t, test.expect.PlantID, actual.PlantID)
					assert.Equal(t, test.expect.OperatorID, actual.OperatorID)
					assert.Equal(t, test.expect.PlantName, actual.PlantName)
					assert.Equal(t, test.expect.PlantAddress, actual.PlantAddress)
					assert.Equal(t, test.expect.OpenPlantID, actual.OpenPlantID)
					assert.Equal(t, test.expect.GlobalPlantID, actual.GlobalPlantID)
					assert.Equal(t, test.expect.DeletedAt, actual.DeletedAt)
					assert.Equal(t, test.expect.CreatedAt, actual.CreatedAt)
					assert.Equal(t, test.expect.CreatedUserID, actual.CreatedUserID)
					assert.Equal(t, test.expect.UpdatedAt, actual.UpdatedAt)
					assert.Equal(t, test.expect.UpdatedUserID, actual.UpdatedUserID)
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Plant CreatePlant テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1: 異常系：削除失敗の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Plant_CreatePlant_Abnormal(tt *testing.T) {

	plantUUID, _ := uuid.Parse("00000000-0000-0000-0000-000000000103")
	operatorUUID, _ := uuid.Parse("b39e6248-c888-56ca-d9d0-89de1b1adc8e")

	tests := []struct {
		name      string
		input     traceability.PlantEntityModel
		dropQuery string
		expect    error
	}{
		{
			name: "2-1: 異常系：削除失敗の場合",
			input: traceability.PlantEntityModel{
				PlantID:       plantUUID,
				OperatorID:    operatorUUID,
				PlantName:     "事業所A3",
				PlantAddress:  "工場A3",
				OpenPlantID:   "000000000003",
				GlobalPlantID: common.StringPtr("000000000003"),
				DeletedAt:     gorm.DeletedAt{},
				CreatedAt:     time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
				CreatedUserID: "test",
				UpdatedAt:     time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
				UpdatedUserID: "test",
			},
			dropQuery: "DROP TABLE IF EXISTS plants",
			expect:    fmt.Errorf("no such table: plants"),
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				db, err := testhelper.NewMockDB()
				if err != nil {
					assert.Fail(t, "Errors occured by creating Mock DB")
				}
				err = db.Exec(test.dropQuery).Error
				if err != nil {
					assert.Fail(t, "Errors occured by deleting DB")
				}
				r := datastore.NewOuranosRepository(db)
				_, err = r.CreatePlant(test.input)
				if assert.Error(t, err) {
					assert.Equal(t, test.expect.Error(), err.Error())
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Plant GetPlant テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1. 200: 正常系：1件の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Plant_GetPlant(tt *testing.T) {

	plantUUID, _ := uuid.Parse("00000000-0000-0000-0000-000000000101")
	operatorUUID, _ := uuid.Parse("b39e6248-c888-56ca-d9d0-89de1b1adc8e")

	tests := []struct {
		name            string
		inputPlantID    string
		inputOperatorID string
		expect          traceability.PlantEntityModel
	}{
		{
			name:            "1-1: 正常系：1件の場合",
			inputPlantID:    "00000000-0000-0000-0000-000000000101",
			inputOperatorID: "b39e6248-c888-56ca-d9d0-89de1b1adc8e",
			expect: traceability.PlantEntityModel{
				PlantID:       plantUUID,
				OperatorID:    operatorUUID,
				PlantName:     "事業所A1",
				PlantAddress:  "工場A1",
				OpenPlantID:   "000000000001",
				GlobalPlantID: common.StringPtr("000000000001"),
				DeletedAt:     gorm.DeletedAt{},
				CreatedAt:     time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
				CreatedUserID: "seed",
				UpdatedAt:     time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
				UpdatedUserID: "seed",
			},
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				db, err := testhelper.NewMockDB()
				if err != nil {
					assert.Fail(t, err.Error())
				}
				r := datastore.NewOuranosRepository(db)
				actual, err := r.GetPlant(test.inputOperatorID, test.inputPlantID)
				if assert.NoError(t, err) {
					assert.Equal(t, test.expect.PlantID, actual.PlantID)
					assert.Equal(t, test.expect.OperatorID, actual.OperatorID)
					assert.Equal(t, test.expect.PlantName, actual.PlantName)
					assert.Equal(t, test.expect.PlantAddress, actual.PlantAddress)
					assert.Equal(t, test.expect.OpenPlantID, actual.OpenPlantID)
					assert.Equal(t, test.expect.GlobalPlantID, actual.GlobalPlantID)
					assert.Equal(t, test.expect.DeletedAt, actual.DeletedAt)
					assert.Equal(t, test.expect.CreatedAt, actual.CreatedAt)
					assert.Equal(t, test.expect.CreatedUserID, actual.CreatedUserID)
					assert.Equal(t, test.expect.UpdatedAt, actual.UpdatedAt)
					assert.Equal(t, test.expect.UpdatedUserID, actual.UpdatedUserID)
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Plant GetPlant テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1: 異常系：0件の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Plant_GetPlant_Abnormal(tt *testing.T) {

	tests := []struct {
		name            string
		inputPlantID    string
		inputOperatorID string
		expect          error
	}{
		{
			name:            "2-1: 異常系：0件の場合",
			inputPlantID:    "00000000-0000-0000-0000-000000000000",
			inputOperatorID: "b39e6248-c888-56ca-d9d0-89de1b1adc8e",
			expect:          fmt.Errorf("record not found"),
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				db, err := testhelper.NewMockDB()
				if err != nil {
					assert.Fail(t, err.Error())
				}
				r := datastore.NewOuranosRepository(db)
				_, err = r.GetPlant(test.inputOperatorID, test.inputPlantID)
				if assert.Error(t, err) {
					assert.Equal(t, test.expect.Error(), err.Error())
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Plant ListPlantsByOperatorID テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1. 正常系：1件以上の場合
// [x] 1-2. 正常系：0件の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Plant_ListPlantsByOperatorID(tt *testing.T) {

	plantUUID1, _ := uuid.Parse("00000000-0000-0000-0000-000000000101")
	plantUUID2, _ := uuid.Parse("00000000-0000-0000-0000-000000000102")
	operatorUUID, _ := uuid.Parse("b39e6248-c888-56ca-d9d0-89de1b1adc8e")

	tests := []struct {
		name   string
		input  string
		expect traceability.PlantEntityModels
	}{
		{
			name:  "1-1: 正常系：1件以上の場合",
			input: "b39e6248-c888-56ca-d9d0-89de1b1adc8e",
			expect: traceability.PlantEntityModels{
				traceability.PlantEntityModel{
					PlantID:       plantUUID1,
					OperatorID:    operatorUUID,
					PlantName:     "事業所A1",
					PlantAddress:  "工場A1",
					OpenPlantID:   "000000000001",
					GlobalPlantID: common.StringPtr("000000000001"),
					DeletedAt:     gorm.DeletedAt{},
					CreatedAt:     time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
					CreatedUserID: "seed",
					UpdatedAt:     time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
					UpdatedUserID: "seed",
				},
				traceability.PlantEntityModel{
					PlantID:       plantUUID2,
					OperatorID:    operatorUUID,
					PlantName:     "事業所A2",
					PlantAddress:  "工場A2",
					OpenPlantID:   "000000000002",
					GlobalPlantID: common.StringPtr("000000000002"),
					DeletedAt:     gorm.DeletedAt{},
					CreatedAt:     time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
					CreatedUserID: "seed",
					UpdatedAt:     time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
					UpdatedUserID: "seed",
				},
			},
		},
		{
			name:   "1-2: 正常系：0件の場合",
			input:  "00000000-0000-0000-0000-000000000000",
			expect: traceability.PlantEntityModels{},
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				db, err := testhelper.NewMockDB()
				if err != nil {
					assert.Fail(t, err.Error())
				}
				r := datastore.NewOuranosRepository(db)
				actual, err := r.ListPlantsByOperatorID(test.input)
				if assert.NoError(t, err) {
					assert.Equal(t, len(test.expect), len(actual))
					for i, data := range test.expect {
						assert.Equal(t, data.PlantID, actual[i].PlantID)
						assert.Equal(t, data.OperatorID, actual[i].OperatorID)
						assert.Equal(t, data.PlantName, actual[i].PlantName)
						assert.Equal(t, data.PlantAddress, actual[i].PlantAddress)
						assert.Equal(t, data.OpenPlantID, actual[i].OpenPlantID)
						assert.Equal(t, data.GlobalPlantID, actual[i].GlobalPlantID)
						assert.Equal(t, data.DeletedAt, actual[i].DeletedAt)
						assert.Equal(t, data.CreatedAt, actual[i].CreatedAt)
						assert.Equal(t, data.CreatedUserID, actual[i].CreatedUserID)
						assert.Equal(t, data.UpdatedAt, actual[i].UpdatedAt)
						assert.Equal(t, data.UpdatedUserID, actual[i].UpdatedUserID)
					}
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Plant ListPlantsByOperatorID テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1: 異常系：取得失敗の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Plant_ListPlantsByOperatorID_Abnormal(tt *testing.T) {

	tests := []struct {
		name      string
		input     string
		dropQuery string
		expect    error
	}{
		{
			name:      "2-1: 異常系：取得失敗の場合",
			input:     "b39e6248-c888-56ca-d9d0-89de1b1adc8e",
			dropQuery: "DROP TABLE IF EXISTS plants",
			expect:    fmt.Errorf("no such table: plants"),
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				db, err := testhelper.NewMockDB()
				if err != nil {
					assert.Fail(t, "Errors occured by creating Mock DB")
				}
				err = db.Exec(test.dropQuery).Error
				if err != nil {
					assert.Fail(t, "Errors occured by deleting DB")
				}
				r := datastore.NewOuranosRepository(db)
				_, err = r.ListPlantsByOperatorID(test.input)
				if assert.Error(t, err) {
					assert.Equal(t, test.expect.Error(), err.Error())
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Plant UpdatePlant テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1: 正常系：更新成功の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Plant_UpdatePlant(tt *testing.T) {

	plantUUID, _ := uuid.Parse("00000000-0000-0000-0000-000000000101")
	operatorUUID, _ := uuid.Parse("b39e6248-c888-56ca-d9d0-89de1b1adc8e")

	tests := []struct {
		name   string
		input  traceability.PlantEntityModel
		expect traceability.PlantEntityModel
	}{
		{
			name: "1-1: 正常系：更新成功の場合",
			input: traceability.PlantEntityModel{
				PlantID:       plantUUID,
				OperatorID:    operatorUUID,
				PlantName:     "事業所A1更新",
				PlantAddress:  "工場A1更新",
				OpenPlantID:   "000000000021",
				GlobalPlantID: common.StringPtr("000000000021"),
			},
			expect: traceability.PlantEntityModel{
				PlantID:       plantUUID,
				OperatorID:    operatorUUID,
				PlantName:     "事業所A1更新",
				PlantAddress:  "工場A1更新",
				OpenPlantID:   "000000000021",
				GlobalPlantID: common.StringPtr("000000000021"),
				DeletedAt:     gorm.DeletedAt{},
				CreatedAt:     time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
				CreatedUserID: "seed",
				UpdatedAt:     time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
				UpdatedUserID: "seed",
			},
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				db, err := testhelper.NewMockDB()
				if err != nil {
					assert.Fail(t, err.Error())
				}
				r := datastore.NewOuranosRepository(db)
				actual, err := r.UpdatePlant(test.input)
				if assert.NoError(t, err) {
					assert.Equal(t, actual.PlantID, test.expect.PlantID)
					assert.Equal(t, actual.OperatorID, test.expect.OperatorID)
					assert.Equal(t, actual.PlantName, test.expect.PlantName)
					assert.Equal(t, actual.PlantAddress, test.expect.PlantAddress)
					assert.Equal(t, actual.OpenPlantID, test.expect.OpenPlantID)
					assert.Equal(t, actual.GlobalPlantID, test.expect.GlobalPlantID)
					assert.Equal(t, actual.DeletedAt, test.expect.DeletedAt)
					assert.Equal(t, actual.CreatedAt, test.expect.CreatedAt)
					assert.Equal(t, actual.CreatedUserID, test.expect.CreatedUserID)
					assert.Equal(t, actual.UpdatedAt, test.expect.UpdatedAt)
					assert.Equal(t, actual.UpdatedUserID, test.expect.UpdatedUserID)
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Plant UpdatePlant テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1: 異常系：更新失敗の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Plant_UpdatePlant_Abnormal(tt *testing.T) {

	plantUUID, _ := uuid.Parse("00000000-0000-0000-0000-000000000101")
	operatorUUID, _ := uuid.Parse("b39e6248-c888-56ca-d9d0-89de1b1adc8e")

	tests := []struct {
		name      string
		input     traceability.PlantEntityModel
		dropQuery string
		expect    error
	}{
		{
			name: "1-1: 正常系：更新成功の場合",
			input: traceability.PlantEntityModel{
				PlantID:       plantUUID,
				OperatorID:    operatorUUID,
				PlantName:     "事業所A1更新",
				PlantAddress:  "工場A1更新",
				OpenPlantID:   "000000000021",
				GlobalPlantID: common.StringPtr("000000000021"),
			},
			dropQuery: "DROP TABLE IF EXISTS plants",
			expect:    fmt.Errorf("no such table: plants"),
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				db, err := testhelper.NewMockDB()
				if err != nil {
					assert.Fail(t, "Errors occured by creating Mock DB")
				}
				err = db.Exec(test.dropQuery).Error
				if err != nil {
					assert.Fail(t, "Errors occured by deleting DB")
				}
				r := datastore.NewOuranosRepository(db)
				_, err = r.UpdatePlant(test.input)
				if assert.Error(t, err) {
					assert.Equal(t, test.expect.Error(), err.Error())
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Plant DeletePlant テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1: 正常系：削除成功の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Plant_DeletePlant(tt *testing.T) {

	tests := []struct {
		name       string
		input      string
		checkQuery string
		before     int
		after      int
	}{
		{
			name:       "1-1: 正常系：削除成功の場合",
			input:      "00000000-0000-0000-0000-000000000101",
			checkQuery: "SELECT COUNT(*) FROM plants WHERE plant_id = ?",
			before:     1,
			after:      0,
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				db, err := testhelper.NewMockDB()
				if err != nil {
					assert.Fail(t, "Errors occured by creating Mock DB")
				}
				r := datastore.NewOuranosRepository(db)
				var count int
				db.Raw(test.checkQuery, test.input).Scan(&count)
				assert.Equal(t, count, test.before)
				err = r.DeletePlant(test.input)
				if assert.NoError(t, err) {
					db.Raw(test.checkQuery, test.input).Scan(&count)
					assert.Equal(t, count, test.after)
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Plant DeletePlant テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1: 異常系：削除失敗の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Plant_DeletePlant_Abnormal(tt *testing.T) {

	tests := []struct {
		name      string
		input     string
		dropQuery string
		expect    error
	}{
		{
			name:      "2-1: 異常系：削除失敗の場合",
			input:     "00000000-0000-0000-0000-000000000000",
			dropQuery: "DROP TABLE IF EXISTS plants",
			expect:    fmt.Errorf("failed to physically delete record from table plants : no such table: plants"),
		},
	}

	for _, test := range tests {
		test := test
		tt.Run(
			test.name,
			func(t *testing.T) {
				db, err := testhelper.NewMockDB()
				if err != nil {
					assert.Fail(t, "Errors occured by creating Mock DB")
				}
				err = db.Exec(test.dropQuery).Error
				if err != nil {
					assert.Fail(t, "Errors occured by deleting DB")
				}
				r := datastore.NewOuranosRepository(db)
				err = r.DeletePlant(test.input)
				if assert.Error(t, err) {
					assert.Equal(t, test.expect.Error(), err.Error())
				}
			},
		)
	}
}
