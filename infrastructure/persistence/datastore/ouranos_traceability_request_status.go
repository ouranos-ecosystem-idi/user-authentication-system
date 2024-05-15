package datastore

import (
	"fmt"
)

func (r *ouranosRepository) DeleteRequestStatusByTradeID(tradeID string) error {
	result := r.db.Unscoped().Table("request_status").Where("trade_id = ?", tradeID).Delete(nil)
	if result.Error != nil {
		return fmt.Errorf("failed to physically delete record from table request_status: %v", result.Error)
	}
	return nil
}
