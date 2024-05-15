package datastore

import (
	"fmt"

	"authenticator-backend/domain/model/traceability"

	"github.com/google/uuid"
)

func (r *ouranosRepository) ListParts(getPlantPartsModel traceability.GetPartsModel) (traceability.PartsModels, error) {
	var (
		partsList traceability.PartsModels
		err       error
	)

	query := r.db.Table("parts").
		Select(`
			parts.trace_id,
			parts.operator_id,
			parts.plant_id,
			parts.parts_name,
			parts.support_parts_name,
			parts.terminated_flag,
			parts.amount_required,
			parts.amount_required_unit
		`).
		Where(`parts.deleted_at IS NULL AND parts.operator_id = ?`, getPlantPartsModel.OperatorID)

	if getPlantPartsModel.TraceID != nil {
		query = query.Where("parts.trace_id = ?", *getPlantPartsModel.TraceID)
	}
	if getPlantPartsModel.PartsName != nil {
		query = query.Where("parts.parts_name = ?", *getPlantPartsModel.PartsName)
	}
	if getPlantPartsModel.PlantID != nil {
		query = query.Where("parts.plant_id = ?", *getPlantPartsModel.PlantID)
	}
	if getPlantPartsModel.ParentFlag != nil {
		if *getPlantPartsModel.ParentFlag {
			query = query.
				Joins("INNER JOIN parts_structures ON parts_structures.trace_id = parts.trace_id").
				Where("parts_structures.parent_trace_id = ?", uuid.Nil.String())
		}
	}

	err = query.
		Limit(getPlantPartsModel.Limit).
		Order(`parts_name ASC`).
		Order(`support_parts_name ASC`).
		Find(&partsList).
		Error

	return partsList, err
}

func (r *ouranosRepository) DeleteParts(traceID string) error {
	result := r.db.Unscoped().Table("parts").Where("trace_id = ?", traceID).Delete(nil)
	if result.Error != nil {
		return fmt.Errorf("failed to physically delete record from table parts: %v", result.Error)
	}
	return nil
}
