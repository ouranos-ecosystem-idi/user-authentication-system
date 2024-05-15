package middleware

import (
	"errors"
	"net/http"
	"strings"

	"authenticator-backend/domain/common"
	"authenticator-backend/extension/logger"
	"authenticator-backend/usecase"
	"authenticator-backend/usecase/input"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware
// Summary: This is the structure which defines the auth middleware.
type AuthMiddleware struct {
	verifyUsecase usecase.IVerifyUsecase
}

// NewAuthMiddleware
// Summary: This is the function which creates the auth middleware.
// input: u(usecase.IVerifyUsecase): verify usecase
// output: (AuthMiddleware) auth middleware
func NewAuthMiddleware(u usecase.IVerifyUsecase) AuthMiddleware {
	return AuthMiddleware{u}
}

// AuthJWT
// Summary: This is the function which authenticates the JWT.
// output: (echo.MiddlewareFunc) echo middleware function
func (m AuthMiddleware) AuthJWT() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// prepare
			method := c.Request().Method
			dataTarget := c.QueryParam("dataTarget")

			idTokenRaw := c.Request().Header.Get("Authorization")
			if idTokenRaw == "" {
				logger.Set(c).Warnf(common.Err401Authentication)

				return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusUnauthorized, common.HTTPErrorSourceAuth, common.Err401Authentication, "", dataTarget, method))
			}
			idToken := strings.TrimPrefix(idTokenRaw, "Bearer ")

			// verify idToken
			claim, err := m.verifyUsecase.IDToken(input.VerifyIDTokenParam{IDToken: idToken})
			if err != nil {
				var customErr *common.CustomError
				if errors.As(err, &customErr) {
					logger.Set(c).Warnf(customErr.Message)

					return echo.NewHTTPError(common.HTTPErrorGenerate(int(customErr.Code), common.HTTPErrorSourceAuth, customErr.Message, "", dataTarget, method))
				}

				logger.Set(c).Error(err.Error())

				return echo.NewHTTPError(common.HTTPErrorGenerate(http.StatusUnauthorized, common.HTTPErrorSourceAuth, common.Err401InvalidToken, "", dataTarget, method))
			}

			// set context
			c.Set("operator", &claim)
			c.Set("operatorID", claim.OperatorID)

			return next(c)
		}
	}
}
