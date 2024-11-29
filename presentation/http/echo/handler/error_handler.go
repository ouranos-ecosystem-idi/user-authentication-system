package handler

import (
	"errors"
	"net/http"
	"syscall"

	"authenticator-backend/domain/common"

	"github.com/jackc/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {

	method := c.Request().Method
	dataTarget := c.QueryParam("dataTarget")

	if errors.Is(err, middleware.ErrJWTMissing) {
		c.Error(echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusUnauthorized, common.HTTPErrorSourceAuth, common.Err401Authentication, "", dataTarget, method)))

		return
	}
	if errors.Is(err, syscall.ECONNREFUSED) {
		c.Error(echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusServiceUnavailable, common.HTTPErrorSourceAuth, common.Err503OuterService, "", dataTarget, method)))

		return
	}
	if he, ok := err.(*echo.HTTPError); ok {
		if he.Code == http.StatusNotFound && he.Message == "Not Found" {
			c.Error(echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusNotFound, common.HTTPErrorSourceDataspace, common.Err404EndpointNotFound, "", dataTarget, method)))

			return
		}
	} else if pe, ok := err.(*pgconn.PgError); ok {
		if pe.Code == common.PgErrorAdminShutdown || pe.Code == common.PgErrorCrashShutdown {
			c.Error(echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusServiceUnavailable, common.HTTPErrorSourceAuth, common.Err503OuterService, "", dataTarget, method)))

			return
		} else {
			c.Error(echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusInternalServerError, common.HTTPErrorSourceAuth, common.Err500Unexpected, "", dataTarget, method)))

			return
		}
	} else {
		c.Error(echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusInternalServerError, common.HTTPErrorSourceAuth, common.Err500Unexpected, "", dataTarget, method)))
		return
	}

	c.Echo().DefaultHTTPErrorHandler(err, c)
}
