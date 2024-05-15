package handler

import (
	"errors"
	"fmt"
	"net/http"

	"authenticator-backend/domain/common"
	"authenticator-backend/domain/model/traceability"
	"authenticator-backend/extension/logger"
	"authenticator-backend/usecase"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// IOperatorHandler
// Summary: This is interface which defines OperatorHandler.
//
//go:generate mockery --name IOperatorHandler --output ../../../../test/mock --case underscore
type IOperatorHandler interface {
	// #2 GetOperatorItem
	// #17 GetOperatorList
	GetOperator(c echo.Context) error
	// #1 PutOperatorItem
	PutOperator(c echo.Context) error
}

// operatorHandler
// Summary: This is structure which defines operatorHandler.
type operatorHandler struct {
	operatorUsecase usecase.IOperatorUsecase
}

// NewOperatorHandler
// Summary: This is function to create new operatorHandler.
// input: u(usecase.IOperatorUsecase) use case interface
// output: (IOperatorHandler) handler interface
func NewOperatorHandler(u usecase.IOperatorUsecase) IOperatorHandler {
	return &operatorHandler{u}
}

// GetOperator
// Summary: This is function which get a list of operator and response status.
// input: c(echo.Context) echo context
// output: (error) error object
func (h *operatorHandler) GetOperator(c echo.Context) error {
	var getOperatorInput traceability.GetOperatorInput
	dataTarget := c.QueryParam("dataTarget")
	method := c.Request().Method
	operatorID := c.Get("operatorID").(string)

	// if there are both operatorIds and openOperatorId, create error.
	isGetOperators := common.QueryParamExists(c, "operatorIds")
	if common.QueryParamExists(c, "openOperatorId") && isGetOperators {
		err := fmt.Errorf(common.OnlyOneCanBeSpecified("operatorIds", "openOperatorId"))
		logger.Set(c).Warnf(err.Error())

		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusBadRequest, common.HTTPErrorSourceAuth, common.Err400InvalidRequest, operatorID, dataTarget, method, err.Error()))
	}

	// if there is operatorIDs, search operator list.
	if isGetOperators {
		operatorIDs, err := common.QueryParamUUIDs(c, "operatorIds")
		if err != nil {
			logger.Set(c).Warnf(err.Error())
			errDetails := common.UnexpectedQueryParameter("operatorIds")

			return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusBadRequest, common.HTTPErrorSourceAuth, common.Err400InvalidRequest, operatorID, dataTarget, method, errDetails))
		}
		if len(operatorIDs) == 0 {
			logger.Set(c).Warnf(common.NotFoundError("operatorIds"))
			errDetails := common.UnexpectedQueryParameter("operatorIds")

			return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusBadRequest, common.HTTPErrorSourceAuth, common.Err400InvalidRequest, operatorID, dataTarget, method, errDetails))
		}

		// if length of operatorIDs exceeds upper limit, create error.
		if len(operatorIDs) > 100 {
			err := fmt.Errorf(common.UnexpectedQueryParameter("operatorIds"))
			logger.Set(c).Warnf(err.Error())

			return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusBadRequest, common.HTTPErrorSourceAuth, common.Err400InvalidRequest, operatorID, dataTarget, method, err.Error()))
		}

		var getOperatorsInput traceability.GetOperatorsInput = traceability.GetOperatorsInput{
			OperatorIDs: operatorIDs,
		}
		m, err := h.operatorUsecase.GetOperators(getOperatorsInput)
		if err != nil {
			logger.Set(c).Errorf(err.Error())

			return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusInternalServerError, common.HTTPErrorSourceAuth, common.Err500Unexpected, operatorID, dataTarget, method))
		}
		return c.JSON(http.StatusOK, m)
	}

	getOperatorInput.OperatorID = operatorID
	openOperatorID := common.QueryParamPtr(c, "openOperatorId")
	if openOperatorID != nil {
		getOperatorInput.OpenOperatorID = openOperatorID
	}

	if err := getOperatorInput.Validate(); err != nil {
		logger.Set(c).Warnf(err.Error())
		errDetails := common.UnexpectedQueryParameter("openOperatorId")

		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusBadRequest, common.HTTPErrorSourceAuth, common.Err400InvalidRequest, operatorID, dataTarget, method, errDetails))
	}

	m, err := h.operatorUsecase.GetOperator(getOperatorInput)
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

	return c.JSON(http.StatusOK, m)
}

// PutOperator
// Summary: This is function which put a object of Operator and response status.
// input: c(echo.Context) echo context
// output: (error) error object
func (h *operatorHandler) PutOperator(c echo.Context) error {
	var i traceability.PutOperatorInput
	dataTarget := c.QueryParam("dataTarget")
	method := c.Request().Method

	operatorID := c.Get("operatorID").(string)

	if err := c.Bind(&i); err != nil {
		logger.Set(c).Warnf(err.Error())
		errDetails := common.FormatBindErrMsg(err)

		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusBadRequest, common.HTTPErrorSourceAuth, common.Err400Validation, operatorID, dataTarget, method, errDetails))
	}
	if err := i.Validate(); err != nil {
		logger.Set(c).Warnf(err.Error())
		errDetails := err.Error()

		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusBadRequest, common.HTTPErrorSourceAuth, common.Err400Validation, operatorID, dataTarget, method, errDetails))
	}

	om, _ := i.ToModel()

	// checks coincidence of each operator values.
	operatorUUID, err := uuid.Parse(operatorID)
	if err != nil {
		logger.Set(c).Warnf(err.Error())

		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusBadRequest, common.HTTPErrorSourceAuth, common.Err401InvalidToken, operatorID, dataTarget, method))
	}
	if operatorUUID != om.OperatorID {
		logger.Set(c).Warnf(common.UnmatchValuesError("operator ID", "operator ID"))

		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusForbidden, common.HTTPErrorSourceAuth, common.Err403AccessDenied, operatorID, dataTarget, method))
	}

	m, err := h.operatorUsecase.PutOperator(om)
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
