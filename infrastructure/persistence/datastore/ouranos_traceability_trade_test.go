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
// Trades ListTradesByOperatorID テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1: 正常系：1件以上の場合
// [x] 1-2: 正常系：0件の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Trade_ListTradesByOperatorID(tt *testing.T) {

	tradeUUID1, _ := uuid.Parse("00000000-0000-0000-0000-000000000301")
	tradeUUID2, _ := uuid.Parse("00000000-0000-0000-0000-000000000302")
	tradeUUID3, _ := uuid.Parse("00000000-0000-0000-0000-000000000311")
	tradeUUID4, _ := uuid.Parse("00000000-0000-0000-0000-000000000312")
	operatorUUID1, _ := uuid.Parse("b39e6248-c888-56ca-d9d0-89de1b1adc8e")
	operatorUUID2, _ := uuid.Parse("15572d1c-ec13-0d78-7f92-dd4278871373")
	dsTraceUUID1, _ := uuid.Parse("00000000-0000-0000-0000-000000000204")
	dsTraceUUID2, _ := uuid.Parse("00000000-0000-0000-0000-000000000205")
	dsTraceUUID3, _ := uuid.Parse("00000000-0000-0000-0000-000000000214")
	dsTraceUUID4, _ := uuid.Parse("00000000-0000-0000-0000-000000000215")
	usTraceUUID1, _ := uuid.Parse("00000000-0000-0000-0000-000000000211")
	usTraceUUID2, _ := uuid.Parse("00000000-0000-0000-0000-000000000201")

	tests := []struct {
		name   string
		input  string
		expect traceability.TradeEntityModels
	}{
		{
			name:  "1-1: 正常系：1件以上の場合",
			input: "b39e6248-c888-56ca-d9d0-89de1b1adc8e",
			expect: traceability.TradeEntityModels{
				traceability.TradeEntityModel{
					TradeID:              &tradeUUID1,
					DownstreamOperatorID: operatorUUID1,
					UpstreamOperatorID:   &operatorUUID2,
					DownstreamTraceID:    dsTraceUUID1,
					UpstreamTraceID:      &usTraceUUID1,
					TradeDate:            common.StringPtr("2024-05-01T00:00:00Z"),
					DeletedAt:            gorm.DeletedAt{},
					CreatedAt:            time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
					CreatedUserID:        "seed",
					UpdatedAt:            time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
					UpdatedUserID:        "seed",
				},
				traceability.TradeEntityModel{
					TradeID:              &tradeUUID2,
					DownstreamOperatorID: operatorUUID1,
					UpstreamOperatorID:   &operatorUUID2,
					DownstreamTraceID:    dsTraceUUID2,
					UpstreamTraceID:      nil,
					TradeDate:            common.StringPtr("2024-05-01T00:00:00Z"),
					DeletedAt:            gorm.DeletedAt{},
					CreatedAt:            time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
					CreatedUserID:        "seed",
					UpdatedAt:            time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
					UpdatedUserID:        "seed",
				},
				traceability.TradeEntityModel{
					TradeID:              &tradeUUID3,
					DownstreamOperatorID: operatorUUID2,
					UpstreamOperatorID:   &operatorUUID1,
					DownstreamTraceID:    dsTraceUUID3,
					UpstreamTraceID:      &usTraceUUID2,
					TradeDate:            common.StringPtr("2024-05-01T00:00:00Z"),
					DeletedAt:            gorm.DeletedAt{},
					CreatedAt:            time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
					CreatedUserID:        "seed",
					UpdatedAt:            time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
					UpdatedUserID:        "seed",
				},
				traceability.TradeEntityModel{
					TradeID:              &tradeUUID4,
					DownstreamOperatorID: operatorUUID2,
					UpstreamOperatorID:   &operatorUUID1,
					DownstreamTraceID:    dsTraceUUID4,
					UpstreamTraceID:      nil,
					TradeDate:            common.StringPtr("2024-05-01T00:00:00Z"),
					DeletedAt:            gorm.DeletedAt{},
					CreatedAt:            time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
					CreatedUserID:        "seed",
					UpdatedAt:            time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
					UpdatedUserID:        "seed",
				},
			},
		},
		{
			name:   "1-2: 正常系：0件の場合",
			input:  "00000000-0000-0000-0000-000000000000",
			expect: traceability.TradeEntityModels{},
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
				actual, err := r.ListTradesByOperatorID(test.input)
				if assert.NoError(t, err) {
					assert.Equal(t, len(test.expect), len(actual))
					for i, data := range test.expect {
						assert.Equal(t, data.TradeID, actual[i].TradeID)
						assert.Equal(t, data.DownstreamOperatorID, actual[i].DownstreamOperatorID)
						assert.Equal(t, data.UpstreamOperatorID, actual[i].UpstreamOperatorID)
						assert.Equal(t, data.DownstreamTraceID, actual[i].DownstreamTraceID)
						assert.Equal(t, data.UpstreamTraceID, actual[i].UpstreamTraceID)
						assert.Equal(t, *data.TradeDate, *actual[i].TradeDate)
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
// Trades ListTradesByOperatorID テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1: 異常系：取得失敗の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Trade_ListTradesByOperatorID_Abnormal(tt *testing.T) {

	tests := []struct {
		name      string
		input     string
		dropQuery string
		expect    error
	}{
		{
			name:      "2-1: 異常系：取得失敗の場合",
			input:     "b39e6248-c888-56ca-d9d0-89de1b1adc8e",
			dropQuery: "DROP TABLE IF EXISTS trades",
			expect:    fmt.Errorf("no such table: trades"),
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
				_, err = r.ListTradesByOperatorID(test.input)
				if assert.Error(t, err) {
					assert.Equal(t, test.expect.Error(), err.Error())
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Trades DeleteTrade テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1: 正常系：削除成功の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Trade_DeleteTrade(tt *testing.T) {

	tests := []struct {
		name       string
		input      string
		checkQuery string
		before     int
		after      int
	}{
		{
			name:       "1-1: 正常系：削除成功の場合",
			input:      "00000000-0000-0000-0000-000000000301",
			checkQuery: "SELECT COUNT(*) FROM trades WHERE trade_id = ?",
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
				err = r.DeleteTrade(test.input)
				if assert.NoError(t, err) {
					db.Raw(test.checkQuery, test.input).Scan(&actualCount)
					assert.Equal(t, test.after, actualCount)
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// Trade DeleteTrade テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1: 異常系：削除失敗の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_Trade_DeleteTrade_Abnormal(tt *testing.T) {

	tests := []struct {
		name      string
		input     string
		dropQuery string
		expect    error
	}{
		{
			name:      "2-1: 異常系：削除失敗の場合",
			input:     "00000000-0000-0000-0000-000000000000",
			dropQuery: "DROP TABLE IF EXISTS trades",
			expect:    fmt.Errorf("failed to physically delete record from table trades: no such table: trades"),
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
				err = r.DeleteTrade(test.input)
				if assert.Error(t, err) {
					assert.Equal(t, test.expect.Error(), err.Error())
				}
			},
		)
	}
}
