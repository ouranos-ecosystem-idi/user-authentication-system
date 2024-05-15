package usecase

import "authenticator-backend/domain/model/traceability"

//go:generate mockery --name IPlantUsecase --output ../test/mock --case underscore
type IPlantUsecase interface {
	// #4 GetPlantList
	ListPlants(getPlantModel traceability.GetPlantModel) ([]traceability.PlantModel, error)
	// #3 PutPlantItem
	PutPlant(plantModel traceability.PlantModel) (traceability.PlantModel, error)
}
