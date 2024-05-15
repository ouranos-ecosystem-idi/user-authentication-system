package testhelper

import (
	"authenticator-backend/presentation/http/echo/handler"
	mocks "authenticator-backend/test/mock"
)

func NewMockHandler(host string) handler.OuranosHandler {

	operatorUsecase := new(mocks.IOperatorUsecase)
	operatorHandler := handler.NewOperatorHandler(operatorUsecase)
	plantUsecase := new(mocks.IPlantUsecase)
	plantHandler := handler.NewPlantHandler(plantUsecase)
	resetUsecase := new(mocks.IResetUsecase)
	resetHandler := handler.NewResetHandler(resetUsecase)
	h := handler.NewOuranosHandler(operatorHandler, plantHandler, resetHandler)

	return h
}
