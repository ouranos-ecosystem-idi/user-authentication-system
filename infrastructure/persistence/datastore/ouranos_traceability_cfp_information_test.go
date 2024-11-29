package datastore_test

import (
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
// CFPInformation GetCFPInformation テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1. 200: 正常系：1件の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_CFPInformation_GetCFPInformation(tt *testing.T) {

	cfpUUID, _ := uuid.Parse("00000000-0000-0000-0000-000000000501")
	traceUUID, _ := uuid.Parse("00000000-0000-0000-0000-000000000201")
	ghgEmission := 0.1
	ter := 1.1
	ger := 2.1
	tir := 3.1

	tests := []struct {
		name   string
		input  string
		expect traceability.CfpEntityModel
	}{
		{
			name:  "1-1: 正常系：1件の場合",
			input: "00000000-0000-0000-0000-000000000201",
			expect: traceability.CfpEntityModel{
				CfpID:              &cfpUUID,
				TraceID:            traceUUID,
				GhgEmission:        &ghgEmission,
				GhgDeclaredUnit:    "kgCO2e/kilogram",
				CfpCertificateList: []string{},
				CfpType:            "preProduction",
				DqrType:            "preProcessing",
				TeR:                &ter,
				GeR:                &ger,
				TiR:                &tir,
				DeletedAt:          gorm.DeletedAt{},
				CreatedAt:          time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
				CreatedUserID:      "seed",
				UpdatedAt:          time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
				UpdatedUserID:      "seed",
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
				actual, err := r.GetCFPInformation(test.input)
				if assert.NoError(t, err) {
					assert.Equal(t, *test.expect.CfpID, *actual.CfpID)
					assert.Equal(t, test.expect.TraceID, actual.TraceID)
					assert.Equal(t, *test.expect.GhgEmission, *actual.GhgEmission)
					assert.Equal(t, test.expect.GhgDeclaredUnit, actual.GhgDeclaredUnit)
					assert.Equal(t, test.expect.CfpType, actual.CfpType)
					assert.Equal(t, test.expect.DqrType, actual.DqrType)
					assert.Equal(t, *test.expect.TeR, *actual.TeR)
					assert.Equal(t, *test.expect.GeR, *actual.GeR)
					assert.Equal(t, *test.expect.TiR, *actual.TiR)
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
// CFPInformation GetCFPInformation テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1: 異常系：0件の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_CFPInformation_GetCFPInformation_Abnormal(tt *testing.T) {

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
				_, err = r.GetCFPInformation(test.input)
				if assert.Error(t, err) {
					assert.Equal(t, test.expect.Error(), err.Error())
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// CFPInformation DeleteCFPInformation テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1: 正常系：削除成功の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_CFPInformation_DeleteCFPInformation(tt *testing.T) {

	tests := []struct {
		name       string
		input      string
		checkQuery string
		before     int
		after      int
	}{
		{
			name:       "1-1: 正常系：削除成功の場合",
			input:      "00000000-0000-0000-0000-000000000501",
			checkQuery: "SELECT COUNT(*) FROM cfp_infomation WHERE cfp_id = ?",
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
				err = r.DeleteCFPInformation(test.input)
				if assert.NoError(t, err) {
					db.Raw(test.checkQuery, test.input).Scan(&count)
					assert.Equal(t, count, test.after)
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// CFPInformation DeleteCFPInformation テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1: 異常系：削除失敗の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_CFPInformation_DeleteCFPInformation_Abnormal(tt *testing.T) {

	tests := []struct {
		name      string
		input     string
		dropQuery string
		expect    error
	}{
		{
			name:      "2-1: 異常系：削除失敗の場合",
			input:     "00000000-0000-0000-0000-000000000000",
			dropQuery: "DROP TABLE IF EXISTS cfp_infomation",
			expect:    fmt.Errorf("failed to physically delete record from table cfp_infomation: no such table: cfp_infomation"),
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
				err = r.DeleteCFPInformation(test.input)
				if assert.Error(t, err) {
					assert.Equal(t, test.expect.Error(), err.Error())
				}
			},
		)
	}
}
