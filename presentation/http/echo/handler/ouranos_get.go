package handler

import (
	"net/http"

	"authenticator-backend/domain/common"

	"github.com/labstack/echo/v4"
)

// GetAuthInfo
// Summary: This is the function which call the handler depending on the dataTarget query parameter.
// input: c(echo.Context): echo context
// output: (error) error object
func (h *ouranosHandler) GetAuthInfo(c echo.Context) error {
	method := c.Request().Method
	operatorId := c.Get("operatorID").(string)

	dataTarget := c.QueryParam("dataTarget")

	switch dataTarget {
	case "operator":
		return h.operatorHandler.GetOperator(c)
	case "plant":
		return h.plantHandler.GetPlant(c)
	default:
		errDetails := common.UnexpectedQueryParameter("dataTarget")
		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusBadRequest, common.HTTPErrorSourceAuth, common.Err400InvalidRequest, operatorId, dataTarget, method, errDetails))
	}
}
