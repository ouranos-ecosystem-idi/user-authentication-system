package handler

import (
	"errors"
	"net/http"

	"authenticator-backend/domain/common"
	"authenticator-backend/domain/model/authentication"
	"authenticator-backend/extension/logger"
	"authenticator-backend/usecase/input"

	"github.com/labstack/echo/v4"
)

// Login
// Summary: This is function which is used to login
// input: c(echo.Context): context
// output: error: error object
func (h *authHandler) Login(c echo.Context) error {
	method := c.Request().Method
	param := input.LoginParam{}

	if err := c.Bind(&param); err != nil {
		logger.Set(c).Warnf(err.Error())

		errDetails := common.FormatBindErrMsg(err)
		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusBadRequest, common.HTTPErrorSourceAuth, common.Err400Validation, "", "", method, errDetails))
	}

	err := param.Validate()
	if err != nil {
		logger.Set(c).Warnf(err.Error())

		errDetails := err.Error()
		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusBadRequest, common.HTTPErrorSourceAuth, common.Err400Validation, "", "", method, errDetails))
	}

	output, err := h.AuthUsecase.Login(param)
	if err != nil {
		var customErr *common.CustomError
		if errors.As(err, &customErr) {
			if customErr.IsWarn() {
				logger.Set(c).Warnf(err.Error())
			} else {
				logger.Set(c).Errorf(err.Error())
			}

			return echo.NewHTTPError(common.HTTPErrorGenerate(int(customErr.Code), common.HTTPErrorSourceAuth, customErr.Message, "", "", method))
		}
		logger.Set(c).Errorf(err.Error())

		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusInternalServerError, common.HTTPErrorSourceAuth, common.Err500Unexpected, "", "", method))
	}
	return c.JSON(http.StatusCreated, output)
}

// Refresh
// Summary: This is function which is used to refresh token
// input: c(echo.Context): context
// output: error: error object
func (h *authHandler) Refresh(c echo.Context) error {
	method := c.Request().Method
	param := input.RefreshParam{}

	if err := c.Bind(&param); err != nil {
		logger.Set(c).Warnf(err.Error())

		errDetails := common.FormatBindErrMsg(err)
		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusBadRequest, common.HTTPErrorSourceAuth, common.Err400Validation, "", "", method, errDetails))
	}

	err := param.Validate()
	if err != nil {
		logger.Set(c).Warnf(err.Error())

		errDetails := err.Error()
		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusBadRequest, common.HTTPErrorSourceAuth, common.Err400Validation, "", "", method, errDetails))
	}

	output, err := h.AuthUsecase.Refresh(param)
	if err != nil {
		var customErr *common.CustomError
		if errors.As(err, &customErr) {
			if customErr.IsWarn() {
				logger.Set(c).Warnf(err.Error())
			} else {
				logger.Set(c).Errorf(err.Error())
			}

			return echo.NewHTTPError(common.HTTPErrorGenerate(int(customErr.Code), common.HTTPErrorSourceAuth, customErr.Message, "", "", method))
		}
		logger.Set(c).Errorf(err.Error())

		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusInternalServerError, common.HTTPErrorSourceAuth, common.Err500Unexpected, "", "", method))
	}
	return c.JSON(http.StatusCreated, output)
}

// ChangePassword
// Summary: This is function which is used to change password
// input: c(echo.Context): context
// output: error: error object
func (h *authHandler) ChangePassword(c echo.Context) error {
	method := c.Request().Method

	claims := c.Get("operator").(*authentication.Claims)
	operatorId := claims.OperatorID
	uid := claims.UID

	var param input.ChangePasswordParam
	if err := c.Bind(&param); err != nil {
		logger.Set(c).Warnf(err.Error())
		errDetails := common.FormatBindErrMsg(err)

		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusBadRequest, common.HTTPErrorSourceAuth, common.Err400Validation, operatorId, "", method, errDetails))
	}
	param.UID = uid

	if err := param.Validate(); err != nil {
		logger.Set(c).Warnf(err.Error())
		errDetails := err.Error()

		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusBadRequest, common.HTTPErrorSourceAuth, common.Err400Validation, operatorId, "", method, errDetails))
	}

	if err := h.AuthUsecase.ChangePassword(param); err != nil {
		logger.Set(c).Errorf(err.Error())

		return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusInternalServerError, common.HTTPErrorSourceAuth, common.Err500Unexpected, operatorId, "", method))
	}

	return c.JSON(http.StatusCreated, common.EmptyBody{})
}
