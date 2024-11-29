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
// Operator GetOperator テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1. 200: 正常系：1件の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Operator_GetOperator(tt *testing.T) {

	operatorUUID, _ := uuid.Parse("b39e6248-c888-56ca-d9d0-89de1b1adc8e")

	tests := []struct {
		name   string
		input  string
		expect traceability.OperatorEntityModel
	}{
		{
			name:  "1-1: 正常系：1件の場合",
			input: "b39e6248-c888-56ca-d9d0-89de1b1adc8e",
			expect: traceability.OperatorEntityModel{
				OperatorID:        operatorUUID,
				OperatorName:      "A社",
				OperatorAddress:   "東京都渋谷区xx",
				OpenOperatorID:    "1234567890123",
				GlobalOperatorID:  common.StringPtr("1234ABCD5678EFGH0123"),
				DeletedAt:         gorm.DeletedAt{},
				CreatedAt:         time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
				CreatedOperatorID: "seed",
				UpdatedAt:         time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
				UpdatedOperatorID: "seed",
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
				actual, err := r.GetOperator(test.input)
				if assert.NoError(t, err) {
					assert.Equal(t, test.expect.OperatorID, actual.OperatorID)
					assert.Equal(t, test.expect.OperatorName, actual.OperatorName)
					assert.Equal(t, test.expect.OperatorAddress, actual.OperatorAddress)
					assert.Equal(t, test.expect.OpenOperatorID, actual.OpenOperatorID)
					assert.Equal(t, test.expect.GlobalOperatorID, actual.GlobalOperatorID)
					assert.Equal(t, test.expect.DeletedAt, actual.DeletedAt)
					assert.Equal(t, test.expect.CreatedAt, actual.CreatedAt)
					//assert.Equal(t, test.expect.CreatedOperatorID, actual.CreatedOperatorID)
					assert.Equal(t, test.expect.UpdatedAt, actual.UpdatedAt)
					//assert.Equal(t, test.expect.UpdatedOperatorID, actual.UpdatedOperatorID)
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Operator GetOperator テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1: 異常系：0件の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Operator_GetOperator_Abnormal(tt *testing.T) {

	tests := []struct {
		name   string
		input  string
		expect error
	}{
		{
			name:   "2-1: 異常系：0件の場合",
			input:  "00000000-0000-0000-0000-000000000000",
			expect: fmt.Errorf("record not found"),
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
				_, err = r.GetOperator(test.input)
				if assert.Error(t, err) {
					assert.Equal(t, test.expect.Error(), err.Error())
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Operator GetOperatorByOpenOperatorID テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1. 200: 正常系：1件の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Operator_GetOperatorByOpenOperatorID(tt *testing.T) {

	operatorUUID, _ := uuid.Parse("b39e6248-c888-56ca-d9d0-89de1b1adc8e")

	tests := []struct {
		name   string
		input  string
		expect traceability.OperatorEntityModel
	}{
		{
			name:  "1-1: 正常系：1件の場合",
			input: "1234567890123",
			expect: traceability.OperatorEntityModel{
				OperatorID:        operatorUUID,
				OperatorName:      "A社",
				OperatorAddress:   "東京都渋谷区xx",
				OpenOperatorID:    "1234567890123",
				GlobalOperatorID:  common.StringPtr("1234ABCD5678EFGH0123"),
				DeletedAt:         gorm.DeletedAt{},
				CreatedAt:         time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
				CreatedOperatorID: "seed",
				UpdatedAt:         time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
				UpdatedOperatorID: "seed",
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
				actual, err := r.GetOperatorByOpenOperatorID(test.input)
				if assert.NoError(t, err) {
					assert.Equal(t, test.expect.OperatorID, actual.OperatorID)
					assert.Equal(t, test.expect.OperatorName, actual.OperatorName)
					assert.Equal(t, test.expect.OperatorAddress, actual.OperatorAddress)
					assert.Equal(t, test.expect.OpenOperatorID, actual.OpenOperatorID)
					assert.Equal(t, test.expect.GlobalOperatorID, actual.GlobalOperatorID)
					assert.Equal(t, test.expect.DeletedAt, actual.DeletedAt)
					assert.Equal(t, test.expect.CreatedAt, actual.CreatedAt)
					//assert.Equal(t, test.expect.CreatedOperatorID, actual.CreatedOperatorID)
					assert.Equal(t, test.expect.UpdatedAt, actual.UpdatedAt)
					//assert.Equal(t, test.expect.UpdatedOperatorID, actual.UpdatedOperatorID)
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Operator GetOperatorByOpenOperatorID テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1: 異常系：0件の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Operator_GetOperatorByOpenOperatorID_Abnormal(tt *testing.T) {

	tests := []struct {
		name   string
		input  string
		expect error
	}{
		{
			name:   "2-1: 異常系：0件の場合",
			input:  "00000000-0000-0000-0000-000000000000",
			expect: fmt.Errorf("record not found"),
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
				_, err = r.GetOperatorByOpenOperatorID(test.input)
				if assert.Error(t, err) {
					assert.Equal(t, test.expect.Error(), err.Error())
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Operator GetOperators テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1. 正常系：1件以上の場合
// [x] 1-2. 正常系：0件の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Operator_GetOperators(tt *testing.T) {

	operatorUUID1, _ := uuid.Parse("b39e6248-c888-56ca-d9d0-89de1b1adc8e")
	operatorUUID2, _ := uuid.Parse("15572d1c-ec13-0d78-7f92-dd4278871373")

	tests := []struct {
		name   string
		input  []string
		expect traceability.OperatorEntityModels
	}{
		{
			name:  "1-1: 正常系：1件以上の場合",
			input: []string{"b39e6248-c888-56ca-d9d0-89de1b1adc8e", "15572d1c-ec13-0d78-7f92-dd4278871373"},
			expect: traceability.OperatorEntityModels{
				&traceability.OperatorEntityModel{
					OperatorID:        operatorUUID2,
					OperatorName:      "B社",
					OperatorAddress:   "東京都渋谷区xx",
					OpenOperatorID:    "1234567890124",
					GlobalOperatorID:  common.StringPtr("1234ABCD5678EFGH0124"),
					DeletedAt:         gorm.DeletedAt{},
					CreatedAt:         time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
					CreatedOperatorID: "seed",
					UpdatedAt:         time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
					UpdatedOperatorID: "seed",
				},
				&traceability.OperatorEntityModel{
					OperatorID:        operatorUUID1,
					OperatorName:      "A社",
					OperatorAddress:   "東京都渋谷区xx",
					OpenOperatorID:    "1234567890123",
					GlobalOperatorID:  common.StringPtr("1234ABCD5678EFGH0123"),
					DeletedAt:         gorm.DeletedAt{},
					CreatedAt:         time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
					CreatedOperatorID: "seed",
					UpdatedAt:         time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
					UpdatedOperatorID: "seed",
				},
			},
		},
		{
			name:   "1-2: 正常系：0件の場合",
			input:  []string{"00000000-0000-0000-0000-000000000000"},
			expect: traceability.OperatorEntityModels{},
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
				actual, err := r.GetOperators(test.input)
				if assert.NoError(t, err) {
					assert.Equal(t, len(actual), len(test.expect))
					for i, data := range test.expect {
						assert.Equal(t, data.OperatorID, actual[i].OperatorID)
						assert.Equal(t, data.OperatorName, actual[i].OperatorName)
						assert.Equal(t, data.OperatorAddress, actual[i].OperatorAddress)
						assert.Equal(t, data.OpenOperatorID, actual[i].OpenOperatorID)
						assert.Equal(t, data.GlobalOperatorID, actual[i].GlobalOperatorID)
						assert.Equal(t, data.DeletedAt, actual[i].DeletedAt)
						assert.Equal(t, data.CreatedAt, actual[i].CreatedAt)
						//assert.Equal(t, data.CreatedOperatorID, actual[i].CreatedOperatorID)
						assert.Equal(t, data.UpdatedAt, actual[i].UpdatedAt)
						//assert.Equal(t, data.UpdatedOperatorID, actual[i].UpdatedOperatorID)
					}
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Operator TestProjectRepository_Operator_GetOperators テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1: 異常系：削除失敗の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Operator_GetOperators_Abnormal(tt *testing.T) {

	tests := []struct {
		name      string
		input     []string
		dropQuery string
		expect    error
	}{
		{
			name:      "2-1: 異常系：削除失敗の場合",
			input:     []string{"00000000-0000-0000-0000-000000000000"},
			dropQuery: "DROP TABLE IF EXISTS operators",
			expect:    fmt.Errorf("no such table: operators"),
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
				_, err = r.GetOperators(test.input)
				if assert.Error(t, err) {
					assert.Equal(t, test.expect.Error(), err.Error())
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Operator PutOperator テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1: 正常系：更新成功の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Operator_PutOperator(tt *testing.T) {

	operatorUUID, _ := uuid.Parse("b39e6248-c888-56ca-d9d0-89de1b1adc8e")

	tests := []struct {
		name   string
		input  traceability.OperatorEntityModel
		expect traceability.OperatorEntityModel
	}{
		{
			name: "1-1: 正常系：更新成功の場合",
			input: traceability.OperatorEntityModel{
				OperatorID:       operatorUUID,
				OperatorName:     "B社更新",
				OperatorAddress:  "東京都渋谷区xx更新",
				OpenOperatorID:   "1234567890125",
				GlobalOperatorID: common.StringPtr("1234ABCD5678EFGH0125"),
			},
			expect: traceability.OperatorEntityModel{
				OperatorID:        operatorUUID,
				OperatorName:      "B社更新",
				OperatorAddress:   "東京都渋谷区xx更新",
				OpenOperatorID:    "1234567890125",
				GlobalOperatorID:  common.StringPtr("1234ABCD5678EFGH0125"),
				DeletedAt:         gorm.DeletedAt{},
				CreatedAt:         time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
				CreatedOperatorID: "seed",
				UpdatedAt:         time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
				UpdatedOperatorID: "seed",
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
				actual, err := r.PutOperator(test.input)
				if assert.NoError(t, err) {
					assert.Equal(t, test.expect.OperatorID, actual.OperatorID)
					assert.Equal(t, test.expect.OperatorName, actual.OperatorName)
					assert.Equal(t, test.expect.OperatorAddress, actual.OperatorAddress)
					assert.Equal(t, test.expect.OpenOperatorID, actual.OpenOperatorID)
					assert.Equal(t, test.expect.GlobalOperatorID, actual.GlobalOperatorID)
					assert.Equal(t, test.expect.DeletedAt, actual.DeletedAt)
					assert.Equal(t, test.expect.CreatedAt, actual.CreatedAt)
					//assert.Equal(t, test.expect.CreatedOperatorID, actual.CreatedOperatorID)
					assert.Equal(t, test.expect.UpdatedAt, actual.UpdatedAt)
					//assert.Equal(t, test.expect.UpdatedOperatorID, actual.UpdatedOperatorID)
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Operator TestProjectRepository_Operator_GetOperators テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1: 異常系：削除失敗の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Operator_PutOperator_Abnormal(tt *testing.T) {

	operatorUUID, _ := uuid.Parse("b39e6248-c888-56ca-d9d0-89de1b1adc8e")

	tests := []struct {
		name      string
		input     traceability.OperatorEntityModel
		dropQuery string
		expect    error
	}{
		{
			name: "2-1: 異常系：削除失敗の場合",
			input: traceability.OperatorEntityModel{
				OperatorID:       operatorUUID,
				OperatorName:     "B社更新",
				OperatorAddress:  "東京都渋谷区xx更新",
				OpenOperatorID:   "1234567890125",
				GlobalOperatorID: common.StringPtr("1234ABCD5678EFGH0125"),
			},
			dropQuery: "DROP TABLE IF EXISTS operators",
			expect:    fmt.Errorf("no such table: operators"),
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
				_, err = r.PutOperator(test.input)
				if assert.Error(t, err) {
					assert.Equal(t, test.expect.Error(), err.Error())
				}
			},
		)
	}
}
