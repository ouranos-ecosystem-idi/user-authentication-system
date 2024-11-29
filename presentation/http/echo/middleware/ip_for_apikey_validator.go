package middleware

import (
	"net/http"
	"os"
	"strings"

	"authenticator-backend/domain/common"
	"authenticator-backend/domain/repository"
	"authenticator-backend/extension/logger"
	"authenticator-backend/infrastructure/persistence/datastore"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// IPForAPIKeyValidator
// Summary: This is the function which validates the IP address related to the APIkey.
// input: db(*gorm.DB): database
// output: (echo.MiddlewareFunc) middleware function
func IPForAPIKeyValidator(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authRepository := datastore.NewAuthRepository(db)
			method := c.Request().Method
			apiKey := c.Request().Header.Get(apiKeyHeader)

			env := os.Getenv("GO_ENV")
			var ip string
			if env == "local" {
				ip = c.RealIP()
			} else {
				xff := c.Request().Header.Get("X-Forwarded-For") // X-Forwarded-For: APP_IP, ALB_IP
				ips := strings.Split(xff, ", ")
				if len(ips) < 2 {
					return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusForbidden, common.HTTPErrorSourceAuth, common.Err403IPNotAuthorizedForKey, "", "", method))
				}
				ip = ips[len(ips)-2]
			}

			dummyBody := dummyBodyApikeyIp{APIKey: apiKey, IP: ip}

			cidrs, err := authRepository.ListCidrs(repository.APIKeyCidrsParam{APIKey: &apiKey})
			if err != nil {
				logger.Set(c).Warnf(common.Err403IPNotAuthorizedForKey)
				authDump(c, dummyBody, nil, eventAPIKey, false)

				return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusForbidden, common.HTTPErrorSourceAuth, common.Err403IPNotAuthorizedForKey, "", "", method))
			}
			if !cidrs.Contains(ip) {
				logger.Set(c).Warnf(common.Err403IPNotAuthorizedForKey)
				authDump(c, dummyBody, nil, eventAPIKey, false)

				return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusForbidden, common.HTTPErrorSourceAuth, common.Err403IPNotAuthorizedForKey, "", "", method))
			}

			authDump(c, dummyBody, nil, eventAPIKey, true)
			return next(c)
		}
	}
}

// dummyBodyApikeyIp
// Summary: This is the structure which defines the dummy body for API key and IP address.
type dummyBodyApikeyIp struct {
	APIKey string `json:"apiKey"`
	IP     string `json:"ipAddress"`
}
