package middleware

import (
	"net/http"

	"authenticator-backend/domain/common"
	"authenticator-backend/domain/model/authentication"
	"authenticator-backend/domain/repository"
	"authenticator-backend/extension/logger"
	"authenticator-backend/infrastructure/persistence/datastore"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const apiKeyHeader = "apiKey"

// APIKeyValidator
// Summary: This is the function which validates the API key.
// input: db(*gorm.DB): database
// output: (echo.MiddlewareFunc) middleware function
func APIKeyValidator(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authRepository := datastore.NewAuthRepository(db)
			method := c.Request().Method
			apiKey := c.Request().Header.Get(apiKeyHeader)

			validAPIKeys, err := authRepository.ListAPIKeys(repository.APIKeysParam{})
			if err != nil {
				logger.Set(c).Errorf(common.Err500Unexpected)

				return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusInternalServerError, common.HTTPErrorSourceAuth, common.Err500Unexpected, "", "", method))
			}

			if apiKey == "" {
				logger.Set(c).Warnf(common.Err403AccessDenied)

				return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusForbidden, common.HTTPErrorSourceAuth, common.Err403AccessDenied, "", "", method))
			}

			// check if the API key is in the ValidAPIKeys
			if !validAPIKeys.ContainsAPIKey(apiKey) {
				logger.Set(c).Warnf(common.Err403InvalidKey)

				return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusForbidden, common.HTTPErrorSourceAuth, common.Err403InvalidKey, "", "", method))
			}

			return next(c)
		}
	}
}

// ApplicationAttributeがDataSpaceとTraceabilityのAPIKEYのみ許可

// SystemAPIKeyValidator
// Summary: This is the function which validates the system API key.
// input: db(*gorm.DB): database
// output: (echo.MiddlewareFunc) middleware function
func SystemAPIKeyValidator(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authRepository := datastore.NewAuthRepository(db)
			method := c.Request().Method
			apiKey := c.Request().Header.Get(apiKeyHeader)

			// Only "DataSpace" and "Traceability" Attribute API keys are valid
			systemUserAttributes := []authentication.ApplicationAttribute{authentication.ApplicationAttributeTraceability, authentication.ApplicationAttributeDataSpace}
			validAPIKeys, err := authRepository.ListAPIKeys(repository.APIKeysParam{
				Attributes: systemUserAttributes,
			})
			if err != nil {
				logger.Set(c).Errorf(common.Err500Unexpected)

				return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusInternalServerError, common.HTTPErrorSourceAuth, common.Err500Unexpected, "", "", method))
			}

			if apiKey == "" {
				logger.Set(c).Warnf(common.Err403AccessDenied)

				return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusForbidden, common.HTTPErrorSourceAuth, common.Err403AccessDenied, "", "", method))
			}

			// check if the API key is in the ValidAPIKeys
			if !validAPIKeys.ContainsAPIKey(apiKey) {
				logger.Set(c).Warnf(common.Err403InvalidKey)

				return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusForbidden, common.HTTPErrorSourceAuth, common.Err403InvalidKey, "", "", method))
			}

			return next(c)
		}
	}
}
