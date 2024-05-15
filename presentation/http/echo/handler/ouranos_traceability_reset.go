package handler

import (
	"net/http"

	"authenticator-backend/domain/common"
	"authenticator-backend/extension/logger"
	"authenticator-backend/usecase"

	"github.com/labstack/echo/v4"
)

// IResetHandler
// Summary: This is the interface for the reset handler.
//
//go:generate mockery --name IResetHandler --output ../../../../test/mock --case underscore
type IResetHandler interface {
	Reset(c echo.Context) error
}

// resetHandler
// Summary: This is the structure which defines the reset handler.
type resetHandler struct {
	resetUsecase usecase.IResetUsecase
}

// NewResetHandler
// Summary: This is the function which creates the reset handler.
// input: u(usecase.IResetUsecase): reset usecase
// output: (IResetHandler) reset handler
func NewResetHandler(u usecase.IResetUsecase) IResetHandler {
	return &resetHandler{u}
}

// Reset
// Summary: This is the function which resets the data.
// input: c(echo.Context): echo context
// output: (error) error object
func (h *ouranosHandler) Reset(c echo.Context) error {
	return h.resetHandler.Reset(c)
}

// Reset
// Summary: This is the function which resets the data.
// input: c(echo.Context): echo context
// output: (error) error object
func (h *resetHandler) Reset(c echo.Context) error {
	method := c.Request().Method
	apikey := c.Request().Header.Get("apiKey")
	operatorID := c.Get("operatorID").(string)

	err := h.resetUsecase.Reset(apikey)
	if err != nil {
		logger.Set(c).Error(err.Error())

		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusInternalServerError, common.HTTPErrorSourceAuth, common.Err500Unexpected, operatorID, "dataReset", method))
	}
	return c.JSON(http.StatusOK, nil)
}
