package datastore

import (
	"authenticator-backend/domain/model/traceability"
	"authenticator-backend/extension/logger"
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
func (r *ouranosRepository) PutOperator(operatorEntityModel traceability.OperatorEntityModel) (traceability.OperatorEntityModel, error) {
	if err := r.db.Table("operators").Where("operator_id = ?", operatorEntityModel.OperatorID).Updates(&operatorEntityModel).Error; err != nil {
		logger.Set(nil).Errorf(err.Error())
		return traceability.OperatorEntityModel{}, err
	}
	return operatorEntityModel, nil
}
