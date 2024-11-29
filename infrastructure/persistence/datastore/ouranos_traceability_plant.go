package datastore

import (
	"fmt"

	"authenticator-backend/domain/common"
	"authenticator-backend/domain/model/traceability"
	"authenticator-backend/extension/logger"

	"gorm.io/gorm"
)

// CreatePlant
// Summary: This is function which insert Plant to data store.
// input: e(traceability.PlantEntityModel) PlantEntityModel object
// output: (traceability.PlantEntityModel) PlantEntityModel Object
// output: (error) error object
func (r *ouranosRepository) CreatePlant(e traceability.PlantEntityModel) (traceability.PlantEntityModel, error) {
	if result := r.db.Table("plants").Create(&e); result.Error != nil {
		return traceability.PlantEntityModel{}, result.Error
	}
	return e, nil
}

// GetPlant
// Summary: This is function which get Plant from data store with operatorID and plantID.
// input: operatorID(string) Value of operatorID
// input: plantID(string) Value of plantID
// output: (traceability.PlantEntityModel) PlantEntityModel Object
// output: (error) error object
func (r *ouranosRepository) GetPlant(operatorID string, plantID string) (traceability.PlantEntityModel, error) {
	var plant traceability.PlantEntityModel

	if err := r.db.Table("plants").Where("operator_id = ?", operatorID).Where("plant_id = ?", plantID).First(&plant).Error; err != nil {
		logger.Set(nil).Errorf(err.Error())
		return traceability.PlantEntityModel{}, err
	}
	return plant, nil
}

// ListPlantsByOperatorID
// Summary: This is function which get array of Plant from data store with operatorID.
// input: operatorID(string) Value of operatorID
// output: (traceability.PlantEntityModels) Array of PlantEntityModel
// output: (error) error object
func (r *ouranosRepository) ListPlantsByOperatorID(operatorID string) (traceability.PlantEntityModels, error) {
	var plants traceability.PlantEntityModels

	if err := r.db.Table("plants").Where("operator_id = ?", operatorID).Find(&plants).Error; err != nil {
		logger.Set(nil).Errorf(err.Error())
		return traceability.PlantEntityModels{}, err
	}
	return plants, nil
}

// UpdatePlant
// Summary: This is function which update Plant to data store with plantID.
// input: e(traceability.PlantEntityModel) PlantEntityModel object
// output: (traceability.PlantEntityModel) PlantEntityModel object
// output: (error) error object
func (r *ouranosRepository) UpdatePlant(e traceability.PlantEntityModel) (traceability.PlantEntityModel, error) {
	updateValues := map[string]interface{}{
		"plant_id":        e.PlantID,
		"operator_id":     e.OperatorID,
		"plant_name":      e.PlantName,
		"plant_address":   e.PlantAddress,
		"open_plant_id":   e.OpenPlantID,
		"global_plant_id": e.GlobalPlantID,
	}
	var updated traceability.PlantEntityModel
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("plants").Where("plant_id = ?", e.PlantID).Updates(updateValues).Error; err != nil {
			logger.Set(nil).Errorf(err.Error())
			return err
		}
		if err := tx.Table("plants").Where("plant_id = ?", e.PlantID).First(&updated).Error; err != nil {
			logger.Set(nil).Errorf(err.Error())
			return err
		}
		return nil
	})

	if err != nil {
		return traceability.PlantEntityModel{}, err
	}
	return updated, nil
}

// DeletePlant
// Summary: This is function which delete Plant from data store with plantID.
// input: plantID(string) Value of plantID
// output: (error) error object
func (r *ouranosRepository) DeletePlant(plantID string) error {
	result := r.db.Unscoped().Table("plants").Where("plant_id = ?", plantID).Delete(nil)
	if result.Error != nil {
		return fmt.Errorf(common.DeleteTableError("plants", result.Error))
	}
	return nil
}
