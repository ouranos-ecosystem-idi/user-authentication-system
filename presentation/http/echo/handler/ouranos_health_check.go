package handler

import (
	"authenticator-backend/domain/common"
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthCheck
// Summary: This is the function which performs a helth check.
// input: c(echo.Context): echo context
// output: (error) error object
func (h *ouranosHandler) HealthCheck(c echo.Context) error {
	healthCheckResponse := common.HealthCheckResponse{
		IsSystemHealthy: true,
	}
	return c.JSON(http.StatusOK, healthCheckResponse)
}
