package repository

import (
	"authenticator-backend/domain/model/traceability"
)

// OuranosRepository
// Summary: This is interface which defines OuranosRepository functions.
//
//go:generate mockery --name OuranosRepository --output ../../test/mock --case underscore
type OuranosRepository interface {
	// Parts
	ListParts(getPlantPartsModel traceability.GetPartsModel) (traceability.PartsModels, error)
	DeleteParts(traceID string) error

	// PartsStructure
	DeletePartsStructure(traceID string) error

	// Operator
	GetOperator(operatorID string) (traceability.OperatorEntityModel, error)
	GetOperatorByOpenOperatorID(openOperatorID string) (traceability.OperatorEntityModel, error)
	GetOperators(operatorIDs []string) (traceability.OperatorEntityModels, error)
	PutOperator(e traceability.OperatorEntityModel) (traceability.OperatorEntityModel, error)

	// Plant
	CreatePlant(e traceability.PlantEntityModel) (traceability.PlantEntityModel, error)
	GetPlant(operatorID string, plantID string) (traceability.PlantEntityModel, error)
	ListPlantsByOperatorID(operatorID string) (traceability.PlantEntityModels, error)
	UpdatePlant(e traceability.PlantEntityModel) (traceability.PlantEntityModel, error)
	DeletePlant(plantID string) error

	// Trade
	ListTradesByOperatorID(operatorID string) (traceability.TradeEntityModels, error)
	DeleteTrade(tradeID string) error

	// RequestStatus
	DeleteRequestStatusByTradeID(tradeID string) error

	// CFPInfomation
	GetCFPInformation(traceID string) (traceability.CfpEntityModel, error)
	DeleteCFPInformation(cfpID string) error

	// CFPCertificate
	DeleteCFPCertificateByCFPID(cfpID string) error
}
