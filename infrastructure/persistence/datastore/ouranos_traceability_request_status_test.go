package datastore_test

import (
	"authenticator-backend/infrastructure/persistence/datastore"
	testhelper "authenticator-backend/test/test_helper"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// /////////////////////////////////////////////////////////////////////////////////
// RequestStatus DeleteRequestStatusByTradeID テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 1-1: 正常系：削除成功の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_RequestStatus_DeleteRequestStatusByTradeID(tt *testing.T) {

	tests := []struct {
		name       string
		input      string
		checkQuery string
		before     int
		after      int
	}{
		{
			name:       "1-1: 正常系：削除成功の場合",
			input:      "00000000-0000-0000-0000-000000000401",
			checkQuery: "SELECT COUNT(*) FROM request_status WHERE status_id = ?",
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
				err = r.DeleteRequestStatusByTradeID(test.input)
				if assert.NoError(t, err) {
					db.Raw(test.checkQuery, test.input).Scan(&actualCount)
					assert.Equal(t, test.after, actualCount)
				}
			},
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////////
// RequestStatus DeleteRequestStatusByTradeID テストケース
// /////////////////////////////////////////////////////////////////////////////////
// [x] 2-1: 異常系：削除失敗の場合
// /////////////////////////////////////////////////////////////////////////////////
func TestProjectRepository_RequestStatus_DeleteRequestStatusByTradeID_Abnormal(tt *testing.T) {

	tests := []struct {
		name      string
		input     string
		dropQuery string
		expect    error
	}{
		{
			name:      "2-1: 異常系：削除失敗の場合",
			input:     "00000000-0000-0000-0000-000000000000",
			dropQuery: "DROP TABLE IF EXISTS request_status",
			expect:    fmt.Errorf("failed to physically delete record from table request_status: no such table: request_status"),
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
				err = r.DeleteRequestStatusByTradeID(test.input)
				if assert.Error(t, err) {
					assert.Equal(t, test.expect.Error(), err.Error())
				}
			},
		)
	}
}
