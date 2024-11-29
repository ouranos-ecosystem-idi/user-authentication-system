package middleware

import (
	"authenticator-backend/extension/logger"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// dumpHandler
// Summary: This is the function which dump the request and response body.
// input: c(echo.Context): echo context
// input: reqBody([]byte): request body
// input: resBody([]byte): response body
func dumpHandler(c echo.Context, reqBody, resBody []byte) {
	if zap.S().Level() == zap.DebugLevel {
		logger.Set(c).Debugf("Auth Access Path: %v, Header: %v, Request Body: %v, Response Body: %v", c.Request().URL.String(), c.Request().Header, string(reqBody), string(resBody))
	} else {
		header := c.Request().Header
		for k := range header {
			if k == "Authorization" {
				header[k] = []string{"Bearer ******"}
			}
		}
		logger.Set(c).Infof("Auth Access Path: %v, Header: %v, Request Body: ******, Response Body: ******", c.Request().URL.String(), c.Request().Header)
	}
}
