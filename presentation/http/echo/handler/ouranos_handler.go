package handler

import "github.com/labstack/echo/v4"

// OuranosHandler
// Summary: This is the interface which defines the handler functions for the Ouranos.
type OuranosHandler interface {
	GetAuthInfo(c echo.Context) error
	PutAuthInfo(c echo.Context) error
	Reset(c echo.Context) error
	HealthCheck(c echo.Context) error
}

// ouranosHandler
// Summary: This is the structure which defines the handler for the Ouranos.
type ouranosHandler struct {
	operatorHandler IOperatorHandler
	plantHandler    IPlantHandler
	resetHandler    IResetHandler
}

// NewOuranosHandler
// Summary: This is the function which creates the Ouranos handler.
// input: operatorHandler(IOperatorHandler) operator handler
// input: plantHandler(IPlantHandler) plant handler
// input: resetHandler(IResetHandler) reset handler
// output: (OuranosHandler) Ouranos handler
func NewOuranosHandler(
	operatorHandler IOperatorHandler,
	plantHandler IPlantHandler,
	resetHandler IResetHandler,
) OuranosHandler {
	return &ouranosHandler{
		operatorHandler,
		plantHandler,
		resetHandler,
	}
}
