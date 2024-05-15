package datastore

import (
	"fmt"

	"authenticator-backend/domain/model/traceability"

	"go.uber.org/zap"
)

func (r *ouranosRepository) GetCFPInformation(traceID string) (traceability.CfpEntityModel, error) {
	var result traceability.CfpEntityModel

	if err := r.db.Table("cfp_infomation").Where("trace_id = ?", traceID).First(&result).Error; err != nil {
		zap.S().Errorf(err.Error())
		return traceability.CfpEntityModel{}, err
	}
	return result, nil

}
func (r *ouranosRepository) DeleteCFPInformation(cfpID string) error {
	result := r.db.Unscoped().Table("cfp_infomation").Where("cfp_id = ?", cfpID).Delete(nil)
	if result.Error != nil {
		return fmt.Errorf("failed to physically delete record from table cfp_infomation: %v", result.Error)
	}
	return nil
}
