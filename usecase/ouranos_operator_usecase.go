package usecase

import "authenticator-backend/domain/model/traceability"

// IOperatorUsecase
// Summary: This is interface which defines OperatorUsecase.
//
//go:generate mockery --name IOperatorUsecase --output ../test/mock --case underscore
type IOperatorUsecase interface {
	// #2 GetOperatorItem
	GetOperator(getOperatorInput traceability.GetOperatorInput) (traceability.OperatorModel, error)
	// #17 GetOperatorList
	GetOperators(getOperatorsInput traceability.GetOperatorsInput) ([]traceability.OperatorModel, error)
	// #1 PutOperatorItem
	PutOperator(operatorModel traceability.OperatorModel) (traceability.OperatorModel, error)
}
