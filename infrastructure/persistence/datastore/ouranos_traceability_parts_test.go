package datastore_test

import (
	"authenticator-backend/domain/common"
	"authenticator-backend/domain/model/traceability"
	"authenticator-backend/infrastructure/persistence/datastore"
	testhelper "authenticator-backend/test/test_helper"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// /////////////////////////////////////////////////////////////////////////////////
// Parts ListParts テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1. 正常系：1件以上の場合
// [x] 1-2. 正常系：0件の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Parts_ListParts(tt *testing.T) {

	plantUUID, _ := uuid.Parse("00000000-0000-0000-0000-000000000101")
	partsUUID, _ := uuid.Parse("00000000-0000-0000-0000-000000000201")
	operatorUUID, _ := uuid.Parse("b39e6248-c888-56ca-d9d0-89de1b1adc8e")
	parentFlag := true
	var amountRequiredUnit traceability.AmountRequiredUnit = "kilogram"
	tests := []struct {
		name   string
		input  traceability.GetPartsModel
		expect traceability.PartsModels
	}{
		{
			name: "1-1: 正常系：1件以上の場合",
			input: traceability.GetPartsModel{
				OperatorID: "b39e6248-c888-56ca-d9d0-89de1b1adc8e",
				TraceID:    common.StringPtr("00000000-0000-0000-0000-000000000201"),
				PartsName:  common.StringPtr("製品A1"),
				PlantID:    common.StringPtr("00000000-0000-0000-0000-000000000101"),
				ParentFlag: &parentFlag,
				Limit:      1,
			},
			expect: traceability.PartsModels{
				&traceability.PartsModel{
					TraceID:            partsUUID,
					OperatorID:         operatorUUID,
					PlantID:            &plantUUID,
					PartsName:          "製品A1",
					SupportPartsName:   common.StringPtr("品番A1"),
					TerminatedFlag:     true,
					AmountRequired:     nil,
					AmountRequiredUnit: &amountRequiredUnit,
				},
			},
		},
		{
			name: "1-2: 正常系：0件の場合",
			input: traceability.GetPartsModel{
				OperatorID: "b39e6248-c888-56ca-d9d0-89de1b1adc8e",
				TraceID:    common.StringPtr("00000000-0000-0000-0000-000000000000"),
				PartsName:  common.StringPtr("製品A1"),
				PlantID:    common.StringPtr("00000000-0000-0000-0000-000000000101"),
				ParentFlag: &parentFlag,
				Limit:      1,
			},
			expect: traceability.PartsModels{},
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
				actual, err := r.ListParts(test.input)
				if assert.NoError(t, err) {
					assert.Equal(t, len(test.expect), len(actual))
					for i, data := range test.expect {
						assert.Equal(t, data.TraceID, actual[i].TraceID)
						assert.Equal(t, data.OperatorID, actual[i].OperatorID)
						assert.Equal(t, data.PlantID, actual[i].PlantID)
						assert.Equal(t, data.PartsName, actual[i].PartsName)
						assert.Equal(t, data.TerminatedFlag, actual[i].TerminatedFlag)
						assert.Equal(t, data.AmountRequired, actual[i].AmountRequired)
						assert.Equal(t, data.AmountRequiredUnit, actual[i].AmountRequiredUnit)
					}
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Parts ListParts テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1: 異常系：取得失敗の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Parts_ListParts_Abnormal(tt *testing.T) {

	parentFlag := true
	tests := []struct {
		name      string
		input     traceability.GetPartsModel
		dropQuery string
		expect    error
	}{
		{
			name: "2-1: 異常系：取得失敗の場合",
			input: traceability.GetPartsModel{
				OperatorID: "b39e6248-c888-56ca-d9d0-89de1b1adc8e",
				TraceID:    common.StringPtr("00000000-0000-0000-0000-000000000201"),
				PartsName:  common.StringPtr("製品A1"),
				PlantID:    common.StringPtr("00000000-0000-0000-0000-000000000101"),
				ParentFlag: &parentFlag,
				Limit:      1,
			},
			dropQuery: "DROP TABLE IF EXISTS parts",
			expect:    fmt.Errorf("no such table: parts"),
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
				_, err = r.ListParts(test.input)
				if assert.Error(t, err) {
					assert.Equal(t, test.expect.Error(), err.Error())
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Parts DeleteParts テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1: 正常系：削除成功の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Parts_DeleteParts(tt *testing.T) {

	tests := []struct {
		name       string
		input      string
		checkQuery string
		before     int
		after      int
	}{
		{
			name:       "1-1: 正常系：削除成功の場合",
			input:      "00000000-0000-0000-0000-000000000201",
			checkQuery: "SELECT COUNT(*) FROM parts WHERE trace_id = ?",
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
				var actualCount int
				db.Raw(test.checkQuery, test.input).Scan(&actualCount)
				assert.Equal(t, test.before, actualCount)
				err = r.DeleteParts(test.input)
				if assert.NoError(t, err) {
					db.Raw(test.checkQuery, test.input).Scan(&actualCount)
					assert.Equal(t, test.after, actualCount)
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Parts DeleteParts テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1: 異常系：削除失敗の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Parts_DeleteParts_Abnormal(tt *testing.T) {

	tests := []struct {
		name      string
		input     string
		dropQuery string
		expect    error
	}{
		{
			name:      "2-1: 異常系：削除失敗の場合",
			input:     "00000000-0000-0000-0000-000000000000",
			dropQuery: "DROP TABLE IF EXISTS parts",
			expect:    fmt.Errorf("failed to physically delete record from table parts: no such table: parts"),
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
				err = r.DeleteParts(test.input)
				if assert.Error(t, err) {
					assert.Equal(t, test.expect.Error(), err.Error())
				}
			},
		)
	}
}
