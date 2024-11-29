package datastore

import (
	"fmt"

	"authenticator-backend/domain/model/traceability"
	"authenticator-backend/extension/logger"
)

func (r *ouranosRepository) ListTradesByOperatorID(operatorID string) (traceability.TradeEntityModels, error) {
	var result traceability.TradeEntityModels
	if err := r.db.Table("trades").Where("downstream_operator_id = ?", operatorID).Or("upstream_operator_id = ?", operatorID).Find(&result).Error; err != nil {
		logger.Set(nil).Errorf(err.Error())

		return traceability.TradeEntityModels{}, err
	}
	return result, nil

}

func (r *ouranosRepository) DeleteTrade(tradeID string) error {
	result := r.db.Unscoped().Table("trades").Where("trade_id = ?", tradeID).Delete(nil)
	if result.Error != nil {
		return fmt.Errorf("failed to physically delete record from table trades: %v", result.Error)
	}
	return nil
}
