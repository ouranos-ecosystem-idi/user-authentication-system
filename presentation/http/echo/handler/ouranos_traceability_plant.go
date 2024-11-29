package handler

import (
	"errors"
	"net/http"

	"authenticator-backend/domain/common"
	"authenticator-backend/domain/model/traceability"
	"authenticator-backend/extension/logger"
	"authenticator-backend/usecase"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// IPlantHandler
// Summary: This is interface which defines PlantHandler.
//
//go:generate mockery --name IPlantHandler --output ../../../../test/mock --case underscore
type IPlantHandler interface {
	// #4 GetPlantList
	GetPlant(c echo.Context) error
	// #3 PutPlantItem
	PutPlant(c echo.Context) error
}

// plantHandler
// Summary: This is structure which defines plantHandler.
type plantHandler struct {
	plantUsecase usecase.IPlantUsecase
}

// NewPlantHandler
// Summary: This is function to create new plantHandler.
// input: u(usecase.IPlantUsecase) use case interface
// output: (IPlantHandler) handler interface
func NewPlantHandler(u usecase.IPlantUsecase) IPlantHandler {
	return &plantHandler{u}
}

// GetPlant
// Summary: This is function which get a list of plant and response status.
// input: c(echo.Context) echo context
// output: (error) error object
func (h *plantHandler) GetPlant(c echo.Context) error {
	var getPlantModel traceability.GetPlantModel
	dataTarget := c.QueryParam("dataTarget")
	method := c.Request().Method
	operatorID := c.Get("operatorID").(string)

	OperatorUUID, err := uuid.Parse(operatorID)
	if err != nil {
		logger.Set(c).Warnf(err.Error())

		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusBadRequest, common.HTTPErrorSourceAuth, common.Err401InvalidToken, operatorID, dataTarget, method))
	}
	getPlantModel.OperatorID = OperatorUUID

	m, err := h.plantUsecase.ListPlants(getPlantModel)
	if err != nil {
		logger.Set(c).Errorf(err.Error())

		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusInternalServerError, common.HTTPErrorSourceAuth, common.Err500Unexpected, operatorID, dataTarget, method))
	}
	return c.JSON(http.StatusOK, m)
}

// PutPlant
// Summary: This is function which put a object of Plant and response status.
// input: c(echo.Context) echo context
// output: (error) error object
func (h *plantHandler) PutPlant(c echo.Context) error {
	var putPlantInput traceability.PutPlantInput
	dataTarget := c.QueryParam("dataTarget")
	method := c.Request().Method
	operatorID := c.Get("operatorID").(string)

	OperatorUUID, err := uuid.Parse(operatorID)
	if err != nil {
		logger.Set(c).Warnf(err.Error())

		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusBadRequest, common.HTTPErrorSourceAuth, common.Err401InvalidToken, operatorID, dataTarget, method))
	}

	if err := c.Bind(&putPlantInput); err != nil {
		logger.Set(c).Warnf(err.Error())
		errDetails := common.FormatBindErrMsg(err)

		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusBadRequest, common.HTTPErrorSourceAuth, common.Err400Validation, operatorID, dataTarget, method, errDetails))
	}

	if err := putPlantInput.Validate(); err != nil {
		logger.Set(c).Warnf(err.Error())
		errDetails := err.Error()

		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusBadRequest, common.HTTPErrorSourceAuth, common.Err400Validation, operatorID, dataTarget, method, errDetails))
	}
	// Verify that the request's OperatorID matches the jwt's OperatorID
	if OperatorUUID.String() != putPlantInput.OperatorID {
		logger.Set(c).Warnf(common.Err403AccessDenied)

		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusForbidden, common.HTTPErrorSourceAuth, common.Err403AccessDenied, operatorID, dataTarget, method))
	}

	putPlantModel, err := putPlantInput.ToModel()
	if err != nil {
		logger.Set(c).Warnf(err.Error())
		errDetails := err.Error()

		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusBadRequest, common.HTTPErrorSourceAuth, common.Err400Validation, operatorID, dataTarget, method, errDetails))
	}

	m, err := h.plantUsecase.PutPlant(putPlantModel)
	if err != nil {
		var customErr *common.CustomError
		if errors.As(err, &customErr) {
			if customErr.IsWarn() {
				logger.Set(c).Warnf(err.Error())
			} else {
				logger.Set(c).Errorf(err.Error())
			}

			return echo.NewHTTPError(common.HTTPErrorGenerate(int(customErr.Code), common.HTTPErrorSourceAuth, customErr.Message, operatorID, dataTarget, method))
		}
		logger.Set(c).Errorf(err.Error())

		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusInternalServerError, common.HTTPErrorSourceAuth, common.Err500Unexpected, operatorID, dataTarget, method))

	}
	return c.JSON(http.StatusCreated, m)
}
