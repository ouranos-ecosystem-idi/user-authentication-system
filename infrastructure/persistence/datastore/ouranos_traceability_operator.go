package datastore

import (
	"authenticator-backend/domain/model/traceability"
	"authenticator-backend/extension/logger"

	"gorm.io/gorm"
)

// GetOperator
// Summary: This is function which get Operator from data store.
// input: operatorID(string) Value of operatorID
// output: (traceability.OperatorEntityModel) OperatorEntityModel object
// output: (error) error object
func (r *ouranosRepository) GetOperator(operatorID string) (traceability.OperatorEntityModel, error) {
	var e traceability.OperatorEntityModel

	if err := r.db.Table("operators").Where("operator_id = ?", operatorID).First(&e).Error; err != nil {
		logger.Set(nil).Errorf(err.Error())
		return traceability.OperatorEntityModel{}, err
	}
	return e, nil
}

// GetOperatorByOpenOperatorID
// Summary: This is function which get Operator from data store with OpenOperatorID(.
// input: openOperatorID(string) Value of openOperatorID
// output: (traceability.OperatorEntityModel) OperatorEntityModel object
// output: (error) error object
func (r *ouranosRepository) GetOperatorByOpenOperatorID(openOperatorID string) (traceability.OperatorEntityModel, error) {
	var e traceability.OperatorEntityModel

	if err := r.db.Table("operators").Where("open_operator_id = ?", openOperatorID).First(&e).Error; err != nil {
		logger.Set(nil).Errorf(err.Error())
		return traceability.OperatorEntityModel{}, err
	}
	return e, nil
}

// GetOperators
// Summary: This is function which get Array of Operator from data store(.
// input: operatorIDs([]string) Array of operatorIDs
// output: (traceability.OperatorEntityModels) OperatorEntityModels object
// output: (error) error object
func (r *ouranosRepository) GetOperators(operatorIDs []string) (traceability.OperatorEntityModels, error) {
	var es traceability.OperatorEntityModels

	if err := r.db.Table("operators").Where("operator_id IN ?", operatorIDs).Find(&es).Error; err != nil {
		logger.Set(nil).Errorf(err.Error())
		return traceability.OperatorEntityModels{}, err
	}

	return es, nil
}

// PutOperator
// Summary: This is function which update Operator to data store.
// input: e(traceability.OperatorEntityModel) OperatorEntityModel object
// output: (traceability.OperatorEntityModel) OperatorEntityModel object
// output: (error) error object
func (r *ouranosRepository) PutOperator(e traceability.OperatorEntityModel) (traceability.OperatorEntityModel, error) {
	updateValues := map[string]interface{}{
		"operator_id":        e.OperatorID,
		"operator_name":      e.OperatorName,
		"operator_address":   e.OperatorAddress,
		"open_operator_id":   e.OpenOperatorID,
		"global_operator_id": e.GlobalOperatorID,
	}

	var updated traceability.OperatorEntityModel
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("operators").Where("operator_id = ?", e.OperatorID).Updates(updateValues).Error; err != nil {
			logger.Set(nil).Errorf(err.Error())
			return err
		}
		if err := tx.Table("operators").Where("operator_id = ?", e.OperatorID).First(&updated).Error; err != nil {
			logger.Set(nil).Errorf(err.Error())
			return err
		}
		return nil
	})

	if err != nil {
		return traceability.OperatorEntityModel{}, err
	}
	return updated, nil
}
